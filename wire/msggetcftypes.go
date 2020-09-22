// Copyright (c) 2017 The btcsuite developers
// Copyright (c) 2017 The Lightning Network Developers
// Copyright (c) 2018-2020 The Hdfchain developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"fmt"
	"io"
)

// MsgGetCFTypes is the getcftypes message.
type MsgGetCFTypes struct{}

// BtcDecode decodes the receiver from w using the wire protocol encoding.
// This is part of the Message interface implementation.
func (msg *MsgGetCFTypes) BtcDecode(r io.Reader, pver uint32) error {
	const op = "MsgGetCFTypes.BtcDecode"
	if pver < NodeCFVersion {
		msg := fmt.Sprintf("getcftypes message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	return nil
}

// BtcEncode encodes the receiver to w using the wire protocol encoding. This is
// part of the Message interface implementation.
func (msg *MsgGetCFTypes) BtcEncode(w io.Writer, pver uint32) error {
	const op = "MsgGetCFTypes.BtcEncode"
	if pver < NodeCFVersion {
		msg := fmt.Sprintf("getcftypes message invalid for protocol "+
			"version %d", pver)
		return messageError(op, ErrMsgInvalidForPVer, msg)
	}

	return nil
}

// Command returns the protocol command string for the message.  This is part
// of the Message interface implementation.
func (msg *MsgGetCFTypes) Command() string {
	return CmdGetCFTypes
}

// MaxPayloadLength returns the maximum length the payload can be for the
// receiver.  This is part of the Message interface implementation.
func (msg *MsgGetCFTypes) MaxPayloadLength(pver uint32) uint32 {
	// Empty message.
	return 0
}

// NewMsgGetCFTypes returns a new getcftypes message that conforms to the
// Message interface.
func NewMsgGetCFTypes() *MsgGetCFTypes {
	return &MsgGetCFTypes{}
}
