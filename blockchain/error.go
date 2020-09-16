// Copyright (c) 2014-2016 The btcsuite developers
// Copyright (c) 2015-2019 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"fmt"
)

// VoteVersionError identifies an error that indicates a vote version was
// specified that does not exist.
type VoteVersionError uint32

// Error returns the assertion error as a human-readable string and satisfies
// the error interface.
func (e VoteVersionError) Error() string {
	return fmt.Sprintf("stake version %v does not exist", uint32(e))
}

// HashError identifies an error that indicates a hash was specified that does
// not exist.
type HashError string

// Error returns the error as a human-readable string and satisfies the error
// interface.
func (e HashError) Error() string {
	return fmt.Sprintf("hash %v does not exist", string(e))
}

// DeploymentError identifies an error that indicates a deployment ID was
// specified that does not exist.
type DeploymentError string

// Error returns the assertion error as a human-readable string and satisfies
// the error interface.
func (e DeploymentError) Error() string {
	return fmt.Sprintf("deployment ID %v does not exist", string(e))
}

// DuplicateDeploymentError identifies an error that indicates a duplicate
// deployment ID was specified in the network parameter deployment definitions.
type DuplicateDeploymentError string

// Error returns the assertion error as a human-readable string and satisfies
// the error interface.
func (e DuplicateDeploymentError) Error() string {
	return fmt.Sprintf("deployment ID %v exists in more than one deployment",
		string(e))
}

// NoFilterError identifies an error that indicates a filter for a given block
// hash does not exist.
type NoFilterError string

// Error returns the error as a human-readable string and satisfies the error
// interface.
func (e NoFilterError) Error() string {
	return fmt.Sprintf("no filter available for block %s", string(e))
}

// AssertError identifies an error that indicates an internal code consistency
// issue and should be treated as a critical and unrecoverable error.
type AssertError string

// Error returns the assertion error as a huma-readable string and satisfies
// the error interface.
func (e AssertError) Error() string {
	return "assertion failed: " + string(e)
}

// ErrorCode identifies a kind of error.
type ErrorCode int

