// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015-2016 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package blockchain

import (
	"fmt"

	"github.com/hdfchain/hdfd/chaincfg/chainhash"
	"github.com/hdfchain/hdfd/dcrutil"
)

// NotificationType represents the type of a notification message.
type NotificationType int

// NotificationCallback is used for a caller to provide a callback for
// notifications about various chain events.
type NotificationCallback func(*Notification)

// Constants for the type of a notification message.
const (
	// NTBlockAccepted indicates the associated block was accepted into
	// the block chain.  Note that this does not necessarily mean it was
	// added to the main chain.  For that, use NTBlockConnected.
	NTBlockAccepted NotificationType = iota

	// NTBlockConnected indicates the associated block was connected to the
	// main chain.
	NTBlockConnected

	// NTBlockDisconnected indicates the associated block was disconnected
	// from the main chain.
	NTBlockDisconnected

	// NTReorganization indicates that a blockchain reorganization is in
	// progress.
	NTReorganization

	// NTSpentAndMissedTickets indicates spent or missed tickets from a newly
	// accepted block.
	NTSpentAndMissedTickets

	// NTSpentAndMissedTickets indicates newly maturing tickets from a newly
	// accepted block.
	NTNewTickets
)

// notificationTypeStrings is a map of notification types back to their constant
// names for pretty printing.
var notificationTypeStrings = map[NotificationType]string{
	NTBlockAccepted:         "NTBlockAccepted",
	NTBlockConnected:        "NTBlockConnected",
	NTBlockDisconnected:     "NTBlockDisconnected",
	NTReorganization:        "NTReorganization",
	NTSpentAndMissedTickets: "NTSpentAndMissedTickets",
	NTNewTickets:            "NTNewTickets",
}

// String returns the NotificationType in human-readable form.
func (n NotificationType) String() string {
	if s, ok := notificationTypeStrings[n]; ok {
		return s
	}
	return fmt.Sprintf("Unknown Notification Type (%d)", int(n))
}

// BlockAcceptedNtfnsData is the structure for data indicating information
// about an accepted block.  Note that this does not necessarily mean the block
// that was accepted extended the best chain as it might have created or
// extended a side chain.
type BlockAcceptedNtfnsData struct {
	// BestHeight is the height of the current best chain.  Since the accepted
	// block might be on a side chain, this is not necessarily the same as the
	// height of the accepted block.
	BestHeight int64

	// ForkLen is the length of the side chain the block extended or zero in the
	// case the block extended the main chain.
	//
	// This can be used in conjunction with the height of the accepted block to
	// determine the height at which the side chain the block created or
	// extended forked from the best chain.
	ForkLen int64

	// Block is the block that was accepted into the chain.
	Block *dcrutil.Block
}

// ReorganizationNtfnsData is the structure for data indicating information
// about a reorganization.
type ReorganizationNtfnsData struct {
	OldHash   chainhash.Hash
	OldHeight int64
	NewHash   chainhash.Hash
	NewHeight int64
}

// TicketNotificationsData is the structure for new/spent/missed ticket
// notifications at blockchain HEAD that are outgoing from chain.
type TicketNotificationsData struct {
	Hash            chainhash.Hash
	Height          int64
	StakeDifficulty int64
	TicketsSpent    []chainhash.Hash
	TicketsMissed   []chainhash.Hash
	TicketsNew      []chainhash.Hash
}

// Notification defines notification that is sent to the caller via the callback
// function provided during the call to New and consists of a notification type
// as well as associated data that depends on the type as follows:
// 	- NTBlockAccepted:         *BlockAcceptedNtfnsData
// 	- NTBlockConnected:        []*dcrutil.Block of len 2
// 	- NTBlockDisconnected:     []*dcrutil.Block of len 2
//  - NTReorganization:        *ReorganizationNtfnsData
//  - NTSpentAndMissedTickets: *TicketNotificationsData
//  - NTNewTickets:            *TicketNotificationsData
type Notification struct {
	Type NotificationType
	Data interface{}
}

// sendNotification sends a notification with the passed type and data if the
// caller requested notifications by providing a callback function in the call
// to New.
func (b *BlockChain) sendNotification(typ NotificationType, data interface{}) {
	// Ignore it if the caller didn't request notifications.
	if b.notifications == nil {
		return
	}

	// Generate and send the notification.
	n := Notification{Type: typ, Data: data}
	b.notifications(&n)
}
