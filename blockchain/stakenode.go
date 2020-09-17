// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"fmt"

	"github.com/hdfchain/hdfd/blockchain/stake"
	"github.com/hdfchain/hdfd/chaincfg/chainhash"
	"github.com/hdfchain/hdfd/database"
)

// fetchNewTicketsForNode fetches the list of newly maturing tickets for a
// given node by traversing backwards through its parents until it finds the
// block that contains the original tickets to mature.
//
// This function is NOT safe for concurrent access and must be called with
// the chainLock held for writes.
func (b *BlockChain) fetchNewTicketsForNode(node *blockNode) ([]chainhash.Hash, error) {
	// If we're before the stake enabled height, there can be no
	// tickets in the live ticket pool.
	if node.height < b.chainParams.StakeEnabledHeight {
		return []chainhash.Hash{}, nil
	}

	// If we already cached the tickets, simply return the cached list.
	// It's important to make the distinction here that nil means the
	// value was never looked up, while an empty slice of pointers means
	// that there were no new tickets at this height.
	if node.newTickets != nil {
		return node.newTickets, nil
	}

	// Calculate block number for where new tickets matured from and retrieve
	// this block from DB or in memory if it's a sidechain.
	matureNode := node.RelativeAncestor(int64(b.chainParams.TicketMaturity))
	if matureNode == nil {
		return nil, fmt.Errorf("unable to obtain previous node; " +
			"ancestor is genesis block")
	}

	matureBlock, errBlock := b.fetchBlockByNode(matureNode)
	if errBlock != nil {
		return nil, errBlock
	}

	tickets := []chainhash.Hash{}
	for _, stx := range matureBlock.MsgBlock().STransactions {
		if stake.IsSStx(stx) {
			h := stx.TxHash()
			tickets = append(tickets, h)
		}
	}

	// Set the new tickets in memory so that they exist for future
	// reference in the node.
	node.newTickets = tickets

	return tickets, nil
}

// fetchStakeNode will scour the blockchain from the best block, for which we
// know that there is valid stake node.  The first step is finding a path to the
// ancestor, or, if on a side chain, the path to the common ancestor, followed
// by the path to the sidechain node.  After this path is established, the
// algorithm walks along the path, regenerating and storing intermediate nodes
// as it does so, until the final stake node of interest is populated with the
// correct data.
//
// This function MUST be called with the chain state lock held (for writes).
func (b *BlockChain) fetchStakeNode(node *blockNode) (*stake.Node, error) {
	// If we already have the stake node fetched, returned the cached result.
	// Stake nodes are immutable.
	if node.stakeNode != nil {
		return node.stakeNode, nil
	}

	// If the parent stake node is cached, connect the stake node
	// from there.
	if node.parent != nil {
		if node.stakeNode == nil && node.parent.stakeNode != nil {
			var err error
			if node.newTickets == nil {
				node.newTickets, err = b.fetchNewTicketsForNode(node)
				if err != nil {
					return nil, err
				}
			}

			node.stakeNode, err = node.parent.stakeNode.ConnectNode(
				node.lotteryIV(), node.ticketsVoted, node.ticketsRevoked,
				node.newTickets)
			if err != nil {
				return nil, err
			}

			return node.stakeNode, nil
		}
	}

	// We need to generate a path to the stake node and restore it
	// it through the entire path.  The bestNode stake node must
	// always be filled in, so assume it is safe to begin working
	// backwards from there.
	detachNodes, attachNodes := b.getReorganizeNodes(node)
	current := b.bestChain.Tip()

	// Flush any potential unsaved changes to the block index due to the
	// call to get the reorganize nodes.  Since the best state is not being
	// being modified, it is safe to ignore any errors here as the changes
	// will be flushed at that point and those errors are not ignored.
	b.flushBlockIndexWarnOnly()

	// Move backwards through the main chain, undoing the ticket
	// treaps for each block.  The database is passed because the
	// undo data and new tickets data for each block may not yet
	// be filled in and may require the database to look up.
	err := b.db.View(func(dbTx database.Tx) error {
		for e := detachNodes.Front(); e != nil; e = e.Next() {
			n := e.Value.(*blockNode)
			if n.stakeNode == nil {
				var errLocal error
				n.stakeNode, errLocal =
					current.stakeNode.DisconnectNode(n.lotteryIV(),
						n.stakeUndoData, n.newTickets, dbTx)
				if errLocal != nil {
					return errLocal
				}
			}
			current = n
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Detach the final block and get the filled in node for the fork
	// point.
	err = b.db.View(func(dbTx database.Tx) error {
		if current.parent.stakeNode == nil {
			var errLocal error
			current.parent.stakeNode, errLocal =
				current.stakeNode.DisconnectNode(current.parent.lotteryIV(),
					current.parent.stakeUndoData, current.parent.newTickets, dbTx)
			if errLocal != nil {
				return errLocal
			}
		}
		current = current.parent
		return nil
	})
	if err != nil {
		return nil, err
	}

	// The node is at a fork point in the block chain, so just return
	// this stake node.
	if attachNodes.Len() == 0 {
		if current.hash != node.hash ||
			current.height != node.height {
			return nil, AssertError("failed to restore stake node to " +
				"fork point when fetching")
		}

		return current.stakeNode, nil
	}

	// The requested node is on a side chain, so we need to apply the
	// transactions and spend information from each of the nodes to attach.
	// Not that side chain ticket data and undo data is always stored
	// in memory, so there is not need to use the database here.
	for e := attachNodes.Front(); e != nil; e = e.Next() {
		n := e.Value.(*blockNode)

		if n.stakeNode == nil {
			if n.newTickets == nil {
				n.newTickets, err = b.fetchNewTicketsForNode(n)
				if err != nil {
					return nil, err
				}
			}

			n.stakeNode, err = current.stakeNode.ConnectNode(n.lotteryIV(),
				n.ticketsVoted, n.ticketsRevoked, n.newTickets)
			if err != nil {
				return nil, err
			}
		}

		current = n
	}

	return current.stakeNode, nil
}