// These constants are used to identify a specific RuleError.
const (
	// ErrDuplicateBlock indicates a block with the same hash already
	// exists.
	ErrDuplicateBlock ErrorCode = iota

	// ErrMissingParent indicates that the block was an orphan.
	ErrMissingParent

	// ErrBlockTooBig indicates the serialized block size exceeds the
	// maximum allowed size.
	ErrBlockTooBig

	// ErrWrongBlockSize indicates that the block size in the header is not
	// the actual serialized size of the block.
	ErrWrongBlockSize

	// ErrBlockVersionTooOld indicates the block version is too old and is
	// no longer accepted since the majority of the network has upgraded
	// to a newer version.
	ErrBlockVersionTooOld

	// ErrBadStakeVersionindicates the block version is too old and is no
	// longer accepted since the majority of the network has upgraded to a
	// newer version.
	ErrBadStakeVersion

	// ErrInvalidTime indicates the time in the passed block has a precision
	// that is more than one second.  The chain consensus rules require
	// timestamps to have a maximum precision of one second.
	ErrInvalidTime

	// ErrTimeTooOld indicates the time is either before the median time of
	// the last several blocks per the chain consensus rules or prior to the
	// most recent checkpoint.
	ErrTimeTooOld

	// ErrTimeTooNew indicates the time is too far in the future as compared
	// the current time.
	ErrTimeTooNew

	// ErrDifficultyTooLow indicates the difficulty for the block is lower
	// than the difficulty required by the most recent checkpoint.
	ErrDifficultyTooLow

	// ErrUnexpectedDifficulty indicates specified bits do not align with
	// the expected value either because it doesn't match the calculated
	// value based on difficulty regarding the rules or it is out of the
	// valid range.
	ErrUnexpectedDifficulty

	// ErrHighHash indicates the block does not hash to a value which is
	// lower than the required target difficultly.
	ErrHighHash

	// ErrBadMerkleRoot indicates the calculated merkle root does not match
	// the expected value.
	ErrBadMerkleRoot

	// ErrBadCommitmentRoot indicates the calculated commitment root does
	// not match the expected value.
	ErrBadCommitmentRoot

	// ErrBadCheckpoint indicates a block that is expected to be at a
	// checkpoint height does not match the expected one.
	ErrBadCheckpoint

	// ErrForkTooOld indicates a block is attempting to fork the block chain
	// before the most recent checkpoint.
	ErrForkTooOld

	// ErrCheckpointTimeTooOld indicates a block has a timestamp before the
	// most recent checkpoint.
	ErrCheckpointTimeTooOld

	// ErrNoTransactions indicates the block does not have a least one
	// transaction.  A valid block must have at least the coinbase
	// transaction.
	ErrNoTransactions

	// ErrTooManyTransactions indicates the block has more transactions than
	// are allowed.
	ErrTooManyTransactions

	// ErrNoTxInputs indicates a transaction does not have any inputs.  A
	// valid transaction must have at least one input.
	ErrNoTxInputs

	// ErrNoTxOutputs indicates a transaction does not have any outputs.  A
	// valid transaction must have at least one output.
	ErrNoTxOutputs

	// ErrTxTooBig indicates a transaction exceeds the maximum allowed size
	// when serialized.
	ErrTxTooBig

	// ErrBadTxOutValue indicates an output value for a transaction is
	// invalid in some way such as being out of range.
	ErrBadTxOutValue

	// ErrDuplicateTxInputs indicates a transaction references the same
	// input more than once.
	ErrDuplicateTxInputs

	// ErrBadTxInput indicates a transaction input is invalid in some way
	// such as referencing a previous transaction outpoint which is out of
	// range or not referencing one at all.
	ErrBadTxInput

	// ErrMissingTxOut indicates a transaction output referenced by an input
	// either does not exist or has already been spent.
	ErrMissingTxOut

	// ErrUnfinalizedTx indicates a transaction has not been finalized.
	// A valid block may only contain finalized transactions.
	ErrUnfinalizedTx

	// ErrDuplicateTx indicates a block contains an identical transaction
	// (or at least two transactions which hash to the same value).  A
	// valid block may only contain unique transactions.
	ErrDuplicateTx

	// ErrOverwriteTx indicates a block contains a transaction that has
	// the same hash as a previous transaction which has not been fully
	// spent.
	ErrOverwriteTx

	// ErrImmatureSpend indicates a transaction is attempting to spend a
	// coinbase that has not yet reached the required maturity.
	ErrImmatureSpend

	// ErrSpendTooHigh indicates a transaction is attempting to spend more
	// value than the sum of all of its inputs.
	ErrSpendTooHigh

	// ErrBadFees indicates the total fees for a block are invalid due to
	// exceeding the maximum possible value.
	ErrBadFees

	// ErrTooManySigOps indicates the total number of signature operations
	// for a transaction or block exceed the maximum allowed limits.
	ErrTooManySigOps

	// ErrFirstTxNotCoinbase indicates the first transaction in a block
	// is not a coinbase transaction.
	ErrFirstTxNotCoinbase

	// ErrCoinbaseHeight indicates that the encoded height in the coinbase
	// is incorrect.
	ErrCoinbaseHeight

	// ErrMultipleCoinbases indicates a block contains more than one
	// coinbase transaction.
	ErrMultipleCoinbases

	// ErrStakeTxInRegularTree indicates a stake transaction was found in
	// the regular transaction tree.
	ErrStakeTxInRegularTree

	// ErrRegTxInStakeTree indicates that a regular transaction was found in
	// the stake transaction tree.
	ErrRegTxInStakeTree

	// ErrBadCoinbaseScriptLen indicates the length of the signature script
	// for a coinbase transaction is not within the valid range.
	ErrBadCoinbaseScriptLen

	// ErrBadCoinbaseValue indicates the amount of a coinbase value does
	// not match the expected value of the subsidy plus the sum of all fees.
	ErrBadCoinbaseValue

	// ErrBadCoinbaseOutpoint indicates that the outpoint used by a coinbase
	// as input was non-null.
	ErrBadCoinbaseOutpoint

	// ErrBadCoinbaseFraudProof indicates that the fraud proof for a coinbase
	// input was non-null.
	ErrBadCoinbaseFraudProof

	// ErrBadCoinbaseAmountIn indicates that the AmountIn (=subsidy) for a
	// coinbase input was incorrect.
	ErrBadCoinbaseAmountIn

	// ErrBadStakebaseAmountIn indicates that the AmountIn (=subsidy) for a
	// stakebase input was incorrect.
	ErrBadStakebaseAmountIn

	// ErrBadStakebaseScriptLen indicates the length of the signature script
	// for a stakebase transaction is not within the valid range.
	ErrBadStakebaseScriptLen

	// ErrBadStakebaseScrVal indicates the signature script for a stakebase
	// transaction was not set to the network consensus value.
	ErrBadStakebaseScrVal

	// ErrScriptMalformed indicates a transaction script is malformed in
	// some way.  For example, it might be longer than the maximum allowed
	// length or fail to parse.
	ErrScriptMalformed

	// ErrScriptValidation indicates the result of executing transaction
	// script failed.  The error covers any failure when executing scripts
	// such signature verification failures and execution past the end of
	// the stack.
	ErrScriptValidation

	// ErrNotEnoughStake indicates that there was for some SStx in a given block,
	// the given SStx did not have enough stake to meet the network target.
	ErrNotEnoughStake

	// ErrStakeBelowMinimum indicates that for some SStx in a given block,
	// the given SStx had an amount of stake below the minimum network target.
	ErrStakeBelowMinimum

	// ErrNonstandardStakeTx indicates that a block contained a stake tx that
	// was not one of the allowed types of a stake transactions.
	ErrNonstandardStakeTx

	// ErrNotEnoughVotes indicates that a block contained less than a majority
	// of voters.
	ErrNotEnoughVotes

	// ErrTooManyVotes indicates that a block contained more than the maximum
	// allowable number of votes.
	ErrTooManyVotes

	// ErrFreshStakeMismatch indicates that a block's header contained a different
	// number of SStx as compared to what was found in the block.
	ErrFreshStakeMismatch

	// ErrTooManySStxs indicates that more than the allowed number of SStx was
	// found in a block.
	ErrTooManySStxs

	// ErrInvalidEarlyStakeTx indicates that a tx type other than SStx was found
	// in the stake tx tree before the period when stake validation begins, or
	// before the stake tx type could possibly be included in the block.
	ErrInvalidEarlyStakeTx

	// ErrTicketUnavailable indicates that a vote in the block spent a ticket
	// that could not be found.
	ErrTicketUnavailable

	// ErrVotesOnWrongBlock indicates that an SSGen voted on a block not the
	// block's parent, and so was ineligible for inclusion into that block.
	ErrVotesOnWrongBlock

	// ErrVotesMismatch indicates that the number of SSGen in the block was not
	// equivalent to the number of votes provided in the block header.
	ErrVotesMismatch

	// ErrIncongruentVotebit indicates that the first votebit in votebits was not
	// the same as that determined by the majority of voters in the SSGen tx
	// included in the block.
	ErrIncongruentVotebit

	// ErrInvalidSSRtx indicates than an SSRtx in a block could not be found to
	// have a valid missed sstx input as per the stake ticket database.
	ErrInvalidSSRtx

	// ErrInvalidRevNum indicates that the number of revocations from the
	// header was not the same as the number of SSRtx included in the block.
	ErrRevocationsMismatch

	// ErrTooManyRevocations indicates more revocations were found in a block
	// than were allowed.
	ErrTooManyRevocations

	// ErrTicketCommitment indicates that a ticket commitment contains an amount
	// that does not coincide with the associated ticket input amount.
	ErrTicketCommitment

	// ErrInvalidVoteInput indicates that an input to a vote transaction is
	// either not a stake ticket submission or is not a supported version.
	ErrInvalidVoteInput

	// ErrBadNumPayees indicates that either a vote or revocation transaction
	// does not make the correct number of payments per the associated ticket
	// commitments.
	ErrBadNumPayees

	// ErrBadPayeeScriptVersion indicates that either a vote or revocation
	// transaction output that corresponds to a ticket commitment does not use
	// a supported script version.
	ErrBadPayeeScriptVersion

	// ErrBadPayeeScriptType indicates that either a vote or revocation
	// transaction output that corresponds to a ticket commitment does not pay
	// to the same script type required by the commitment.
	ErrBadPayeeScriptType

	// ErrBadPayeeScriptType indicates that either a vote or revocation
	// transaction output that corresponds to a ticket commitment does not pay
	// to the hash required by the commitment.
	ErrMismatchedPayeeHash

	// ErrBadPayeeValue indicates that either a vote or revocation transaction
	// output that corresponds to a ticket commitment does not pay the expected
	// amount required by the commitment.
	ErrBadPayeeValue

	// ErrSSGenSubsidy indicates that there was an error in the amount of subsidy
	// generated in the vote.
	ErrSSGenSubsidy

	// ErrImmatureTicketSpend indicates that a vote or revocation is attempting
	// to spend a ticket submission output that has not yet reached the required
	// maturity.
	ErrImmatureTicketSpend

	// ErrTicketInputScript indicates that a ticket input is not one of the
	// supported script forms or versions.
	ErrTicketInputScript

	// ErrInvalidRevokeInput indicates that an input to a revocation transaction
	// is either not a stake ticket submission or is not a supported version.
	ErrInvalidRevokeInput

	// ErrSSRtxPayees indicates that the SSRtx failed to pay out to the committed
	// addresses or amounts from the originating SStx.
	ErrSSRtxPayees

	// ErrTxSStxOutSpend indicates that a non SSGen or SSRtx tx attempted to spend
	// an OP_SSTX tagged output from an SStx.
	ErrTxSStxOutSpend

	// ErrRegTxCreateStakeOut indicates that a regular tx attempted to create
	// a stake tagged output.
	ErrRegTxCreateStakeOut

	// ErrInvalidFinalState indicates that the final state of the PRNG included
	// in the block differed from the calculated final state.
	ErrInvalidFinalState

	// ErrPoolSize indicates an error in the ticket pool size for this block.
	ErrPoolSize

	// ErrForceReorgWrongChain indicates that a reroganization was attempted
	// to be forced, but the chain indicated was not mirrored by b.bestChain.
	ErrForceReorgWrongChain

	// ErrForceReorgMissingChild indicates that a reroganization was attempted
	// to be forced, but the child node to reorganize to could not be found.
	ErrForceReorgMissingChild

	// ErrBadStakebaseValue indicates that a block's stake tx tree has spent
	// more than it is allowed.
	ErrBadStakebaseValue

	// ErrDiscordantTxTree specifies that a given origin tx's content
	// indicated that it should exist in a different tx tree than the
	// one given in the TxIn outpoint.
	ErrDiscordantTxTree

	// ErrStakeFees indicates an error with the fees found in the stake
	// transaction tree.
	ErrStakeFees

	// ErrNoStakeTx indicates there were no stake transactions found in a
	// block after stake validation height.
	ErrNoStakeTx

	// ErrBadBlockHeight indicates that a block header's embedded block height
	// was different from where it was actually embedded in the block chain.
	ErrBadBlockHeight

	// ErrBlockOneTx indicates that block height 1 failed to correct generate
	// the block one initial payout transaction.
	ErrBlockOneTx

	// ErrBlockOneTx indicates that block height 1 coinbase transaction in
	// zero was incorrect in some way.
	ErrBlockOneInputs

	// ErrBlockOneOutputs indicates that block height 1 failed to incorporate
	// the ledger addresses correctly into the transaction's outputs.
	ErrBlockOneOutputs

	// ErrNoTax indicates that there was no tax present in the coinbase of a
	// block after height 1.
	ErrNoTax

	// ErrExpiredTx indicates that the transaction is currently expired.
	ErrExpiredTx

	// ErrExpiryTxSpentEarly indicates that an output from a transaction
	// that included an expiry field was spent before coinbase maturity
	// many blocks had passed in the blockchain.
	ErrExpiryTxSpentEarly

	// ErrFraudAmountIn indicates the witness amount given was fraudulent.
	ErrFraudAmountIn

	// ErrFraudBlockHeight indicates the witness block height given was
	// fraudulent.
	ErrFraudBlockHeight

	// ErrFraudBlockIndex indicates the witness block index given was
	// fraudulent.
	ErrFraudBlockIndex

	// ErrZeroValueOutputSpend indicates that a transaction attempted to spend a
	// zero value output.
	ErrZeroValueOutputSpend

	// ErrInvalidEarlyVoteBits indicates that a block before stake validation
	// height had an unallowed vote bits value.
	ErrInvalidEarlyVoteBits

	// ErrInvalidEarlyFinalState indicates that a block before stake validation
	// height had a non-zero final state.
	ErrInvalidEarlyFinalState

	// ErrKnownInvalidBlock indicates that this block has previously failed
	// validation.
	ErrKnownInvalidBlock

	// ErrInvalidAncestorBlock indicates that an ancestor of this block has
	// failed validation.
	ErrInvalidAncestorBlock

	// ErrInvalidTemplateParent indicates that a block template builds on a
	// block that is either not the current best chain tip or its parent.
	ErrInvalidTemplateParent

	// numErrorCodes is the maximum error code number used in tests.
	numErrorCodes
)

