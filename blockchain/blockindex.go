// Copyright (c) 2013-2017 The btcsuite developers
// Copyright (c) 2018-2019 The Hdfchain developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"bytes"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/hdfchain/hdfd/blockchain/stake/v3"
	"github.com/hdfchain/hdfd/blockchain/standalone/v2"
	"github.com/hdfchain/hdfd/chaincfg/chainhash"
	"github.com/hdfchain/hdfd/database/v2"
	"github.com/hdfchain/hdfd/wire"
)

// blockStatus is a bit field representing the validation state of the block.
type blockStatus byte

// The following constants specify possible status bit flags for a block.
//
// NOTE: This section specifically does not use iota since the block status is
// serialized and must be stable for long-term storage.
const (
	// statusNone indicates that the block has no validation state flags set.
	statusNone blockStatus = 0

	// statusDataStored indicates that the block's payload is stored on disk.
	statusDataStored blockStatus = 1 << 0

	// statusValidated indicates that the block has been fully validated.  It
	// also means that all of its ancestors have also been validated.
	statusValidated blockStatus = 1 << 1

	// statusValidateFailed indicates that the block has failed validation.
	statusValidateFailed blockStatus = 1 << 2

	// statusInvalidAncestor indicates that one of the ancestors of the block
	// has failed validation, thus the block is also invalid.
	statusInvalidAncestor blockStatus = 1 << 3
)

// HaveData returns whether the full block data is stored in the database.  This
// will return false for a block node where only the header is downloaded or
// stored.
func (status blockStatus) HaveData() bool {
	return status&statusDataStored != 0
}

// HasValidated returns whether the block is known to have been successfully
// validated.  A return value of false in no way implies the block is invalid.
// Thus, this will return false for a valid block that has not been fully
// validated yet.
//
// NOTE: A block that is known to have been validated might also be marked as
// known invalid as well if the block is manually invalidated.
func (status blockStatus) HasValidated() bool {
	return status&statusValidated != 0
}

// KnownInvalid returns whether either the block itself is known to be invalid
// or to have an invalid ancestor.  A return value of false in no way implies
// the block is valid or only has valid ancestors.  Thus, this will return false
// for invalid blocks that have not been proven invalid yet as well as return
// false for blocks with invalid ancestors that have not been proven invalid
// yet.
//
// NOTE: A block that is known invalid might also be marked as known to have
// been successfully validated as well if the block is manually invalidated.
func (status blockStatus) KnownInvalid() bool {
	return status&(statusValidateFailed|statusInvalidAncestor) != 0
}

// blockNode represents a block within the block chain and is primarily used to
// aid in selecting the best chain to be the main chain.  The main chain is
// stored into the block database.
type blockNode struct {
	// NOTE: Additions, deletions, or modifications to the order of the
	// definitions in this struct should not be changed without considering
	// how it affects alignment on 64-bit platforms.  The current order is
	// specifically crafted to result in minimal padding.  There will be
	// hundreds of thousands of these in memory, so a few extra bytes of
	// padding adds up.

	// parent is the parent block for this node.
	parent *blockNode

	// skipToAncestor is used to provide a skip list to significantly speed up
	// traversal to ancestors deep in history.
	skipToAncestor *blockNode

	// hash is the hash of the block this node represents.
	hash chainhash.Hash

	// workSum is the total amount of work in the chain up to and including
	// this node.
	workSum *big.Int

	// Some fields from block headers to aid in best chain selection and
	// reconstructing headers from memory.  These must be treated as
	// immutable and are intentionally ordered to avoid padding on 64-bit
	// platforms.
	height       int64
	voteBits     uint16
	finalState   [6]byte
	blockVersion int32
	voters       uint16
	freshStake   uint8
	revocations  uint8
	poolSize     uint32
	bits         uint32
	sbits        int64
	timestamp    int64
	merkleRoot   chainhash.Hash
	stakeRoot    chainhash.Hash
	blockSize    uint32
	nonce        uint32
	extraData    [32]byte
	stakeVersion uint32

	// status is a bitfield representing the validation state of the block.
	// This field, unlike most other fields, may be changed after the block
	// node is created, so it must only be accessed or updated using the
	// concurrent-safe NodeStatus, SetStatusFlags, and UnsetStatusFlags
	// methods on blockIndex once the node has been added to the index.
	status blockStatus

	// stakeNode contains all the consensus information required for the
	// staking system.  The node also caches information required to add or
	// remove stake nodes, so that the stake node itself may be pruneable
	// to save memory while maintaining high throughput efficiency for the
	// evaluation of sidechains.
	stakeNode      *stake.Node
	newTickets     []chainhash.Hash
	ticketsVoted   []chainhash.Hash
	ticketsRevoked []chainhash.Hash

	// Keep track of all vote version and bits in this block.
	votes []stake.VoteVersionTuple
}

