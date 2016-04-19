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

//Structure to request missing messages in a node's process list
type MissingData struct {
	MessageBase
	RequestHash interfaces.IHash
	Timestamp   interfaces.Timestamp

	//Not marshalled
	hash interfaces.IHash
}

var _ interfaces.IMsg = (*MissingData)(nil)

func (m *MissingData) Process(uint32, interfaces.IState) bool {
	return true
}

func (m *MissingData) GetHash() interfaces.IHash {
	if m.hash == nil {
		data, err := m.MarshalBinary()
		if err != nil {
			panic(fmt.Sprintf("Error in MissingData.GetHash(): %s", err.Error()))
		}
		m.hash = primitives.Sha(data)
	}
	return m.hash
}

func (m *MissingData) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
		data, err := m.MarshalBinary()
		if err != nil {
			return nil
		}
		m.MsgHash = primitives.Sha(data)
	}
	return m.MsgHash
}

func (m *MissingData) GetTimestamp() interfaces.Timestamp {
	return m.Timestamp
}

func (m *MissingData) Type() int {
	return constants.MISSING_DATA
}

func (m *MissingData) Int() int {
	return -1
}

func (m *MissingData) Bytes() []byte {
	return nil
}

func (m *MissingData) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
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

	m.RequestHash = primitives.NewHash(constants.ZERO_HASH)
	newData, err = m.RequestHash.UnmarshalBinaryData(newData)
	if err != nil {
		return nil, err
	}

	m.Peer2peer = true // Always a peer2peer request.

	return data, nil
}

func (m *MissingData) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *MissingData) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, byte(m.Type()))

	t := m.GetTimestamp()
	data, err := t.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	//TODO: actually marshal RequestHash properly
	//binary.Write(&buf, binary.BigEndian, m.RequestHash)

	if d, err := m.RequestHash.MarshalBinary(); err != nil {
		return nil, err
	} else {
		buf.Write(d)
	}

	var md MissingData

	bb := buf.Bytes()

	if unmarshalErr := md.UnmarshalBinary(bb); unmarshalErr != nil {
		return nil, unmarshalErr
	}

	return bb, nil
}

func (m *MissingData) String() string {
	return fmt.Sprintf("MissingData: %+v", m.RequestHash)
}

// Validate the message, given the state.  Three possible results:
//  < 0 -- Message is invalid.  Discard
//  0   -- Cannot tell if message is Valid
//  1   -- Message is valid
func (m *MissingData) Validate(state interfaces.IState) int {
	return 1
}

// Returns true if this is a message for this server to execute as
// a leader.
func (m *MissingData) Leader(state interfaces.IState) bool {
	return false
}

// Execute the leader functions of the given message
func (m *MissingData) LeaderExecute(state interfaces.IState) error {
	return nil
}

// Returns true if this is a message for this server to execute as a follower
func (m *MissingData) Follower(interfaces.IState) bool {
	return true
}

func (m *MissingData) FollowerExecute(state interfaces.IState) error {
	var dataObject interface{}
	var dataHash interfaces.IHash
	rawObject, dataType, err := state.LoadDataByHash(m.RequestHash)

	if rawObject != nil && err == nil { // If I don't have this message, ignore.
		//msg.SetOrigin(m.GetOrigin())
		//state.NetworkOutMsgQueue() <- msg
		//TODO: actually send the data to the requesting node
		fmt.Println(dataType, "MESSAGE FROM MD: ", rawObject)

		switch dataType {
		case "dblock":
			dataObject = rawObject.(interfaces.IDirectoryBlock)
			dataHash = dataObject.(interfaces.IDirectoryBlock).GetHash()
		case "entry":
			dataObject = rawObject.(interfaces.IEntry)
			dataHash = dataObject.(interfaces.IEntry).GetHash()
		case "eblock":
			dataObject = rawObject.(interfaces.IEntryBlock)
			dataHash, _ = dataObject.(interfaces.IEntryBlock).Hash()
		case "fblock":
			dataObject = rawObject.(interfaces.IFBlock)
			dataHash = dataObject.(interfaces.IFBlock).GetHash()
		case "ecblock":
			dataObject = rawObject.(interfaces.IEntryCreditBlock)
			dataHash = dataObject.(interfaces.IEntryCreditBlock).GetHash()
		case "ablock":
			dataObject = rawObject.(interfaces.IAdminBlock)
			dataHash = dataObject.(interfaces.IAdminBlock).GetHash()
		default:
			return fmt.Errorf("Datatype unsupported")
		}

		msg := NewDataResponse(dataObject, dataType, dataHash)
		msg.SetOrigin(m.GetOrigin())
		state.NetworkOutMsgQueue() <- msg
	} else {
		return err
	}

	return nil
}

func (e *MissingData) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *MissingData) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *MissingData) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func NewMissingData(state interfaces.IState, requestHash interfaces.IHash) interfaces.IMsg {

	msg := new(MissingData)

	msg.Peer2peer = true // Always a peer2peer request.
	msg.Timestamp = state.GetTimestamp()
	msg.RequestHash = requestHash

	return msg
}
