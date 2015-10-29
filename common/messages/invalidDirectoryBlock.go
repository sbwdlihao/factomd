// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages

import (
	"bytes"
	"fmt"
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

//A placeholder structure for messages
type InvalidDirectoryBlock struct {
}

var _ interfaces.IMsg = (*InvalidDirectoryBlock)(nil)

func (m *InvalidDirectoryBlock) Type() int {
	return constants.INVALID_DIRECTORY_BLOCK_MSG
}

func (m *InvalidDirectoryBlock) Int() int {
	return -1
}

func (m *InvalidDirectoryBlock) Bytes() []byte {
	return nil
}

func (m *InvalidDirectoryBlock) UnmarshalBinaryData(data []byte) (newdata []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()

	return nil, nil
}

func (m *InvalidDirectoryBlock) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *InvalidDirectoryBlock) MarshalBinary() (data []byte, err error) {
	return nil, nil
}

func (m *InvalidDirectoryBlock) String() string {
	return ""
}

func (m *InvalidDirectoryBlock) DBHeight() int {
	return 0
}

func (m *InvalidDirectoryBlock) ChainID() []byte {
	return nil
}

func (m *InvalidDirectoryBlock) ListHeight() int {
	return 0
}

func (m *InvalidDirectoryBlock) SerialHash() []byte {
	return nil
}

func (m *InvalidDirectoryBlock) Signature() []byte {
	return nil
}

// Validate the message, given the state.  Three possible results:
//  < 0 -- Message is invalid.  Discard
//  0   -- Cannot tell if message is Valid
//  1   -- Message is valid
func (m *InvalidDirectoryBlock) Validate(interfaces.IState) int {
	return 0
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *InvalidDirectoryBlock) Leader(state interfaces.IState) bool {
	switch state.GetNetworkNumber() {
	case 0: // Main Network
		panic("Not implemented yet")
	case 1: // Test Network
		panic("Not implemented yet")
	case 2: // Local Network
		panic("Not implemented yet")
	default:
		panic("Not implemented yet")
	}

}

// Execute the leader functions of the given message
func (m *InvalidDirectoryBlock) LeaderExecute(state interfaces.IState) error {
	return nil
}

// Returns true if this is a message for this server to execute as a follower
func (m *InvalidDirectoryBlock) Follower(interfaces.IState) bool {
	return true
}

func (m *InvalidDirectoryBlock) FollowerExecute(interfaces.IState) error {
	return nil
}

func (e *InvalidDirectoryBlock) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *InvalidDirectoryBlock) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *InvalidDirectoryBlock) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}