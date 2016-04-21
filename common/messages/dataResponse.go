// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

// Communicate a Directory Block State

type DataResponse struct {
	MessageBase
	Timestamp interfaces.Timestamp

	DataType   int // 0 = Entry, 1 = EntryBlock
	DataHash   interfaces.IHash
	DataObject interface{}
}

var _ interfaces.IMsg = (*DataResponse)(nil)

func (m *DataResponse) IsSameAs(b *DataResponse) bool {
	return true
}

func (m *DataResponse) GetHash() interfaces.IHash {
	return nil
}

func (m *DataResponse) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
		data, err := m.MarshalBinary()
		if err != nil {
			return nil
		}
		m.MsgHash = primitives.Sha(data)
	}
	return m.MsgHash
}

func (m *DataResponse) Type() int {
	return constants.DATA_RESPONSE
}

func (m *DataResponse) Int() int {
	return -1
}

func (m *DataResponse) Bytes() []byte {
	return nil
}

func (m *DataResponse) GetTimestamp() interfaces.Timestamp {
	return m.Timestamp
}

// Validate the message, given the state.  Three possible results:
//  < 0 -- Message is invalid.  Discard
//  0   -- Cannot tell if message is Valid
//  1   -- Message is valid
func (m *DataResponse) Validate(state interfaces.IState) int {
	var dataHash interfaces.IHash
	var err error
	switch m.DataType {
	case 0: // DataType = entry
		dataObject, ok := m.DataObject.(interfaces.IEBEntry)
		if !ok {
			return -1
		}
		dataHash = dataObject.GetHash()
	case 1: // DataType = eblock
		dataObject, ok := m.DataObject.(interfaces.IEntryBlock)
		if !ok {
			return -1
		}
		dataHash, err = dataObject.Hash()
		if err != nil {
			return -1
		}
	default:
		// DataType currently not supported, treat as invalid
		return -1
	}

	if dataHash.IsSameAs(m.DataHash) {
		return 1
	}

	return -1
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *DataResponse) Leader(state interfaces.IState) bool {
	return false
}

// Execute the leader functions of the given message
func (m *DataResponse) LeaderExecute(state interfaces.IState) error {
	return fmt.Errorf("Should never execute a DataResponse in the Leader")
}

// Returns true if this is a message for this server to execute as a follower
func (m *DataResponse) Follower(interfaces.IState) bool {
	return true
}

func (m *DataResponse) FollowerExecute(state interfaces.IState) error {
	return state.FollowerExecuteDBState(m)
}

// Acknowledgements do not go into the process list.
func (e *DataResponse) Process(dbheight uint32, state interfaces.IState) bool {
	panic("Should never have its Process() method called")
}

func (e *DataResponse) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *DataResponse) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *DataResponse) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func (m *DataResponse) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()

	newData = data[1:]

	newData, err = m.Timestamp.UnmarshalBinaryData(newData)
	if err != nil {
		return nil, err
	}

	m.DataType = int(newData[0])
	newData = newData[1:]

	m.DataHash = primitives.NewHash(constants.ZERO_HASH)
	newData, err = m.DataHash.UnmarshalBinaryData(newData)
	if err != nil {
		return nil, err
	}

	//UNMARSHAL m.DataObject (interface{})

	m.Peer2peer = true // Always a peer2peer request.

	return data, nil
}

func (m *DataResponse) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *DataResponse) MarshalBinary() ([]byte, error) {

	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, byte(m.Type()))

	t := m.GetTimestamp()
	data, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	buf.Write(data)

	return buf.Bytes(), nil
}

func (m *DataResponse) String() string {
	return fmt.Sprintf("DataResponse Type: %v\n Hash: %x\n Object: %v\n",
		m.DataType,
		m.DataHash.Bytes()[:5],
		m.DataObject)
}

func NewDataResponse(dataObject interface{},
	dataType int,
	dataHash interfaces.IHash) interfaces.IMsg {

	msg := new(DataResponse)

	msg.Peer2peer = true

	msg.DataHash = dataHash
	msg.DataType = dataType
	msg.DataObject = dataObject

	return msg
}