// clearLowestOneBit clears the lowest set bit in the passed value.
func clearLowestOneBit(n int64) int64 {
	return n & (n - 1)
}

// calcSkipListHeight calculates the height of an ancestor block to use when
// constructing the ancestor traversal skip list.
func calcSkipListHeight(height int64) int64 {
	if height < 0 {
		return 0
	}

	// Traditional skip lists create multiple levels to achieve expected average
	// search, insert, and delete costs of O(log n).  Since the blockchain is
	// append only, there is no need to handle random insertions or deletions,
	// so this takes advantage of that to effectively create a deterministic
	// skip list with a single level that is reasonably close to O(log n) in
	// order to reduce the number of pointers and implementation complexity.
	//
	// This calculation is definitely not the most optimal possible in terms of
	// the number of steps in the worst case, however, it is predominantly
	// logarithmic, easy to reason about, deterministic, blazing fast to
	// calculate and can easily be shown to have a worst case performance of
	// 420 steps for heights up to 4,294,967,296 (2^32) and 1580 steps for
	// heights up to 2^63 - 1.
	//
	// Finally, it also satisfies the only real requirement for proper operation
	// of the skip list which is for the calculated height to be less than the
	// provided height.
	return clearLowestOneBit(clearLowestOneBit(height))
}

// initBlockNode initializes a block node from the given header, initialization
// vector for the ticket lottery, and parent node.  The workSum is calculated
// based on the parent, or, in the case no parent is provided, it will just be
// the work for the passed block.
//
// This function is NOT safe for concurrent access.  It must only be called when
// initially creating a node.
func initBlockNode(node *blockNode, blockHeader *wire.BlockHeader, parent *blockNode) {
	*node = blockNode{
		hash:         blockHeader.BlockHash(),
		workSum:      standalone.CalcWork(blockHeader.Bits),
		height:       int64(blockHeader.Height),
		blockVersion: blockHeader.Version,
		voteBits:     blockHeader.VoteBits,
		finalState:   blockHeader.FinalState,
		voters:       blockHeader.Voters,
		freshStake:   blockHeader.FreshStake,
		poolSize:     blockHeader.PoolSize,
		bits:         blockHeader.Bits,
		sbits:        blockHeader.SBits,
		timestamp:    blockHeader.Timestamp.Unix(),
		merkleRoot:   blockHeader.MerkleRoot,
		stakeRoot:    blockHeader.StakeRoot,
		revocations:  blockHeader.Revocations,
		blockSize:    blockHeader.Size,
		nonce:        blockHeader.Nonce,
		extraData:    blockHeader.ExtraData,
		stakeVersion: blockHeader.StakeVersion,
		status:       statusNone,
	}
	if parent != nil {
		node.parent = parent
		node.skipToAncestor = parent.Ancestor(calcSkipListHeight(node.height))
		node.workSum = node.workSum.Add(parent.workSum, node.workSum)
	}
}

// newBlockNode returns a new block node for the given block header and parent
// node.  The workSum is calculated based on the parent, or, in the case no
// parent is provided, it will just be the work for the passed block.
func newBlockNode(blockHeader *wire.BlockHeader, parent *blockNode) *blockNode {
	var node blockNode
	initBlockNode(&node, blockHeader, parent)
	return &node
}

