// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2020 The Hdfchain developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"github.com/hdfchain/hdfd/blockchain/stake/v3"
	"github.com/hdfchain/hdfd/blockchain/standalone/v2"
	"github.com/hdfchain/hdfd/blockchain/v3/internal/progresslog"
	"github.com/hdfchain/hdfd/chaincfg/chainhash"
	"github.com/hdfchain/hdfd/chaincfg/v3"
	"github.com/hdfchain/hdfd/database/v2"
	"github.com/hdfchain/hdfd/dcrutil/v3"
	"github.com/hdfchain/hdfd/gcs/v2"
	"github.com/hdfchain/hdfd/gcs/v2/blockcf2"
	"github.com/hdfchain/hdfd/wire"
)

// errInterruptRequested indicates that an operation was cancelled due
// to a user-requested interrupt.
var errInterruptRequested = errors.New("interrupt requested")

// errBatchFinished indicates that a foreach database loop was exited due to
// reaching the maximum batch size.
var errBatchFinished = errors.New("batch finished")

// interruptRequested returns true when the provided channel has been closed.
// This simplifies early shutdown slightly since the caller can just use an if
// statement instead of a select.
func interruptRequested(ctx context.Context) bool {
	return ctx.Err() != nil
}

// deserializeDatabaseInfoV2 deserializes a database information struct from the
// passed serialized byte slice according to the legacy version 2 format.
//
// The legacy format is as follows:
//
//   Field      Type     Size      Description
//   version    uint32   4 bytes   The version of the database
//   compVer    uint32   4 bytes   The script compression version of the database
//   created    uint32   4 bytes   The date of the creation of the database
//
// The high bit (0x80000000) is used on version to indicate that an upgrade
// is in progress and used to confirm the database fidelity on start up.
func deserializeDatabaseInfoV2(dbInfoBytes []byte) (*databaseInfo, error) {
	// upgradeStartedBit if the bit flag for whether or not a database
	// upgrade is in progress. It is used to determine if the database
	// is in an inconsistent state from the update.
	const upgradeStartedBit = 0x80000000

	byteOrder := binary.LittleEndian

	rawVersion := byteOrder.Uint32(dbInfoBytes[0:4])
	upgradeStarted := (upgradeStartedBit & rawVersion) > 0
	version := rawVersion &^ upgradeStartedBit
	compVer := byteOrder.Uint32(dbInfoBytes[4:8])
	ts := byteOrder.Uint32(dbInfoBytes[8:12])

	if upgradeStarted {
		return nil, AssertError("database is in the upgrade started " +
			"state before resumable upgrades were supported - " +
			"delete the database and resync the blockchain")
	}

	return &databaseInfo{
		version: version,
		compVer: compVer,
		created: time.Unix(int64(ts), 0),
	}, nil
}

// ticketsVotedInBlock fetches a list of tickets that were voted in the
// block.
func ticketsVotedInBlock(bl *dcrutil.Block) []chainhash.Hash {
	var tickets []chainhash.Hash
	for _, stx := range bl.MsgBlock().STransactions {
		if stake.IsSSGen(stx) {
			tickets = append(tickets, stx.TxIn[1].PreviousOutPoint.Hash)
		}
	}

	return tickets
}

// ticketsRevokedInBlock fetches a list of tickets that were revoked in the
// block.
func ticketsRevokedInBlock(bl *dcrutil.Block) []chainhash.Hash {
	var tickets []chainhash.Hash
	for _, stx := range bl.MsgBlock().STransactions {
		if stake.DetermineTxType(stx) == stake.TxTypeSSRtx {
			tickets = append(tickets, stx.TxIn[0].PreviousOutPoint.Hash)
		}
	}

	return tickets
}

