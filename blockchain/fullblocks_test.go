// Copyright (c) 2016 The Hdfchain developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain_test

import (
	"testing"

	"github.com/hdfchain/hdfd/blockchain"
	"github.com/hdfchain/hdfd/blockchain/fullblocktests"
	"github.com/hdfchain/hdfd/chaincfg"
	"github.com/hdfchain/dcrutil"
)

// TestFullBlocks ensures all tests generated by the fullblocktests package
// have the expected result when processed via ProcessBlock.
func TestFullBlocks(t *testing.T) {
	tests, err := fullblocktests.Generate()
	if err != nil {
		t.Fatalf("failed to generate tests: %v", err)
	}

	// Create a new database and chain instance to run tests against.
	chain, teardownFunc, err := chainSetup("fullblocktest",
		&chaincfg.SimNetParams)
	if err != nil {
		t.Fatalf("Failed to setup chain instance: %v", err)
	}
	defer teardownFunc()

	// testAcceptedBlock attempts to process the block in the provided test
	// instance and ensures that it was accepted according to the flags
	// specified in the test.
	testAcceptedBlock := func(item fullblocktests.AcceptedBlock) {
		blockHeight := item.Block.Header.Height
		block := dcrutil.NewBlock(item.Block)
		t.Logf("Testing block %s (hash %s, height %d)",
			item.Name, block.Hash(), blockHeight)

		isMainChain, isOrphan, err := chain.ProcessBlock(block,
			blockchain.BFNone)
		if err != nil {
			t.Fatalf("block %q (hash %s, height %d) should "+
				"have been accepted: %v", item.Name,
				block.Hash(), blockHeight, err)
		}

		// Ensure the main chain and orphan flags match the values
		// specified in the test.
		if isMainChain != item.IsMainChain {
			t.Fatalf("block %q (hash %s, height %d) unexpected main "+
				"chain flag -- got %v, want %v", item.Name,
				block.Hash(), blockHeight, isMainChain,
				item.IsMainChain)
		}
		if isOrphan != item.IsOrphan {
			t.Fatalf("block %q (hash %s, height %d) unexpected "+
				"orphan flag -- got %v, want %v", item.Name,
				block.Hash(), blockHeight, isOrphan,
				item.IsOrphan)
		}
	}

	// testRejectedBlock attempts to process the block in the provided test
	// instance and ensures that it was rejected with the reject code
	// specified in the test.
	testRejectedBlock := func(item fullblocktests.RejectedBlock) {
		blockHeight := item.Block.Header.Height
		block := dcrutil.NewBlock(item.Block)
		t.Logf("Testing block %s (hash %s, height %d)",
			item.Name, block.Hash(), blockHeight)

		_, _, err := chain.ProcessBlock(block, blockchain.BFNone)
		if err == nil {
			t.Fatalf("block %q (hash %s, height %d) should not "+
				"have been accepted", item.Name, block.Hash(),
				blockHeight)
		}

		// Ensure the error code is of the expected type and the reject
		// code matches the value specified in the test instance.
		rerr, ok := err.(blockchain.RuleError)
		if !ok {
			t.Fatalf("block %q (hash %s, height %d) returned "+
				"unexpected error type -- got %T, want "+
				"blockchain.RuleError", item.Name, block.Hash(),
				blockHeight, err)
		}
		if rerr.ErrorCode != item.RejectCode {
			t.Fatalf("block %q (hash %s, height %d) does not have "+
				"expected reject code -- got %v, want %v",
				item.Name, block.Hash(), blockHeight,
				rerr.ErrorCode, item.RejectCode)
		}
	}

	// testOrphanOrRejectedBlock attempts to process the block in the
	// provided test instance and ensures that it was either accepted as an
	// orphan or rejected with a rule violation.
	testOrphanOrRejectedBlock := func(item fullblocktests.OrphanOrRejectedBlock) {
		blockHeight := item.Block.Header.Height
		block := dcrutil.NewBlock(item.Block)
		t.Logf("Testing block %s (hash %s, height %d)",
			item.Name, block.Hash(), blockHeight)

		_, isOrphan, err := chain.ProcessBlock(block, blockchain.BFNone)
		if err != nil {
			// Ensure the error code is of the expected type.
			if _, ok := err.(blockchain.RuleError); !ok {
				t.Fatalf("block %q (hash %s, height %d) "+
					"returned unexpected error type -- "+
					"got %T, want blockchain.RuleError",
					item.Name, block.Hash(), blockHeight,
					err)
			}
		}

		if !isOrphan {
			t.Fatalf("block %q (hash %s, height %d) was accepted, "+
				"but is not considered an orphan", item.Name,
				block.Hash(), blockHeight)
		}
	}

	// testExpectedTip ensures the current tip of the blockchain is the
	// block specified in the provided test instance.
	testExpectedTip := func(item fullblocktests.ExpectedTip) {
		blockHeight := item.Block.Header.Height
		block := dcrutil.NewBlock(item.Block)
		t.Logf("Testing tip for block %s (hash %s, height %d)",
			item.Name, block.Hash(), blockHeight)

		// Ensure hash and height match.
		best := chain.BestSnapshot()
		if *best.Hash != item.Block.BlockHash() ||
			best.Height != int64(item.Block.Header.Height) {

			t.Fatalf("block %q (hash %s, height %d) should be "+
				"the current tip -- got (hash %s, height %d)",
				item.Name, block.Hash(), blockHeight, best.Hash,
				best.Height)
		}
	}

	for testNum, test := range tests {
		for itemNum, item := range test {
			switch item := item.(type) {
			case fullblocktests.AcceptedBlock:
				testAcceptedBlock(item)
			case fullblocktests.RejectedBlock:
				testRejectedBlock(item)
			case fullblocktests.OrphanOrRejectedBlock:
				testOrphanOrRejectedBlock(item)
			case fullblocktests.ExpectedTip:
				testExpectedTip(item)
			default:
				t.Fatalf("test #%d, item #%d is not one of "+
					"the supported test instance types -- "+
					"got type: %T", testNum, itemNum, item)
			}
		}
	}
}