// Header constructs a block header from the node and returns it.
//
// This function is safe for concurrent access.
func (node *blockNode) Header() wire.BlockHeader {
	// No lock is needed because all accessed fields are immutable.
	prevHash := zeroHash
	if node.parent != nil {
		prevHash = &node.parent.hash
	}
	return wire.BlockHeader{
		Version:      node.blockVersion,
		PrevBlock:    *prevHash,
		MerkleRoot:   node.merkleRoot,
		StakeRoot:    node.stakeRoot,
		VoteBits:     node.voteBits,
		FinalState:   node.finalState,
		Voters:       node.voters,
		FreshStake:   node.freshStake,
		Revocations:  node.revocations,
		PoolSize:     node.poolSize,
		Bits:         node.bits,
		SBits:        node.sbits,
		Height:       uint32(node.height),
		Size:         node.blockSize,
		Timestamp:    time.Unix(node.timestamp, 0),
		Nonce:        node.nonce,
		ExtraData:    node.extraData,
		StakeVersion: node.stakeVersion,
	}
}

// lotteryIV returns the initialization vector for the deterministic PRNG used
// to determine winning tickets.
//
// This function is safe for concurrent access.
func (node *blockNode) lotteryIV() chainhash.Hash {
	// Serialize the block header for use in calculating the initialization
	// vector for the ticket lottery.  The only way this can fail is if the
	// process is out of memory in which case it would panic anyways, so
	// although panics are generally frowned upon in package code, it is
	// acceptable here.
	buf := bytes.NewBuffer(make([]byte, 0, wire.MaxBlockHeaderPayload))
	header := node.Header()
	if err := header.Serialize(buf); err != nil {
		panic(err)
	}

	return stake.CalcHash256PRNGIV(buf.Bytes())
}

// populateTicketInfo sets pruneable ticket information in the provided block
// node.
//
// This function is NOT safe for concurrent access.  It must only be called when
// initially creating a node or when protected by the chain lock.
func (node *blockNode) populateTicketInfo(spentTickets *stake.SpentTicketsInBlock) {
	node.ticketsVoted = spentTickets.VotedTickets
	node.ticketsRevoked = spentTickets.RevokedTickets
	node.votes = spentTickets.Votes
}

// Ancestor returns the ancestor block node at the provided height by following
// the chain backwards from this node.  The returned block will be nil when a
// height is requested that is after the height of the passed node or is less
// than zero.
//
// This function is safe for concurrent access.
func (node *blockNode) Ancestor(height int64) *blockNode {
	if height < 0 || height > node.height {
		return nil
	}

	n := node
	for n != nil && n.height != height {
		// Skip to the linked ancestor when it won't overshoot the target
		// height.
		if n.skipToAncestor != nil && calcSkipListHeight(n.height) >= height {
			n = n.skipToAncestor
			continue
		}

		n = n.parent
	}

	return n
}

// RelativeAncestor returns the ancestor block node a relative 'distance' blocks
// before this node.  This is equivalent to calling Ancestor with the node's
// height minus provided distance.
//
// This function is safe for concurrent access.
func (node *blockNode) RelativeAncestor(distance int64) *blockNode {
	return node.Ancestor(node.height - distance)
}

// CalcPastMedianTime calculates the median time of the previous few blocks
// prior to, and including, the block node.
//
// This function is safe for concurrent access.
func (node *blockNode) CalcPastMedianTime() time.Time {
	// Create a slice of the previous few block timestamps used to calculate
	// the median per the number defined by the constant medianTimeBlocks.
	timestamps := make([]int64, medianTimeBlocks)
	numNodes := 0
	iterNode := node
	for i := 0; i < medianTimeBlocks && iterNode != nil; i++ {
		timestamps[i] = iterNode.timestamp
		numNodes++

		iterNode = iterNode.parent
	}

	// Prune the slice to the actual number of available timestamps which
	// will be fewer than desired near the beginning of the block chain
	// and sort them.
	timestamps = timestamps[:numNodes]
	sort.Sort(timeSorter(timestamps))

	// NOTE: The consensus rules incorrectly calculate the median for even
	// numbers of blocks.  A true median averages the middle two elements
	// for a set with an even number of elements in it.   Since the constant
	// for the previous number of blocks to be used is odd, this is only an
	// issue for a few blocks near the beginning of the chain.  I suspect
	// this is an optimization even though the result is slightly wrong for
	// a few of the first blocks since after the first few blocks, there
	// will always be an odd number of blocks in the set per the constant.
	//
	// This code follows suit to ensure the same rules are used, however, be
	// aware that should the medianTimeBlocks constant ever be changed to an
	// even number, this code will be wrong.
	medianTimestamp := timestamps[numNodes/2]
	return time.Unix(medianTimestamp, 0)
}