// upgradeToVersion2 upgrades a version 1 blockchain to version 2, allowing
// use of the new on-disk ticket database.
func upgradeToVersion2(db database.DB, chainParams *chaincfg.Params, dbInfo *databaseInfo) error {
	// Hardcoded so updates to the global values do not affect old upgrades.
	byteOrder := binary.LittleEndian
	chainStateKeyName := []byte("chainstate")
	heightIdxBucketName := []byte("heightidx")

	// These are legacy functions that relied on information in the database
	// that is no longer available in more recent code.
	dbFetchHashByHeight := func(dbTx database.Tx, height int64) (*chainhash.Hash, error) {
		var serializedHeight [4]byte
		byteOrder.PutUint32(serializedHeight[:], uint32(height))

		meta := dbTx.Metadata()
		heightIndex := meta.Bucket(heightIdxBucketName)
		hashBytes := heightIndex.Get(serializedHeight[:])
		if hashBytes == nil {
			str := fmt.Sprintf("no block at height %d exists", height)
			return nil, errNotInMainChain(str)
		}

		var hash chainhash.Hash
		copy(hash[:], hashBytes)
		return &hash, nil
	}
	dbFetchBlockByHeight := func(dbTx database.Tx, height int64) (*dcrutil.Block, error) {
		// First find the hash associated with the provided height in the index.
		hash, err := dbFetchHashByHeight(dbTx, height)
		if err != nil {
			return nil, err
		}

		// Load the raw block bytes from the database.
		blockBytes, err := dbTx.FetchBlock(hash)
		if err != nil {
			return nil, err
		}

		// Create the encapsulated block and set the height appropriately.
		block, err := dcrutil.NewBlockFromBytes(blockBytes)
		if err != nil {
			return nil, err
		}

		return block, nil
	}

	log.Info("Initializing upgrade to database version 2")
	progressLogger := progresslog.NewBlockProgressLogger("Upgraded", log)

	// The upgrade is atomic, so there is no need to set the flag that
	// the database is undergoing an upgrade here.  Get the stake node
	// for the genesis block, and then begin connecting stake nodes
	// incrementally.
	err := db.Update(func(dbTx database.Tx) error {
		// Fetch the stored best chain state from the database metadata.
		serializedData := dbTx.Metadata().Get(chainStateKeyName)
		best, err := deserializeBestChainState(serializedData)
		if err != nil {
			return err
		}

		bestStakeNode, errLocal := stake.InitDatabaseState(dbTx, chainParams,
			&chainParams.GenesisHash)
		if errLocal != nil {
			return errLocal
		}

		parent, errLocal := dbFetchBlockByHeight(dbTx, 0)
		if errLocal != nil {
			return errLocal
		}

		for i := int64(1); i <= int64(best.height); i++ {
			block, errLocal := dbFetchBlockByHeight(dbTx, i)
			if errLocal != nil {
				return errLocal
			}

			// If we need the tickets, fetch them too.
			var newTickets []chainhash.Hash
			if i >= chainParams.StakeEnabledHeight {
				matureHeight := i - int64(chainParams.TicketMaturity)
				matureBlock, errLocal := dbFetchBlockByHeight(dbTx, matureHeight)
				if errLocal != nil {
					return errLocal
				}
				for _, stx := range matureBlock.MsgBlock().STransactions {
					if stake.IsSStx(stx) {
						h := stx.TxHash()
						newTickets = append(newTickets, h)
					}
				}
			}

			// Iteratively connect the stake nodes in memory.
			header := block.MsgBlock().Header
			hB, errLocal := header.Bytes()
			if errLocal != nil {
				return errLocal
			}
			bestStakeNode, errLocal = bestStakeNode.ConnectNode(
				stake.CalcHash256PRNGIV(hB), ticketsVotedInBlock(block),
				ticketsRevokedInBlock(block), newTickets)
			if errLocal != nil {
				return errLocal
			}

			// Write the top block stake node to the database.
			errLocal = stake.WriteConnectedBestNode(dbTx, bestStakeNode,
				best.hash)
			if errLocal != nil {
				return errLocal
			}

			progressLogger.LogBlockHeight(block.MsgBlock(), parent.MsgBlock())
			parent = block
		}

		// Write the new database version.
		dbInfo.version = 2
		return dbPutDatabaseInfo(dbTx, dbInfo)
	})
	if err != nil {
		return err
	}

	log.Info("Upgrade to new stake database was successful!")

	return nil
}

// -----------------------------------------------------------------------------
// The legacy version 2 block index consists of an entry for every known block.
// which includes information such as the block header and hashes of tickets
// voted and revoked.
//
// The serialized key format is:
//
//   <block height><block hash>
//
//   Field           Type              Size
//   block height    uint32            4 bytes
//   block hash      chainhash.Hash    chainhash.HashSize
//
// The serialized value format is:
//
//   <block header><status><num votes><votes info><num revoked><revoked tickets>
//
//   Field              Type                Size
//   block header       wire.BlockHeader    180 bytes
//   status             blockStatus         1 byte
//   num votes          VLQ                 variable
//   vote info
//     ticket hash      chainhash.Hash      chainhash.HashSize
//     vote version     VLQ                 variable
//     vote bits        VLQ                 variable
//   num revoked        VLQ                 variable
//   revoked tickets
//     ticket hash      chainhash.Hash      chainhash.HashSize
// -----------------------------------------------------------------------------
type blockIndexEntryV2 struct {
	header         wire.BlockHeader
	status         blockStatus
	voteInfo       []stake.VoteVersionTuple
	ticketsVoted   []chainhash.Hash
	ticketsRevoked []chainhash.Hash
}

// blockIndexEntrySerializeSizV2e returns the number of bytes it would take to
// serialize the passed block index entry according to the legacy version 2
// format described above.
func blockIndexEntrySerializeSizeV2(entry *blockIndexEntryV2) int {
	voteInfoSize := 0
	for i := range entry.voteInfo {
		voteInfoSize += chainhash.HashSize +
			serializeSizeVLQ(uint64(entry.voteInfo[i].Version)) +
			serializeSizeVLQ(uint64(entry.voteInfo[i].Bits))
	}

	return blockHdrSize + 1 + serializeSizeVLQ(uint64(len(entry.voteInfo))) +
		voteInfoSize + serializeSizeVLQ(uint64(len(entry.ticketsRevoked))) +
		chainhash.HashSize*len(entry.ticketsRevoked)
}

// putBlockIndexEntryV2 serializes the passed block index entry according to the
// legacy version 2 format described above directly into the passed target byte
// slice.  The target byte slice must be at least large enough to handle the
// number of bytes returned by the blockIndexEntrySerializeSizeV2 function or it
// will panic.
func putBlockIndexEntryV2(target []byte, entry *blockIndexEntryV2) (int, error) {
	if len(entry.voteInfo) != len(entry.ticketsVoted) {
		return 0, AssertError("putBlockIndexEntry called with " +
			"mismatched number of tickets voted and vote info")
	}

	// Serialize the entire block header.
	w := bytes.NewBuffer(target[0:0])
	if err := entry.header.Serialize(w); err != nil {
		return 0, err
	}

	// Serialize the status.
	offset := blockHdrSize
	target[offset] = byte(entry.status)
	offset++

	// Serialize the number of votes and associated vote information.
	offset += putVLQ(target[offset:], uint64(len(entry.voteInfo)))
	for i := range entry.voteInfo {
		offset += copy(target[offset:], entry.ticketsVoted[i][:])
		offset += putVLQ(target[offset:], uint64(entry.voteInfo[i].Version))
		offset += putVLQ(target[offset:], uint64(entry.voteInfo[i].Bits))
	}

	// Serialize the number of revocations and associated revocation
	// information.
	offset += putVLQ(target[offset:], uint64(len(entry.ticketsRevoked)))
	for i := range entry.ticketsRevoked {
		offset += copy(target[offset:], entry.ticketsRevoked[i][:])
	}

	return offset, nil
}

