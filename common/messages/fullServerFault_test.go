// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages_test

import (
	"testing"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	. "github.com/FactomProject/factomd/common/messages"
	"github.com/FactomProject/factomd/common/primitives"
)

func TestMarshalUnmarshalFullServerFault(t *testing.T) {
	ts := primitives.NewTimestampNow()
	vmIndex := int(*ts) % 10
	sf := NewServerFault(ts, primitives.NewHash([]byte("a test")), vmIndex, 10, 100)

	sl := coupleOfSigs(t)

	fsf := NewFullServerFault(sf, sl)
	hex, err := fsf.MarshalBinary()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Marshalled - %x", hex)

	fsf2, err := UnmarshalMessage(hex)
	if err != nil {
		t.Error(err)
	}
	str := fsf2.String()
	t.Logf("str - %v", str)

	if fsf2.Type() != constants.FULL_SERVER_FAULT_MSG {
		t.Errorf("Invalid message type unmarshalled - got %v, expected %v", fsf2.Type(), constants.FULL_SERVER_FAULT_MSG)
	}

	hex2, err := fsf2.(*FullServerFault).MarshalBinary()
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
}

func coupleOfSigs(t *testing.T) []interfaces.IFullSignature {
	priv1 := new(primitives.PrivateKey)

	err := priv1.GenerateKey()
	if err != nil {
		t.Fatalf("%v", err)
	}

	msg1 := "Test Message Sign1"
	msg2 := "Test Message Sign2"

	sig1 := priv1.Sign([]byte(msg1))
	sig2 := priv1.Sign([]byte(msg2))

	var twoSigs []interfaces.IFullSignature
	twoSigs = append(twoSigs, sig1)
	twoSigs = append(twoSigs, sig2)
	return twoSigs
}

func makeSigList(t *testing.T) SigList {
	priv1 := new(primitives.PrivateKey)

	err := priv1.GenerateKey()
	if err != nil {
		t.Fatalf("%v", err)
	}

	msg1 := "Test Message Sign1"
	msg2 := "Test Message Sign2"

	sig1 := priv1.Sign([]byte(msg1))
	sig2 := priv1.Sign([]byte(msg2))

	var twoSigs []interfaces.IFullSignature
	twoSigs = append(twoSigs, sig1)
	twoSigs = append(twoSigs, sig2)

	sl := new(SigList)
	sl.Length = 2
	sl.List = twoSigs
	return *sl
}

func TestThatFullAndFaultCoreHashesMatch(t *testing.T) {
	ts := primitives.NewTimestampNow()
	vmIndex := int(*ts) % 10
	sf := NewServerFault(ts, primitives.NewHash([]byte("a test")), vmIndex, 10, 100)

	sl := coupleOfSigs(t)

	fsf := NewFullServerFault(sf, sl)

	if !sf.GetCoreHash().IsSameAs(fsf.GetCoreHash()) {
		t.Error("CoreHashes do not match between FullServerFault and ServerFault")
	}
}
