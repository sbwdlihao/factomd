package state

// This file is for the simulator to attach identities properly to the state.
// Each state has its own set of keys that need to match the ones in the
// identitiy to properly test identities/authorities
import (
	"github.com/FactomProject/factomd/common/primitives"
)

func (s *State) SimSetNewKeys(p *primitives.PrivateKey) {
	s.serverPrivKey = p
	s.serverPubKey = p.Pub
}

func (s *State) SimGetSigKey() string {
	return s.serverPrivKey.Pub.String()
}
