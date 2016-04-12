package adminBlock

import (
	"bytes"
	"fmt"
	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

// DB Signature Entry -------------------------
type AddFederatedServer struct {
	IdentityChainID            interfaces.IHash
}

var _ interfaces.IABEntry = (*AddFederatedServer)(nil)
var _ interfaces.BinaryMarshallable = (*AddFederatedServer)(nil)

func (c *AddFederatedServer) UpdateState(state interfaces.IState) {
    if len(state.GetFedServers())==0 {
        state.AddFedServer(c.IdentityChainID)
    }
    state.Println(fmt.Sprintf("Adding Federaed Server: %x",c.IdentityChainID.Bytes()[:3]))
}

// Create a new DB Signature Entry
func NewAddFederatedServer(identityChainID interfaces.IHash) (e *AddFederatedServer) {
	e = new(AddFederatedServer)
    e.IdentityChainID = primitives.NewHash(identityChainID.Bytes())
	return
}

func (e *AddFederatedServer) Type() byte {
	return constants.TYPE_ADD_FED_SERVER 
}

func (e *AddFederatedServer) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer

	data, err = e.IdentityChainID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	buf.Write(data)
    
	return buf.Bytes(), nil
}

func (e *AddFederatedServer) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling: %v", r)
		}
	}()

	newData = data
	newData =  newData[1:]

    e.IdentityChainID = new(primitives.Hash)
	newData, err = e.IdentityChainID.UnmarshalBinaryData(newData)
	if err != nil {
		return
	}
	return
}

func (e *AddFederatedServer) UnmarshalBinary(data []byte) (err error) {
	_, err = e.UnmarshalBinaryData(data)
	return
}

func (e *AddFederatedServer) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *AddFederatedServer) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *AddFederatedServer) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func (e *AddFederatedServer) String() string {
	str := fmt.Sprintf("Add Server with Identity Chain ID = %x",e.IdentityChainID.Bytes()[:5])
	return str
}

func (e *AddFederatedServer) IsInterpretable() bool {
	return false
}

func (e *AddFederatedServer) Interpret() string {
	return ""
}

func (e *AddFederatedServer) Hash() interfaces.IHash {
	bin, err := e.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return primitives.Sha(bin)
}