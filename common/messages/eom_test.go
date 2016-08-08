// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages_test

import (
	"testing"

	"github.com/FactomProject/factomd/common/constants"
	. "github.com/FactomProject/factomd/common/messages"

	"github.com/FactomProject/factomd/common/primitives"
)

func TestMarshalUnmarshalEOM(t *testing.T) {
	msg := newEOM()

	hex, err := msg.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Marshalled - %x", hex)

	msg2, err := UnmarshalMessage(hex)
	if err != nil {
		t.Error(err)
	}
	str := msg2.String()
	t.Logf("str - %v", str)

	if msg2.Type() != constants.EOM_MSG {
		t.Error("Invalid message type unmarshalled")
	}

	hex2, err := msg2.(*EOM).MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	if len(hex) != len(hex2) {
		t.Error("Hexes aren't of identical length")
	}
	for i := range hex {
		if hex[i] != hex2[i] {
			t.Error("Hexes do not match")
		}
	}

	if msg.IsSameAs(msg2.(*EOM)) != true {
		t.Errorf("EOM messages are not identical")
	}
}

func TestSignAndVerifyEOM(t *testing.T) {
	msg := newSignedEOM()
	hex, err := msg.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Marshalled - %x", hex)

	t.Logf("Sig - %x", *msg.Signature.GetSignature())
	if len(*msg.Signature.GetSignature()) == 0 {
		t.Error("Signature not present")
	}

	valid, err := msg.VerifySignature()
	if err != nil {
		t.Error(err)
	}
	if valid == false {
		t.Error("Signature is not valid")
	}

	msg2, err := UnmarshalMessage(hex)
	if err != nil {
		t.Error(err)
	}

	if msg2.Type() != constants.EOM_MSG {
		t.Error("Invalid message type unmarshalled")
	}

	valid, err = msg2.(*EOM).VerifySignature()
	if err != nil {
		t.Error(err)
	}
	if valid == false {
		t.Error("Signature 2 is not valid")
	}
}

func newEOM() *EOM {
	eom := new(EOM)
	eom.Timestamp = primitives.NewTimestampFromMilliseconds(0xFF22100122FF)
	eom.Minute = 3
	h, err := primitives.NewShaHashFromStr("deadbeef00000000000000000000000000000000000000000000000000000000")
	if err != nil {
		panic(err)
	}
	eom.ChainID = h
	eom.DBHeight = 123456

	return eom
}

func newSignedEOM() *EOM {
	ack := newEOM()

	key, err := primitives.NewPrivateKeyFromHex("07c0d52cb74f4ca3106d80c4a70488426886bccc6ebc10c6bafb37bf8a65f4c38cee85c62a9e48039d4ac294da97943c2001be1539809ea5f54721f0c5477a0a")
	if err != nil {
		panic(err)
	}
	err = ack.Sign(key)
	if err != nil {
		panic(err)
	}

	return ack
}