// chainTipEntry defines an entry used to track the chain tips and is structured
// such that there is a single statically-allocated field to house a tip, and a
// dynamically-allocated slice for the rare case when there are multiple
// tips at the same height.
//
// This is done to reduce the number of allocations for the common case since
// there is typically only a single tip at a given height.
type chainTipEntry struct {
	tip       *blockNode
	otherTips []*blockNode
}

// blockIndex provides facilities for keeping track of an in-memory index of the
// block chain.  Although the name block chain suggests a single chain of
// blocks, it is actually a tree-shaped structure where any node can have
// multiple children.  However, there can only be one active branch which does
// indeed form a chain from the tip all the way back to the genesis block.
type blockIndex struct {
	// The following fields are set when the instance is created and can't
	// be changed afterwards, so there is no need to protect them with a
	// separate mutex.
	db database.DB

	// These following fields are protected by the embedded mutex.
	//
	// index contains an entry for every known block tracked by the block
	// index.
	//
	// modified contains an entry for all nodes that have been modified
	// since the last time the index was flushed to disk.
	//
	// chainTips contains an entry with the tip of all known side chains.
	//
	// totalTips tracks the total number of all known chain tips.
	sync.RWMutex
	index     map[chainhash.Hash]*blockNode
	modified  map[*blockNode]struct{}
	chainTips map[int64]chainTipEntry
	totalTips uint64
}

// newBlockIndex returns a new empty instance of a block index.  The index will
// be dynamically populated as block nodes are loaded from the database and
// manually added.
func newBlockIndex(db database.DB) *blockIndex {
	return &blockIndex{
		db:        db,
		index:     make(map[chainhash.Hash]*blockNode),
		modified:  make(map[*blockNode]struct{}),
		chainTips: make(map[int64]chainTipEntry),
	}
}

// HaveBlock returns whether or not the block index contains the provided hash
// and the block data is available.
//
// This function is safe for concurrent access.
func (bi *blockIndex) HaveBlock(hash *chainhash.Hash) bool {
	bi.RLock()
	node := bi.index[*hash]
	hasBlock := node != nil && node.status.HaveData()
	bi.RUnlock()
	return hasBlock
}

// addNode adds the provided node to the block index.  Duplicate entries are not
// checked so it is up to caller to avoid adding them.
//
// This function MUST be called with the block index lock held (for writes).
func (bi *blockIndex) addNode(node *blockNode) {
	bi.index[node.hash] = node

	// Since the block index does not support nodes that do not connect to
	// an existing node (except the genesis block), all new nodes are either
	// extending an existing chain or are on a side chain, but in either
	// case, are a new chain tip.  In the case the node is extending a
	// chain, the parent is no longer a tip.
	bi.addChainTip(node)
	if node.parent != nil {
		bi.removeChainTip(node.parent)
	}
}

// AddNode adds the provided node to the block index and marks it as modified.
// Duplicate entries are not checked so it is up to caller to avoid adding them.
//
// This function is safe for concurrent access.
func (bi *blockIndex) AddNode(node *blockNode) {
	bi.Lock()
	bi.addNode(node)
	bi.modified[node] = struct{}{}
	bi.Unlock()
}

// addChainTip adds the passed block node as a new chain tip.
//
// This function MUST be called with the block index lock held (for writes).
func (bi *blockIndex) addChainTip(tip *blockNode) {
	bi.totalTips++

	// When an entry does not already exist for the given tip height, add an
	// entry to the map with the tip stored in the statically-allocated field.
	entry, ok := bi.chainTips[tip.height]
	if !ok {
		bi.chainTips[tip.height] = chainTipEntry{tip: tip}
		return
	}

	// Otherwise, an entry already exists for the given tip height, so store the
	// tip in the dynamically-allocated slice.
	entry.otherTips = append(entry.otherTips, tip)
	bi.chainTips[tip.height] = entry
}