// Map of ErrorCode values back to their constant names for pretty printing.
var errorCodeStrings = map[ErrorCode]string{
	ErrDuplicateBlock:         "ErrDuplicateBlock",
	ErrMissingParent:          "ErrMissingParent",
	ErrBlockTooBig:            "ErrBlockTooBig",
	ErrWrongBlockSize:         "ErrWrongBlockSize",
	ErrBlockVersionTooOld:     "ErrBlockVersionTooOld",
	ErrBadStakeVersion:        "ErrBadStakeVersion",
	ErrInvalidTime:            "ErrInvalidTime",
	ErrTimeTooOld:             "ErrTimeTooOld",
	ErrTimeTooNew:             "ErrTimeTooNew",
	ErrDifficultyTooLow:       "ErrDifficultyTooLow",
	ErrUnexpectedDifficulty:   "ErrUnexpectedDifficulty",
	ErrHighHash:               "ErrHighHash",
	ErrBadMerkleRoot:          "ErrBadMerkleRoot",
	ErrBadCommitmentRoot:      "ErrBadCommitmentRoot",
	ErrBadCheckpoint:          "ErrBadCheckpoint",
	ErrForkTooOld:             "ErrForkTooOld",
	ErrCheckpointTimeTooOld:   "ErrCheckpointTimeTooOld",
	ErrNoTransactions:         "ErrNoTransactions",
	ErrTooManyTransactions:    "ErrTooManyTransactions",
	ErrNoTxInputs:             "ErrNoTxInputs",
	ErrNoTxOutputs:            "ErrNoTxOutputs",
	ErrTxTooBig:               "ErrTxTooBig",
	ErrBadTxOutValue:          "ErrBadTxOutValue",
	ErrDuplicateTxInputs:      "ErrDuplicateTxInputs",
	ErrBadTxInput:             "ErrBadTxInput",
	ErrMissingTxOut:           "ErrMissingTxOut",
	ErrUnfinalizedTx:          "ErrUnfinalizedTx",
	ErrDuplicateTx:            "ErrDuplicateTx",
	ErrOverwriteTx:            "ErrOverwriteTx",
	ErrImmatureSpend:          "ErrImmatureSpend",
	ErrSpendTooHigh:           "ErrSpendTooHigh",
	ErrBadFees:                "ErrBadFees",
	ErrTooManySigOps:          "ErrTooManySigOps",
	ErrFirstTxNotCoinbase:     "ErrFirstTxNotCoinbase",
	ErrCoinbaseHeight:         "ErrCoinbaseHeight",
	ErrMultipleCoinbases:      "ErrMultipleCoinbases",
	ErrStakeTxInRegularTree:   "ErrStakeTxInRegularTree",
	ErrRegTxInStakeTree:       "ErrRegTxInStakeTree",
	ErrBadCoinbaseScriptLen:   "ErrBadCoinbaseScriptLen",
	ErrBadCoinbaseValue:       "ErrBadCoinbaseValue",
	ErrBadCoinbaseOutpoint:    "ErrBadCoinbaseOutpoint",
	ErrBadCoinbaseFraudProof:  "ErrBadCoinbaseFraudProof",
	ErrBadCoinbaseAmountIn:    "ErrBadCoinbaseAmountIn",
	ErrBadStakebaseAmountIn:   "ErrBadStakebaseAmountIn",
	ErrBadStakebaseScriptLen:  "ErrBadStakebaseScriptLen",
	ErrBadStakebaseScrVal:     "ErrBadStakebaseScrVal",
	ErrScriptMalformed:        "ErrScriptMalformed",
	ErrScriptValidation:       "ErrScriptValidation",
	ErrNotEnoughStake:         "ErrNotEnoughStake",
	ErrStakeBelowMinimum:      "ErrStakeBelowMinimum",
	ErrNonstandardStakeTx:     "ErrNonstandardStakeTx",
	ErrNotEnoughVotes:         "ErrNotEnoughVotes",
	ErrTooManyVotes:           "ErrTooManyVotes",
	ErrFreshStakeMismatch:     "ErrFreshStakeMismatch",
	ErrTooManySStxs:           "ErrTooManySStxs",
	ErrInvalidEarlyStakeTx:    "ErrInvalidEarlyStakeTx",
	ErrTicketUnavailable:      "ErrTicketUnavailable",
	ErrVotesOnWrongBlock:      "ErrVotesOnWrongBlock",
	ErrVotesMismatch:          "ErrVotesMismatch",
	ErrIncongruentVotebit:     "ErrIncongruentVotebit",
	ErrInvalidSSRtx:           "ErrInvalidSSRtx",
	ErrRevocationsMismatch:    "ErrRevocationsMismatch",
	ErrTooManyRevocations:     "ErrTooManyRevocations",
	ErrTicketCommitment:       "ErrTicketCommitment",
	ErrInvalidVoteInput:       "ErrInvalidVoteInput",
	ErrBadNumPayees:           "ErrBadNumPayees",
	ErrBadPayeeScriptVersion:  "ErrBadPayeeScriptVersion",
	ErrBadPayeeScriptType:     "ErrBadPayeeScriptType",
	ErrMismatchedPayeeHash:    "ErrMismatchedPayeeHash",
	ErrBadPayeeValue:          "ErrBadPayeeValue",
	ErrSSGenSubsidy:           "ErrSSGenSubsidy",
	ErrImmatureTicketSpend:    "ErrImmatureTicketSpend",
	ErrTicketInputScript:      "ErrTicketInputScript",
	ErrInvalidRevokeInput:     "ErrInvalidRevokeInput",
	ErrSSRtxPayees:            "ErrSSRtxPayees",
	ErrTxSStxOutSpend:         "ErrTxSStxOutSpend",
	ErrRegTxCreateStakeOut:    "ErrRegTxCreateStakeOut",
	ErrInvalidFinalState:      "ErrInvalidFinalState",
	ErrPoolSize:               "ErrPoolSize",
	ErrForceReorgWrongChain:   "ErrForceReorgWrongChain",
	ErrForceReorgMissingChild: "ErrForceReorgMissingChild",
	ErrBadStakebaseValue:      "ErrBadStakebaseValue",
	ErrDiscordantTxTree:       "ErrDiscordantTxTree",
	ErrStakeFees:              "ErrStakeFees",
	ErrNoStakeTx:              "ErrNoStakeTx",
	ErrBadBlockHeight:         "ErrBadBlockHeight",
	ErrBlockOneTx:             "ErrBlockOneTx",
	ErrBlockOneInputs:         "ErrBlockOneInputs",
	ErrBlockOneOutputs:        "ErrBlockOneOutputs",
	ErrNoTax:                  "ErrNoTax",
	ErrExpiredTx:              "ErrExpiredTx",
	ErrExpiryTxSpentEarly:     "ErrExpiryTxSpentEarly",
	ErrFraudAmountIn:          "ErrFraudAmountIn",
	ErrFraudBlockHeight:       "ErrFraudBlockHeight",
	ErrFraudBlockIndex:        "ErrFraudBlockIndex",
	ErrZeroValueOutputSpend:   "ErrZeroValueOutputSpend",
	ErrInvalidEarlyVoteBits:   "ErrInvalidEarlyVoteBits",
	ErrInvalidEarlyFinalState: "ErrInvalidEarlyFinalState",
	ErrKnownInvalidBlock:      "ErrKnownInvalidBlock",
	ErrInvalidAncestorBlock:   "ErrInvalidAncestorBlock",
	ErrInvalidTemplateParent:  "ErrInvalidTemplateParent",
}

// String returns the ErrorCode as a human-readable name.
func (e ErrorCode) String() string {
	if s := errorCodeStrings[e]; s != "" {
		return s
	}
	return fmt.Sprintf("Unknown ErrorCode (%d)", int(e))
}

// RuleError identifies a rule violation.  It is used to indicate that
// processing of a block or transaction failed due to one of the many validation
// rules.  The caller can use type assertions to determine if a failure was
// specifically due to a rule violation and access the ErrorCode field to
// ascertain the specific reason for the rule violation.
type RuleError struct {
	ErrorCode   ErrorCode // Describes the kind of error
	Description string    // Human readable description of the issue
}

// Error satisfies the error interface and prints human-readable errors.
func (e RuleError) Error() string {
	return e.Description
}

// ruleError creates a RuleError given a set of arguments.
func ruleError(c ErrorCode, desc string) RuleError {
	return RuleError{ErrorCode: c, Description: desc}
}
