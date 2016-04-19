// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages

import (
	"bytes"
	//	"encoding/binary"
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

	DataType   string
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
	switch m.DataType {
	case "dblock":
		fmt.Println("Dblock")
	case "entry":
		fmt.Println("Entry")
	case "eblock":
		fmt.Println("Eblock")
	case "fblock":
		fmt.Println("Fblock")
	case "ecblock":
		fmt.Println("ECblock")
	case "ablock":
		fmt.Println("Ablock")
	default:
		// DataType currently not supported, treat as invalid
		return -1
	}

	return 1
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
		return
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()

	m.Peer2peer = true

	newData = data[1:] // Skip our type;  Someone else's problem.

	//TODO: Unmarshal relevant data
	return
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
	dataType string,
	dataHash interfaces.IHash) interfaces.IMsg {

	msg := new(DataResponse)

	msg.Peer2peer = true

	msg.DataHash = dataHash
	msg.DataType = dataType
	msg.DataObject = dataObject

	return msg
}