// decodeBlockIndexEntryV2 decodes the passed serialized block index entry into
// the passed struct according to the legacy version 2 format described above.
// It returns the number of bytes read.
func decodeBlockIndexEntryV2(serialized []byte, entry *blockIndexEntryV2) (int, error) {
	// Hardcoded value so updates do not affect old upgrades.
	const blockHdrSize = 180

	// Ensure there are enough bytes to decode header.
	if len(serialized) < blockHdrSize {
		return 0, errDeserialize("unexpected end of data while reading block " +
			"header")
	}
	hB := serialized[0:blockHdrSize]

	// Deserialize the header.
	var header wire.BlockHeader
	if err := header.Deserialize(bytes.NewReader(hB)); err != nil {
		return 0, err
	}
	offset := blockHdrSize

	// Deserialize the status.
	if offset+1 > len(serialized) {
		return offset, errDeserialize("unexpected end of data while reading " +
			"status")
	}
	status := blockStatus(serialized[offset])
	offset++

	// Deserialize the number of tickets spent.
	var ticketsVoted []chainhash.Hash
	var votes []stake.VoteVersionTuple
	numVotes, bytesRead := deserializeVLQ(serialized[offset:])
	if bytesRead == 0 {
		return offset, errDeserialize("unexpected end of data while reading " +
			"num votes")
	}
	offset += bytesRead
	if numVotes > 0 {
		ticketsVoted = make([]chainhash.Hash, numVotes)
		votes = make([]stake.VoteVersionTuple, numVotes)
		for i := uint64(0); i < numVotes; i++ {
			// Deserialize the ticket hash associated with the vote.
			if offset+chainhash.HashSize > len(serialized) {
				return offset, errDeserialize(fmt.Sprintf("unexpected end of "+
					"data while reading vote #%d hash", i))
			}
			copy(ticketsVoted[i][:], serialized[offset:])
			offset += chainhash.HashSize

			// Deserialize the vote version.
			version, bytesRead := deserializeVLQ(serialized[offset:])
			if bytesRead == 0 {
				return offset, errDeserialize(fmt.Sprintf("unexpected end of "+
					"data while reading vote #%d version", i))
			}
			offset += bytesRead

			// Deserialize the vote bits.
			voteBits, bytesRead := deserializeVLQ(serialized[offset:])
			if bytesRead == 0 {
				return offset, errDeserialize(fmt.Sprintf("unexpected end of "+
					"data while reading vote #%d bits", i))
			}
			offset += bytesRead

			votes[i].Version = uint32(version)
			votes[i].Bits = uint16(voteBits)
		}
	}

	// Deserialize the number of tickets revoked.
	var ticketsRevoked []chainhash.Hash
	numTicketsRevoked, bytesRead := deserializeVLQ(serialized[offset:])
	if bytesRead == 0 {
		return offset, errDeserialize("unexpected end of data while reading " +
			"num tickets revoked")
	}
	offset += bytesRead
	if numTicketsRevoked > 0 {
		ticketsRevoked = make([]chainhash.Hash, numTicketsRevoked)
		for i := uint64(0); i < numTicketsRevoked; i++ {
			// Deserialize the ticket hash associated with the
			// revocation.
			if offset+chainhash.HashSize > len(serialized) {
				return offset, errDeserialize(fmt.Sprintf("unexpected end of "+
					"data while reading revocation #%d", i))
			}
			copy(ticketsRevoked[i][:], serialized[offset:])
			offset += chainhash.HashSize
		}
	}

	entry.header = header
	entry.status = status
	entry.voteInfo = votes
	entry.ticketsVoted = ticketsVoted
	entry.ticketsRevoked = ticketsRevoked
	return offset, nil
}

