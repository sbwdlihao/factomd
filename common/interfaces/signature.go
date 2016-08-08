// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package interfaces

/**************************************
 * ISign
 *
 * Interface for RCB Signatures
 *
 * Data structure to support signatures, including multisig.
 **************************************/

// Verifier objects can Verify signed messages
type Verifier interface {
	Verify(msg []byte, sig *[64]byte) bool
}

// Signer object can Sign msg
type Signer interface {
	Sign(msg []byte) IFullSignature
}

type ISignature interface {
	BinaryMarshallable
	SetSignature(sig []byte) error // Set or update the signature
	GetSignature() *[64]byte
	CustomMarshalText() ([]byte, error)
	Bytes() []byte
}

type IFullSignature interface {
	ISignature
	SetPub(publicKey []byte)
	// Get the public key
	GetKey() []byte
	// Validate data against this signature
	Verify(data []byte) bool
	IsSameAs(IFullSignature) bool
}

/**************************************
 * ISign
 *
 * Interface for RCB Signatures
 *
 * The signature block holds the signatures that validate one of the RCBs.
 * Each signature has an index, so if the RCD is a multisig, you can know
 * how to apply the signatures to the addresses in the RCD.
 **************************************/
type ISignatureBlock interface {
	IBlock
	GetSignatures() []ISignature
	AddSignature(sig ISignature)
	GetSignature(int) ISignature
}

type ISignable interface {
	Sign(privateKey []byte) error
	MarshalBinarySig() ([]byte, error)
	ValidateSignatures() error
}
