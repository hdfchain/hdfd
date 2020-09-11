// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Hdfchain developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain_test

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hdfchain/hdfd/blockchain/v2"
	"github.com/hdfchain/hdfd/chaincfg/v2"
	"github.com/hdfchain/hdfd/database/v2"
	_ "github.com/hdfchain/hdfd/database/v2/ffldb"
	"github.com/hdfchain/hdfd/dcrutil/v2"
)

// This example demonstrates how to create a new chain instance and use
// ProcessBlock to attempt to add a block to the chain.  As the package
// overview documentation describes, this includes all of the Hdfchain consensus
// rules.  This example intentionally attempts to insert a duplicate genesis
// block to illustrate how an invalid block is handled.
func ExampleBlockChain_ProcessBlock() {
	// Create a new database to store the accepted blocks into.  Typically
	// this would be opening an existing database and would not be deleting
	// and creating a new database like this, but it is done here so this is
	// a complete working example and does not leave temporary files laying
	// around.
	mainNetParams := chaincfg.MainNetParams()
	dbPath := filepath.Join(os.TempDir(), "exampleprocessblock")
	_ = os.RemoveAll(dbPath)
	db, err := database.Create("ffldb", dbPath, mainNetParams.Net)
	if err != nil {
		fmt.Printf("Failed to create database: %v\n", err)
		return
	}
	defer os.RemoveAll(dbPath)
	defer db.Close()

	// Create a new BlockChain instance using the underlying database for
	// the main bitcoin network.  This example does not demonstrate some
	// of the other available configuration options such as specifying a
	// notification callback and signature cache.  Also, the caller would
	// ordinarily keep a reference to the median time source and add time
	// values obtained from other peers on the network so the local time is
	// adjusted to be in agreement with other peers.
	chain, err := blockchain.New(&blockchain.Config{
		DB:          db,
		ChainParams: mainNetParams,
		TimeSource:  blockchain.NewMedianTime(),
	})
	if err != nil {
		fmt.Printf("Failed to create chain instance: %v\n", err)
		return
	}

	// Process a block.  For this example, we are going to intentionally
	// cause an error by trying to process the genesis block which already
	// exists.
	genesisBlock := dcrutil.NewBlock(mainNetParams.GenesisBlock)
	forkLen, isOrphan, err := chain.ProcessBlock(genesisBlock,
		blockchain.BFNone)
	if err != nil {
		fmt.Printf("Failed to create chain instance: %v\n", err)
		return
	}
	isMainChain := !isOrphan && forkLen == 0
	fmt.Printf("Block accepted. Is it on the main chain?: %v", isMainChain)
	fmt.Printf("Block accepted. Is it an orphan?: %v", isOrphan)

	// This output is dependent on the genesis block, and needs to be
	// updated if the mainnet genesis block is updated.
	// Output:
	// Failed to process block: already have block 267a53b5ee86c24a48ec37aee4f4e7c0c4004892b7259e695e9f5b321f1ab9d2
}