// migrateBlockIndex migrates all block entries from the v1 block index bucket
// managed by ffldb to the v2 bucket managed by this package.  The v1 bucket
// stored all block entries keyed by block hash, whereas the v2 bucket stores
// them keyed by block height + hash.  Also, the old block index only stored the
// header, while the new one stores all info needed to recreate block nodes.
//
// The new block index is guaranteed to be fully updated if this returns without
// failure.
func migrateBlockIndex(ctx context.Context, db database.DB) error {
	// blkHdrOffset defines the offsets into a v1 block index row for the block
	// header.
	//
	// The serialized block index row format is:
	//   <blocklocation><blockheader>
	const blkHdrOffset = 12

	// blkHdrHeightStart is the offset of the height in the serialized block
	// header bytes as it existed at the time of this migration.  It is hard
	// coded here so potential future changes do not affect old upgrades.
	const blkHdrHeightStart = 128

	// Hardcoded bucket names so updates to the global values do not affect old
	// upgrades.
	v1BucketName := []byte("ffldb-blockidx")
	v2BucketName := []byte("blockidx")
	hashIdxBucketName := []byte("hashidx")

	log.Info("Reindexing block information in the database.  This will take " +
		"a while...")
	start := time.Now()

	// Create the new block index bucket as needed.
	err := db.Update(func(dbTx database.Tx) error {
		_, err := dbTx.Metadata().CreateBucketIfNotExists(v2BucketName)
		return err
	})
	if err != nil {
		return err
	}

	// doBatch contains the primary logic for upgrading the block index from
	// version 1 to 2 in batches.  This is done because attempting to migrate in
	// a single database transaction could result in massive memory usage and
	// could potentially crash on many systems due to ulimits.
	//
	// It returns the number of entries processed.
	const maxEntries = 20000
	var resumeOffset uint32
	doBatch := func(dbTx database.Tx) (uint32, error) {
		meta := dbTx.Metadata()
		v1BlockIdxBucket := meta.Bucket(v1BucketName)
		if v1BlockIdxBucket == nil {
			return 0, fmt.Errorf("bucket %s does not exist", v1BucketName)
		}

		v2BlockIdxBucket := meta.Bucket(v2BucketName)
		if v2BlockIdxBucket == nil {
			return 0, fmt.Errorf("bucket %s does not exist", v2BucketName)
		}

		hashIdxBucket := meta.Bucket(hashIdxBucketName)
		if hashIdxBucket == nil {
			return 0, fmt.Errorf("bucket %s does not exist", hashIdxBucketName)
		}

		// Migrate block index entries so long as the max number of entries for
		// this batch has not been exceeded.
		var numMigrated, numIterated uint32
		err := v1BlockIdxBucket.ForEach(func(hashBytes, blockRow []byte) error {
			if numMigrated >= maxEntries {
				return errBatchFinished
			}

			// Skip entries that have already been migrated in previous batches.
			numIterated++
			if numIterated-1 < resumeOffset {
				return nil
			}
			resumeOffset++

			// Skip entries that have already been migrated in previous
			// interrupted upgrades.
			var blockHash chainhash.Hash
			copy(blockHash[:], hashBytes)
			endOffset := blkHdrOffset + blockHdrSize
			headerBytes := blockRow[blkHdrOffset:endOffset:endOffset]
			heightBytes := headerBytes[blkHdrHeightStart : blkHdrHeightStart+4]
			height := binary.LittleEndian.Uint32(heightBytes)
			key := blockIndexKey(&blockHash, height)
			if v2BlockIdxBucket.Get(key) != nil {
				return nil
			}

			// Load the raw full block from the database.
			blockBytes, err := dbTx.FetchBlock(&blockHash)
			if err != nil {
				return err
			}

			// Deserialize the block bytes.
			var block wire.MsgBlock
			err = block.Deserialize(bytes.NewReader(blockBytes))
			if err != nil {
				return err
			}

			// Mark the block as valid if it's part of the main chain.  While it
			// is possible side chain blocks were validated too, there was
			// previously no tracking of that information, so there is no way to
			// know for sure.  It's better to be safe and just assume side chain
			// blocks were never validated.
			status := statusDataStored
			if hashIdxBucket.Get(blockHash[:]) != nil {
				status |= statusValidated
			}

			// Write the serialized block index entry to the new bucket keyed by
			// its hash and height.
			ticketInfo := stake.FindSpentTicketsInBlock(&block)
			entry := &blockIndexEntryV2{
				header:         block.Header,
				status:         status,
				voteInfo:       ticketInfo.Votes,
				ticketsVoted:   ticketInfo.VotedTickets,
				ticketsRevoked: ticketInfo.RevokedTickets,
			}
			serialized := make([]byte, blockIndexEntrySerializeSizeV2(entry))
			if _, err = putBlockIndexEntryV2(serialized, entry); err != nil {
				return err
			}
			err = v2BlockIdxBucket.Put(key, serialized)
			if err != nil {
				return err
			}

			numMigrated++

			if interruptRequested(ctx) {
				return errInterruptRequested
			}

			return nil
		})
		return numMigrated, err
	}

	// Migrate all entries in batches for the reasons mentioned above.
	var totalMigrated uint64
	for {
		var numMigrated uint32
		err := db.Update(func(dbTx database.Tx) error {
			var err error
			numMigrated, err = doBatch(dbTx)
			if errors.Is(err, errInterruptRequested) ||
				errors.Is(err, errBatchFinished) {
				// No error here so the database transaction is
				// not cancelled and therefore outstanding work
				// is written to disk.  The outer function will
				// exit with an interrupted error below due to
				// another interrupted check.
				err = nil
			}
			return err
		})
		if err != nil {
			return err
		}

		if interruptRequested(ctx) {
			return errInterruptRequested
		}

		if numMigrated == 0 {
			break
		}

		totalMigrated += uint64(numMigrated)
		log.Infof("Migrated %d entries (%d total)", numMigrated, totalMigrated)
	}

	seconds := int64(time.Since(start) / time.Second)
	log.Infof("Done upgrading block index.  Total entries: %d in %d seconds",
		totalMigrated, seconds)
	return nil
}

// upgradeToVersion3 upgrades a version 2 blockchain to version 3 along with
// upgrading the block index to version 2.
func upgradeToVersion3(ctx context.Context, db database.DB, dbInfo *databaseInfo) error {
	if err := migrateBlockIndex(ctx, db); err != nil {
		return err
	}

	// Update and persist the updated database versions.
	dbInfo.version = 3
	dbInfo.bidxVer = 2
	return db.Update(func(dbTx database.Tx) error {
		return dbPutDatabaseInfo(dbTx, dbInfo)
	})
}