// removeChainTip removes the passed block node from the available chain tips.
//
// This function MUST be called with the block index lock held (for writes).
func (bi *blockIndex) removeChainTip(tip *blockNode) {
	// Nothing to do if no tips exist at the given height.
	entry, ok := bi.chainTips[tip.height]
	if !ok {
		return
	}

	// The most common case is a single tip at the given height, so handle the
	// case where the tip that is being removed is the tip that is stored in the
	// statically-allocated field first.
	if entry.tip == tip {
		bi.totalTips--
		entry.tip = nil

		// Remove the map entry altogether if there are no more tips left.
		if len(entry.otherTips) == 0 {
			delete(bi.chainTips, tip.height)
			return
		}

		// There are still tips stored in the dynamically-allocated slice, so
		// move the first tip from it to the statically-allocated field, nil the
		// slice so it can be garbage collected when there are no more items in
		// it, and update the map with the modified entry accordingly.
		entry.tip = entry.otherTips[0]
		entry.otherTips = entry.otherTips[1:]
		if len(entry.otherTips) == 0 {
			entry.otherTips = nil
		}
		bi.chainTips[tip.height] = entry
		return
	}

	// The tip being removed is not the tip stored in the statically-allocated
	// field, so attempt to remove it from the dyanimcally-allocated slice.
	for i, n := range entry.otherTips {
		if n == tip {
			bi.totalTips--

			copy(entry.otherTips[i:], entry.otherTips[i+1:])
			entry.otherTips[len(entry.otherTips)-1] = nil
			entry.otherTips = entry.otherTips[:len(entry.otherTips)-1]
			if len(entry.otherTips) == 0 {
				entry.otherTips = nil
			}
			bi.chainTips[tip.height] = entry
			return
		}
	}
}

// lookupNode returns the block node identified by the provided hash.  It will
// return nil if there is no entry for the hash.
//
// This function MUST be called with the block index lock held (for reads).
func (bi *blockIndex) lookupNode(hash *chainhash.Hash) *blockNode {
	return bi.index[*hash]
}

// LookupNode returns the block node identified by the provided hash.  It will
// return nil if there is no entry for the hash.
//
// This function is safe for concurrent access.
func (bi *blockIndex) LookupNode(hash *chainhash.Hash) *blockNode {
	bi.RLock()
	node := bi.lookupNode(hash)
	bi.RUnlock()
	return node
}

// NodeStatus returns the status associated with the provided node.
//
// This function is safe for concurrent access.
func (bi *blockIndex) NodeStatus(node *blockNode) blockStatus {
	bi.RLock()
	status := node.status
	bi.RUnlock()
	return status
}

// SetStatusFlags sets the provided status flags for the given block node
// regardless of their previous state.  It does not unset any flags.
//
// This function is safe for concurrent access.
func (bi *blockIndex) SetStatusFlags(node *blockNode, flags blockStatus) {
	bi.Lock()
	origStatus := node.status
	node.status |= flags
	if node.status != origStatus {
		bi.modified[node] = struct{}{}
	}
	bi.Unlock()
}

// UnsetStatusFlags unsets the provided status flags for the given block node
// regardless of their previous state.
//
// This function is safe for concurrent access.
func (bi *blockIndex) UnsetStatusFlags(node *blockNode, flags blockStatus) {
	bi.Lock()
	origStatus := node.status
	node.status &^= flags
	if node.status != origStatus {
		bi.modified[node] = struct{}{}
	}
	bi.Unlock()
}

// flush writes all of the modified block nodes to the database and clears the
// set of modified nodes if it succeeds.
func (bi *blockIndex) flush() error {
	// Nothing to flush if there are no modified nodes.
	bi.Lock()
	if len(bi.modified) == 0 {
		bi.Unlock()
		return nil
	}

	// Write all of the nodes in the set of modified nodes to the database.
	err := bi.db.Update(func(dbTx database.Tx) error {
		for node := range bi.modified {
			err := dbPutBlockNode(dbTx, node)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		bi.Unlock()
		return err
	}

	// Clear the set of modified nodes.
	bi.modified = make(map[*blockNode]struct{})
	bi.Unlock()
	return nil
}
