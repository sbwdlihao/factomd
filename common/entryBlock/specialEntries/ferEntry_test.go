package specialEntries_test

import (
	"testing"

	//"github.com/FactomProject/factomd/common/interfaces"
	//"github.com/FactomProject/factomd/common/primitives"
	. "github.com/FactomProject/factomd/common/entryBlock/specialEntries"
)

func TestMarshalUnmarshalFEREntry(t *testing.T) {
	fe := new(FEREntry)
	fe.Version = "1"
	fe.ExpirationHeight = 1
	fe.ResidentHeight = 2
	fe.TargetActivationHeight = 3
	fe.Priority = 4
	fe.TargetPrice = 5

	b, err := fe.MarshalBinary()
	if err != nil {
		t.Errorf("%v", err)
	}

	fe2 := new(FEREntry)
	err = fe2.UnmarshalBinary(b)
	if err != nil {
		t.Errorf("%v", err)
	}

	if fe.Version != fe2.Version {
		t.Error("Version is not the same")
	}
	if fe.ExpirationHeight != fe2.ExpirationHeight {
		t.Error("ExpirationHeight is not the same")
	}
	if fe.ResidentHeight != fe2.ResidentHeight {
		t.Error("ResidentHeight is not the same")
	}
	if fe.TargetActivationHeight != fe2.TargetActivationHeight {
		t.Error("TargetActivationHeight is not the same")
	}
	if fe.Priority != fe2.Priority {
		t.Error("Priority is not the same")
	}
	if fe.TargetPrice != fe2.TargetPrice {
		t.Error("TargetPrice is not the same")
	}
}