// removeMainChainIndex removes the main chain hash index and height index
// buckets.  These are no longer needed due to using the full block index in
// memory.
//
// The database is guaranteed to be fully updated if this returns without
// failure.
func removeMainChainIndex(ctx context.Context, db database.DB) error {
	// Hardcoded bucket names so updates to the global values do not affect old
	// upgrades.
	hashIdxBucketName := []byte("hashidx")
	heightIdxBucketName := []byte("heightidx")

	log.Info("Removing unneeded indexes in the database...")
	start := time.Now()

	// Delete the main chain index buckets.
	err := db.Update(func(dbTx database.Tx) error {
		// Delete the main chain hash to height index.
		meta := dbTx.Metadata()
		hashIdxBucket := meta.Bucket(hashIdxBucketName)
		if hashIdxBucket != nil {
			if err := meta.DeleteBucket(hashIdxBucketName); err != nil {
				return err
			}
			log.Info("Removed hash index.")
		}

		if interruptRequested(ctx) {
			// No error here so the database transaction is not cancelled
			// and therefore outstanding work is written to disk.  The
			// outer function will exit with an interrupted error below due
			// to another interrupted check.
			return nil
		}

		// Delete the main chain hash to height index.
		heightIdxBucket := meta.Bucket(heightIdxBucketName)
		if heightIdxBucket != nil {
			if err := meta.DeleteBucket(heightIdxBucketName); err != nil {
				return err
			}
			log.Info("Removed height index.")
		}

		return nil
	})
	if err != nil {
		return err
	}

	if interruptRequested(ctx) {
		return errInterruptRequested
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	log.Infof("Done upgrading database in %v.", elapsed)
	return nil
}

// upgradeToVersion4 upgrades a version 3 blockchain database to version 4.
func upgradeToVersion4(ctx context.Context, db database.DB, dbInfo *databaseInfo) error {
	if err := removeMainChainIndex(ctx, db); err != nil {
		return err
	}

	// Update and persist the updated database versions.
	dbInfo.version = 4
	return db.Update(func(dbTx database.Tx) error {
		return dbPutDatabaseInfo(dbTx, dbInfo)
	})
}

// incrementalFlatDrop uses multiple database updates to remove key/value pairs
// saved to a flag bucket.
func incrementalFlatDrop(ctx context.Context, db database.DB, bucketKey []byte, humanName string) error {
	const maxDeletions = 2000000
	var totalDeleted uint64
	for numDeleted := maxDeletions; numDeleted == maxDeletions; {
		numDeleted = 0
		err := db.Update(func(dbTx database.Tx) error {
			bucket := dbTx.Metadata().Bucket(bucketKey)
			cursor := bucket.Cursor()
			for ok := cursor.First(); ok; ok = cursor.Next() &&
				numDeleted < maxDeletions {

				if err := cursor.Delete(); err != nil {
					return err
				}
				numDeleted++
			}
			return nil
		})
		if err != nil {
			return err
		}

		if numDeleted > 0 {
			totalDeleted += uint64(numDeleted)
			log.Infof("Deleted %d keys (%d total) from %s", numDeleted,
				totalDeleted, humanName)
		}

		if interruptRequested(ctx) {
			return errInterruptRequested
		}
	}
	return nil
}

// upgradeToVersion5 upgrades a version 4 blockchain database to version 5.
func upgradeToVersion5(ctx context.Context, db database.DB, chainParams *chaincfg.Params, dbInfo *databaseInfo) error {
	// Hardcoded bucket and key names so updates to the global values do not
	// affect old upgrades.
	utxoSetBucketName := []byte("utxoset")
	spendJournalBucketName := []byte("spendjournal")
	chainStateKeyName := []byte("chainstate")
	v5ReindexTipKeyName := []byte("v5reindextip")

	log.Info("Clearing database utxoset and spend journal for upgrade...")
	start := time.Now()

	// Clear the utxoset.
	err := incrementalFlatDrop(ctx, db, utxoSetBucketName, "utxoset")
	if err != nil {
		return err
	}
	log.Info("Cleared utxoset.")

	if interruptRequested(ctx) {
		return errInterruptRequested
	}

	// Clear the spend journal.
	err = incrementalFlatDrop(ctx, db, spendJournalBucketName, "spend journal")
	if err != nil {
		return err
	}
	log.Info("Cleared spend journal.")

	if interruptRequested(ctx) {
		return errInterruptRequested
	}

	err = db.Update(func(dbTx database.Tx) error {
		// Reset the ticket database to the genesis block.
		log.Info("Resetting the ticket database.  This might take a while...")
		err := stake.ResetDatabase(dbTx, chainParams, &chainParams.GenesisHash)
		if err != nil {
			return err
		}

		// Fetch the stored best chain state from the database metadata.
		meta := dbTx.Metadata()
		serializedData := meta.Get(chainStateKeyName)
		best, err := deserializeBestChainState(serializedData)
		if err != nil {
			return err
		}

		// Store the current best chain tip as the reindex target.
		if err := meta.Put(v5ReindexTipKeyName, best.hash[:]); err != nil {
			return err
		}

		// Reset the state related to the best block to the genesis block.
		genesisBlock := chainParams.GenesisBlock
		numTxns := uint64(len(genesisBlock.Transactions))
		serializedData = serializeBestChainState(bestChainState{
			hash:         genesisBlock.BlockHash(),
			height:       0,
			totalTxns:    numTxns,
			totalSubsidy: 0,
			workSum:      standalone.CalcWork(genesisBlock.Header.Bits),
		})
		err = meta.Put(chainStateKeyName, serializedData)
		if err != nil {
			return err
		}

		// Update and persist the updated database versions.
		dbInfo.version = 5
		return dbPutDatabaseInfo(dbTx, dbInfo)
	})
	if err != nil {
		return err
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	log.Infof("Done upgrading database in %v.", elapsed)
	return nil
}

// maybeFinishV5Upgrade potentially reindexes the chain due to a version 5
// database upgrade.  It will resume previously uncompleted attempts.
func (b *BlockChain) maybeFinishV5Upgrade(ctx context.Context) error {
	// Nothing to do if the database is not version 5.
	if b.dbInfo.version != 5 {
		return nil
	}

	// Hardcoded key name so updates to the global values do not affect old
	// upgrades.
	v5ReindexTipKeyName := []byte("v5reindextip")

	// Finish the version 5 reindex as needed.
	var v5ReindexTipHash *chainhash.Hash
	err := b.db.View(func(dbTx database.Tx) error {
		hash := dbTx.Metadata().Get(v5ReindexTipKeyName)
		if hash != nil {
			v5ReindexTipHash = new(chainhash.Hash)
			copy(v5ReindexTipHash[:], hash)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if v5ReindexTipHash != nil {
		// Look up the final target tip to reindex to in the block index.
		targetTip := b.index.LookupNode(v5ReindexTipHash)
		if targetTip == nil {
			return AssertError(fmt.Sprintf("maybeFinishV5Upgrade: cannot find "+
				"chain tip %s in block index", v5ReindexTipHash))
		}

		// Ensure all ancestors of the current best chain tip are marked as
		// valid.  This is necessary due to older software versions not marking
		// nodes before the final checkpoint as valid.
		for node := targetTip; node != nil; node = node.parent {
			b.index.SetStatusFlags(node, statusValidated)
		}
		if err := b.index.flush(); err != nil {
			return err
		}

		// Disable notifications during the reindex.
		ntfnCallback := b.notifications
		b.notifications = nil
		defer func() {
			b.notifications = ntfnCallback
		}()

		tip := b.bestChain.Tip()
		for tip != targetTip {
			if interruptRequested(ctx) {
				return errInterruptRequested
			}

			// Limit to a reasonable number of blocks at a time.
			const maxReindexBlocks = 250
			intermediateTip := targetTip
			if intermediateTip.height-tip.height > maxReindexBlocks {
				intermediateTip = intermediateTip.Ancestor(tip.height +
					maxReindexBlocks)
			}

			log.Infof("Reindexing to height %d of %d (progress %.2f%%)...",
				intermediateTip.height, targetTip.height,
				float64(intermediateTip.height)/float64(targetTip.height)*100)
			b.chainLock.Lock()
			if err := b.reorganizeChainInternal(intermediateTip); err != nil {
				b.chainLock.Unlock()
				return err
			}
			b.chainLock.Unlock()

			tip = b.bestChain.Tip()
		}

		// Mark the v5 reindex as complete by removing the associated key.
		err := b.db.Update(func(dbTx database.Tx) error {
			return dbTx.Metadata().Delete(v5ReindexTipKeyName)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// scriptSourceEntry houses a script and its associated version.
type scriptSourceEntry struct {
	version uint16
	script  []byte
}

// scriptSource provides a source of transaction output scripts and their
// associated script version for given outpoints and implements the PrevScripter
// interface so it may be used in cases that require access to said scripts.
type scriptSource map[wire.OutPoint]scriptSourceEntry

// PrevScript returns the script and script version associated with the provided
// previous outpoint along with a bool that indicates whether or not the
// requested entry exists.  This ensures the caller is able to distinguish
// between missing entry and empty v0 scripts.
func (s scriptSource) PrevScript(prevOut *wire.OutPoint) (uint16, []byte, bool) {
	entry, ok := s[*prevOut]
	if !ok {
		return 0, nil, false
	}
	return entry.version, entry.script, true
}

// clearFailedBlockFlags unmarks all blocks previously marked failed so they are
// eligible for validation again under new consensus rules.  This ensures
// clients that did not update prior to new rules activating are able to
// automatically recover under the new rules without having to download the
// entire chain again.
func clearFailedBlockFlags(index *blockIndex) error {
	for _, node := range index.index {
		index.UnsetStatusFlags(node, statusValidateFailed|statusInvalidAncestor)
	}

	return index.flush()
}

// initializeGCSFilters creates and stores version 2 GCS filters for all blocks
// in the main chain.  This ensures they are immediately available to clients
// and simplifies the rest of the related code since it can rely on the filters
// being available once the upgrade completes.
//
// The database is  guaranteed to have a filter entry for every block in the
// main chain if this returns without failure.
func initializeGCSFilters(ctx context.Context, db database.DB, index *blockIndex, bestChain *chainView) error {
	// Hardcoded values so updates to the global values do not affect old
	// upgrades.
	gcsBucketName := []byte("gcsfilters")
	const compressionVersion = 1

	log.Info("Creating and storing GCS filters.  This will take a while...")
	start := time.Now()

	// Create the new filter bucket as needed.
	err := db.Update(func(dbTx database.Tx) error {
		_, err := dbTx.Metadata().CreateBucketIfNotExists(gcsBucketName)
		return err
	})
	if err != nil {
		return err
	}

	// newFilter loads the full block for the provided node from the db along
	// with its spend journal information and uses it to create a v2 GCS filter.
	newFilter := func(dbTx database.Tx, node *blockNode) (*gcs.FilterV2, error) {
		// Load the full block from the database.
		block, err := dbFetchBlockByNode(dbTx, node)
		if err != nil {
			return nil, err
		}

		// Load all of the spent transaction output data from the database.
		stxos, err := dbFetchSpendJournalEntry(dbTx, block)
		if err != nil {
			return nil, err
		}

		// Use the combination of the block and the stxos to create a source
		// of previous scripts spent by the block needed to create the
		// filter.
		prevScripts := stxosToScriptSource(block, stxos, compressionVersion)

		// Create the filter from the block and referenced previous output
		// scripts.
		filter, err := blockcf2.Regular(block.MsgBlock(), prevScripts)
		if err != nil {
			return nil, err
		}

		return filter, nil
	}

	// doBatch contains the primary logic for creating the GCS filters when
	// moving from database version 5 to 6 in batches.  This is done because
	// attempting to create them all in a single database transaction could
	// result in massive memory usage and could potentially crash on many
	// systems due to ulimits.
	//
	// It returns the number of entries processed as well as the total number
	// bytes occupied by all of the processed filters.
	const maxEntries = 20000
	node := bestChain.Genesis()
	doBatch := func(dbTx database.Tx) (uint64, uint64, error) {
		filterBucket := dbTx.Metadata().Bucket(gcsBucketName)
		if filterBucket == nil {
			return 0, 0, fmt.Errorf("bucket %s does not exist", gcsBucketName)
		}

		var numCreated, totalBytes uint64
		for ; node != nil; node = bestChain.Next(node) {
			if numCreated >= maxEntries {
				break
			}

			// Create the filter from the block and referenced previous output
			// scripts.
			filter, err := newFilter(dbTx, node)
			if err != nil {
				return numCreated, totalBytes, err
			}

			// Store the filter to the database.
			serialized := filter.Bytes()
			err = filterBucket.Put(node.hash[:], serialized)
			if err != nil {
				return numCreated, totalBytes, err
			}
			totalBytes += uint64(len(serialized))

			numCreated++

			if interruptRequested(ctx) {
				return numCreated, totalBytes, errInterruptRequested
			}
		}

		return numCreated, totalBytes, nil
	}

	// Migrate all entries in batches for the reasons mentioned above.
	var totalCreated, totalFilterBytes uint64
	for {
		var numCreated, numFilterBytes uint64
		err := db.Update(func(dbTx database.Tx) error {
			var err error
			numCreated, numFilterBytes, err = doBatch(dbTx)
			if errors.Is(err, errInterruptRequested) {
				// No error here so the database transaction is not cancelled
				// and therefore outstanding work is written to disk.  The
				// outer function will exit with an interrupted error below due
				// to another interrupted check.
				err = nil
			}
			return err
		})
		if err != nil {
			return err
		}

		if interruptRequested(ctx) {
			return errInterruptRequested
		}

		if numCreated == 0 {
			break
		}

		totalCreated += numCreated
		totalFilterBytes += numFilterBytes
		log.Infof("Created %d entries (%d total)", numCreated, totalCreated)
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	log.Infof("Done creating GCS filters in %v.  Total entries: %d (%d bytes)",
		elapsed, totalCreated, totalFilterBytes)
	return nil
}

// upgradeToVersion6 upgrades a version 5 blockchain database to version 6.
func upgradeToVersion6(ctx context.Context, db database.DB, chainParams *chaincfg.Params, dbInfo *databaseInfo) error {
	if interruptRequested(ctx) {
		return errInterruptRequested
	}

	log.Info("Upgrading database to version 6...")
	start := time.Now()

	// Load the chain state and block index from the database.
	bestChain := newChainView(nil)
	index := newBlockIndex(db)
	err := db.View(func(dbTx database.Tx) error {
		// Fetch the stored best chain state from the database.
		state, err := dbFetchBestState(dbTx)
		if err != nil {
			return err
		}

		// Load all of the block index entries from the database and
		// construct the block index.
		err = loadBlockIndex(dbTx, &chainParams.GenesisHash, index)
		if err != nil {
			return err
		}

		// Set the best chain to the stored best state.
		tip := index.lookupNode(&state.hash)
		if tip == nil {
			return AssertError(fmt.Sprintf("initChainState: cannot find "+
				"chain tip %s in block index", state.hash))
		}
		bestChain.SetTip(tip)

		return nil
	})
	if err != nil {
		return err
	}

	// Unmark all blocks previously marked failed so they are eligible for
	// validation again under the new consensus rules.
	if err := clearFailedBlockFlags(index); err != nil {
		return err
	}

	// Create and store version 2 GCS filters for all blocks in the main chain.
	err = initializeGCSFilters(ctx, db, index, bestChain)
	if err != nil {
		return err
	}

	err = db.Update(func(dbTx database.Tx) error {
		// Update and persist the updated database versions.
		dbInfo.version = 6
		return dbPutDatabaseInfo(dbTx, dbInfo)
	})
	if err != nil {
		return err
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	log.Infof("Done upgrading database in %v.", elapsed)
	return nil
}

// migrateBlockIndexVersion2To3 migrates all block entries from the v2 block
// index bucket to a v3 bucket and removes the old v2 bucket.  As compared to
// the v2 block index, the v3 index removes the ticket hashes associated with
// vote info and revocations.
//
// The new block index is guaranteed to be fully updated if this returns without
// failure.
func migrateBlockIndexVersion2To3(ctx context.Context, db database.DB, dbInfo *databaseInfo) error {
	// Hardcoded bucket names so updates do not affect old upgrades.
	v2BucketName := []byte("blockidx")
	v3BucketName := []byte("blockidxv3")

	log.Info("Reindexing block information in the database.  This may take a " +
		"while...")
	start := time.Now()

	// Create the new block index bucket as needed.
	err := db.Update(func(dbTx database.Tx) error {
		_, err := dbTx.Metadata().CreateBucketIfNotExists(v3BucketName)
		return err
	})
	if err != nil {
		return err
	}

	// doBatch contains the primary logic for upgrading the block index from
	// version 2 to 3 in batches.  This is done because attempting to migrate in
	// a single database transaction could result in massive memory usage and
	// could potentially crash on many systems due to ulimits.
	//
	// It returns the number of entries processed.
	const maxEntries = 20000
	var resumeOffset uint32
	doBatch := func(dbTx database.Tx) (uint32, error) {
		meta := dbTx.Metadata()
		v2BlockIdxBucket := meta.Bucket(v2BucketName)
		if v2BlockIdxBucket == nil {
			return 0, fmt.Errorf("bucket %s does not exist", v2BucketName)
		}

		v3BlockIdxBucket := meta.Bucket(v3BucketName)
		if v3BlockIdxBucket == nil {
			return 0, fmt.Errorf("bucket %s does not exist", v3BucketName)
		}

		// Migrate block index entries so long as the max number of entries for
		// this batch has not been exceeded.
		var numMigrated, numIterated uint32
		err := v2BlockIdxBucket.ForEach(func(key, oldSerialized []byte) error {
			if numMigrated >= maxEntries {
				return errBatchFinished
			}

			// Skip entries that have already been migrated in previous batches.
			numIterated++
			if numIterated-1 < resumeOffset {
				return nil
			}
			resumeOffset++

			// Skip entries that have already been migrated in previous
			// interrupted upgrades.
			if v3BlockIdxBucket.Get(key) != nil {
				return nil
			}

			// Decode the old block index entry.
			var entry blockIndexEntryV2
			_, err := decodeBlockIndexEntryV2(oldSerialized, &entry)
			if err != nil {
				return err
			}

			// Write the block index entry seriliazed with the new format to the
			// new bucket.
			serialized, err := serializeBlockIndexEntry(&blockIndexEntry{
				header:   entry.header,
				status:   entry.status,
				voteInfo: entry.voteInfo,
			})
			if err != nil {
				return err
			}
			err = v3BlockIdxBucket.Put(key, serialized)
			if err != nil {
				return err
			}

			numMigrated++

			if interruptRequested(ctx) {
				return errInterruptRequested
			}

			return nil
		})
		return numMigrated, err
	}

	// Migrate all entries in batches for the reasons mentioned above.
	var totalMigrated uint64
	for {
		var numMigrated uint32
		err := db.Update(func(dbTx database.Tx) error {
			var err error
			numMigrated, err = doBatch(dbTx)
			if errors.Is(err, errInterruptRequested) ||
				errors.Is(err, errBatchFinished) {
				// No error here so the database transaction is not cancelled
				// and therefore outstanding work is written to disk.  The outer
				// function will exit with an interrupted error below due to
				// another interrupted check.
				err = nil
			}
			return err
		})
		if err != nil {
			return err
		}

		if interruptRequested(ctx) {
			return errInterruptRequested
		}

		if numMigrated == 0 {
			break
		}

		totalMigrated += uint64(numMigrated)
		log.Infof("Migrated %d entries (%d total)", numMigrated, totalMigrated)
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	log.Infof("Done migrating block index.  Total entries: %d in %v",
		totalMigrated, elapsed)

	if interruptRequested(ctx) {
		return errInterruptRequested
	}

	// Drop version 2 block index.
	log.Info("Removing old block index entries...")
	start = time.Now()
	err = incrementalFlatDrop(ctx, db, v2BucketName, "old block index")
	if err != nil {
		return err
	}
	elapsed = time.Since(start).Round(time.Millisecond)
	log.Infof("Done removing old block index entries in %v", elapsed)

	// Update and persist the database versions.
	err = db.Update(func(dbTx database.Tx) error {
		dbInfo.bidxVer = 3
		return dbPutDatabaseInfo(dbTx, dbInfo)
	})
	return err
}

// upgradeDB upgrades old database versions to the newest version by applying
// all possible upgrades iteratively.
//
// NOTE: The passed database info will be updated with the latest versions.
func upgradeDB(ctx context.Context, db database.DB, chainParams *chaincfg.Params, dbInfo *databaseInfo) error {
	if dbInfo.version == 1 {
		if err := upgradeToVersion2(db, chainParams, dbInfo); err != nil {
			return err
		}
	}

	// Migrate to the new v2 block index format if needed.  That database
	// version was bumped because prior versions of the software did not have
	// a block index version.
	if dbInfo.version == 2 && dbInfo.bidxVer < 2 {
		if err := upgradeToVersion3(ctx, db, dbInfo); err != nil {
			return err
		}
	}

	// Remove the main chain index from the database if needed.
	if dbInfo.version == 3 {
		if err := upgradeToVersion4(ctx, db, dbInfo); err != nil {
			return err
		}
	}

	// Clear the utxoset, clear the spend journal, reset the best chain back to
	// the genesis block, and mark that a v5 reindex is required if needed.
	if dbInfo.version == 4 {
		err := upgradeToVersion5(ctx, db, chainParams, dbInfo)
		if err != nil {
			return err
		}
	}

	// Update to a version 6 database if needed.  This entails unmarking all
	// blocks previously marked failed so they are eligible for validation again
	// under the new consensus rules and creating and storing version 2 GCS
	// filters for all blocks in the main chain.
	if dbInfo.version == 5 {
		err := upgradeToVersion6(ctx, db, chainParams, dbInfo)
		if err != nil {
			return err
		}
	}

	// Update to the version 3 block index format if needed.
	if dbInfo.version == 6 && dbInfo.bidxVer == 2 {
		err := migrateBlockIndexVersion2To3(ctx, db, dbInfo)
		if err != nil {
			return err
		}
	}

	return nil
}
