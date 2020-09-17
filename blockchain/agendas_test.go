// Copyright (c) 2017-2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"math"
	"testing"
	"time"

	"github.com/hdfchain/hdfd/blockchain/stake"
	"github.com/hdfchain/hdfd/chaincfg"
	"github.com/hdfchain/hdfd/txscript"
)

// testLNFeaturesDeployment ensures the deployment of the LN features agenda
// activates the expected changes for the provided network parameters and
// expected deployment version.
func testLNFeaturesDeployment(t *testing.T, params *chaincfg.Params, deploymentVer uint32) {
	// baseConsensusScriptVerifyFlags are the expected script flags when the
	// agenda is not active.
	const baseConsensusScriptVerifyFlags = txscript.ScriptVerifyCleanStack |
		txscript.ScriptVerifyCheckLockTimeVerify

	// Find the correct deployment for the LN features agenda and ensure it
	// never expires to prevent test failures when the real expiration time
	// passes.  Also, clone the provided parameters first to avoid mutating
	// them.
	params = cloneParams(params)
	var deployment *chaincfg.ConsensusDeployment
	deployments := params.Deployments[deploymentVer]
	for deploymentID, depl := range deployments {
		if depl.Vote.Id == chaincfg.VoteIDLNFeatures {
			deployment = &deployments[deploymentID]
			break
		}
	}
	if deployment == nil {
		t.Fatalf("Unable to find consensus deployement for %s",
			chaincfg.VoteIDLNFeatures)
	}
	deployment.ExpireTime = math.MaxUint64 // Never expires.

	// Find the correct choice for the yes vote.
	const yesVoteID = "yes"
	var yesChoice chaincfg.Choice
	for _, choice := range deployment.Vote.Choices {
		if choice.Id == yesVoteID {
			yesChoice = choice
		}
	}
	if yesChoice.Id != yesVoteID {
		t.Fatalf("Unable to find vote choice for id %q", yesVoteID)
	}

	// Shorter versions of params for convenience.
	stakeValidationHeight := uint32(params.StakeValidationHeight)
	ruleChangeActivationInterval := params.RuleChangeActivationInterval

	tests := []struct {
		name          string
		numNodes      uint32 // num fake nodes to create
		curActive     bool   // whether agenda active for current block
		nextActive    bool   // whether agenda active for NEXT block
		expectedFlags txscript.ScriptFlags
	}{
		{
			name:          "stake validation height",
			numNodes:      stakeValidationHeight,
			curActive:     false,
			nextActive:    false,
			expectedFlags: baseConsensusScriptVerifyFlags,
		},
		{
			name:          "started",
			numNodes:      ruleChangeActivationInterval,
			curActive:     false,
			nextActive:    false,
			expectedFlags: baseConsensusScriptVerifyFlags,
		},
		{
			name:          "lockedin",
			numNodes:      ruleChangeActivationInterval,
			curActive:     false,
			nextActive:    false,
			expectedFlags: baseConsensusScriptVerifyFlags,
		},
		{
			name:          "one before active",
			numNodes:      ruleChangeActivationInterval - 1,
			curActive:     false,
			nextActive:    true,
			expectedFlags: baseConsensusScriptVerifyFlags,
		},
		{
			name:       "exactly active",
			numNodes:   1,
			curActive:  true,
			nextActive: true,
			expectedFlags: baseConsensusScriptVerifyFlags |
				txscript.ScriptVerifyCheckSequenceVerify |
				txscript.ScriptVerifySHA256,
		},
		{
			name:       "one after active",
			numNodes:   1,
			curActive:  true,
			nextActive: true,
			expectedFlags: baseConsensusScriptVerifyFlags |
				txscript.ScriptVerifyCheckSequenceVerify |
				txscript.ScriptVerifySHA256,
		},
	}

	curTimestamp := time.Now()
	bc := newFakeChain(params)
	node := bc.bestChain.Tip()
	for _, test := range tests {
		for i := uint32(0); i < test.numNodes; i++ {
			node = newFakeNode(node, int32(deploymentVer),
				deploymentVer, 0, curTimestamp)

			// Create fake votes that vote yes on the agenda to
			// ensure it is activated.
			for j := uint16(0); j < params.TicketsPerBlock; j++ {
				node.votes = append(node.votes, stake.VoteVersionTuple{
					Version: deploymentVer,
					Bits:    yesChoice.Bits | 0x01,
				})
			}
			bc.bestChain.SetTip(node)
			curTimestamp = curTimestamp.Add(time.Second)
		}

		// Ensure the agenda reports the expected activation status for
		// the current block.
		gotActive, err := bc.isLNFeaturesAgendaActive(node.parent)
		if err != nil {
			t.Errorf("%s: unexpected err: %v", test.name, err)
			continue
		}
		if gotActive != test.curActive {
			t.Errorf("%s: mismatched current active status - got: "+
				"%v, want: %v", test.name, gotActive,
				test.curActive)
			continue
		}

		// Ensure the agenda reports the expected activation status for
		// the NEXT block
		gotActive, err = bc.IsLNFeaturesAgendaActive()
		if err != nil {
			t.Errorf("%s: unexpected err: %v", test.name, err)
			continue
		}
		if gotActive != test.nextActive {
			t.Errorf("%s: mismatched next active status - got: %v, "+
				"want: %v", test.name, gotActive,
				test.nextActive)
			continue
		}

		// Ensure the consensus script verify flags are as expected.
		gotFlags, err := bc.consensusScriptVerifyFlags(node)
		if err != nil {
			t.Errorf("%s: unexpected err: %v", test.name, err)
			continue
		}
		if gotFlags != test.expectedFlags {
			t.Errorf("%s: mismatched flags - got %v, want %v",
				test.name, gotFlags, test.expectedFlags)
			continue
		}
	}
}

// TestLNFeaturesDeployment ensures the deployment of the LN features agenda
// activate the expected changes.
func TestLNFeaturesDeployment(t *testing.T) {
	testLNFeaturesDeployment(t, &chaincfg.MainNetParams, 5)
	testLNFeaturesDeployment(t, &chaincfg.RegNetParams, 6)
}
