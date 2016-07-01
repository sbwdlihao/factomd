// Copyright 2015 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package directoryBlock

import (
	"bytes"
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

var _ = fmt.Print

type DirectoryBlock struct {
	//Marshalized
	Header    interfaces.IDirectoryBlockHeader
	DBEntries []interfaces.IDBEntry

	//Not Marshalized
	DBHash interfaces.IHash
	KeyMR  interfaces.IHash
}

var _ interfaces.Printable = (*DirectoryBlock)(nil)
var _ interfaces.BinaryMarshallableAndCopyable = (*DirectoryBlock)(nil)
var _ interfaces.IDirectoryBlock = (*DirectoryBlock)(nil)
var _ interfaces.DatabaseBatchable = (*DirectoryBlock)(nil)
var _ interfaces.DatabaseBlockWithEntries = (*DirectoryBlock)(nil)

func (c *DirectoryBlock) SetEntryHash(hash, chainID interfaces.IHash, index int) {
	if len(c.DBEntries) < index {
		ent := make([]interfaces.IDBEntry, index)
		copy(ent, c.DBEntries)
		c.DBEntries = ent
	}
	dbe := new(DBEntry)
	dbe.ChainID = chainID
	dbe.KeyMR = hash
	c.DBEntries[index] = dbe
}

func (c *DirectoryBlock) SetABlockHash(aBlock interfaces.IAdminBlock) error {
	hash, err := aBlock.PartialHash()
	if err != nil {
		return err
	}
	fmt.Println("Justin zzzzzzzzzzzzzzzzzzzzzzzzzzzz SetABlockHash", hash.String())
	fmt.Println("Ab:", aBlock.String())
	c.SetEntryHash(hash, aBlock.GetChainID(), 0)
	return nil
}

func (c *DirectoryBlock) SetECBlockHash(ecBlock interfaces.IEntryCreditBlock) error {
	hash, err := ecBlock.HeaderHash()
	if err != nil {
		return err
	}

	fmt.Println("Justin zzzzzzzzzzzzzzzzzzzzzzzzzzzz SetECBlockHash", hash.String())
	fmt.Println("Eb:", ecBlock.String())

	c.SetEntryHash(hash, ecBlock.GetChainID(), 1)
	return nil
}

func (c *DirectoryBlock) SetFBlockHash(fBlock interfaces.IFBlock) error {
	hash := fBlock.GetKeyMR()
	fmt.Println("Justin zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz SetFBlockHash", hash.String())
	fmt.Println("Fb:", fBlock.String())

	c.SetEntryHash(hash, fBlock.GetChainID(), 2)
	return nil
}

func (c *DirectoryBlock) GetEntryHashes() []interfaces.IHash {
	entries := c.DBEntries[:]
	answer := make([]interfaces.IHash, len(entries))
	for i, entry := range entries {
		answer[i] = entry.GetKeyMR()
	}
	return answer
}

func (c *DirectoryBlock) GetEntrySigHashes() []interfaces.IHash {
	return nil
}

func (c *DirectoryBlock) Sort() {
	done := false
	for i := 3; !done && i < len(c.DBEntries)-1; i++ {
		done = true
		for j := 3; j < len(c.DBEntries)-1-i+3; j++ {
			comp := bytes.Compare(c.DBEntries[j].GetChainID().Bytes(),
				c.DBEntries[j+1].GetChainID().Bytes())
			if comp > 0 {
				h := c.DBEntries[j]
				c.DBEntries[j] = c.DBEntries[j+1]
				c.DBEntries[j+1] = h
			}
			if comp != 0 {
				done = false
			}
		}
	}
}

func (c *DirectoryBlock) GetEntryHashesForBranch() []interfaces.IHash {
	entries := c.DBEntries[:]
	answer := make([]interfaces.IHash, 2*len(entries))
	for i, entry := range entries {
		answer[2*i] = entry.GetChainID()
		answer[2*i+1] = entry.GetKeyMR()
	}
	return answer
}

func (c *DirectoryBlock) GetDBEntries() []interfaces.IDBEntry {
	return c.DBEntries
}

func (c *DirectoryBlock) GetKeyMR() interfaces.IHash {
	fmt.Println("Justin aaaaaaaaaaaaaaaaaaaaaaaaaaa: GetKeyMR()")
	keyMR, err := c.BuildKeyMerkleRoot()
	if err != nil {
		panic("Failed to build the key MR")
	}

	c.KeyMR = keyMR

	return c.KeyMR
}

func (c *DirectoryBlock) GetHeader() interfaces.IDirectoryBlockHeader {
	return c.Header
}

func (c *DirectoryBlock) SetHeader(header interfaces.IDirectoryBlockHeader) {
	c.Header = header
}

func (c *DirectoryBlock) SetDBEntries(dbEntries []interfaces.IDBEntry) error {
	fmt.Println("Justin aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa: SetDBEntries")
	c.DBEntries = dbEntries
	c.GetHeader().SetBlockCount(uint32(len(dbEntries)))
	_, err := c.BuildBodyMR()
	if err != nil {
		return err
	}
	return nil
}

func (c *DirectoryBlock) New() interfaces.BinaryMarshallableAndCopyable {
	dBlock := new(DirectoryBlock)
	dBlock.Header = NewDBlockHeader()
	dBlock.DBHash = primitives.NewZeroHash()
	dBlock.KeyMR = primitives.NewZeroHash()
	return dBlock
}

func (c *DirectoryBlock) GetDatabaseHeight() uint32 {
	return c.GetHeader().GetDBHeight()
}

func (c *DirectoryBlock) GetChainID() interfaces.IHash {
	return primitives.NewHash(constants.D_CHAINID)
}

func (c *DirectoryBlock) DatabasePrimaryIndex() interfaces.IHash {
	return c.GetKeyMR()
}

func (c *DirectoryBlock) DatabaseSecondaryIndex() interfaces.IHash {
	return c.GetHash()
}

func (e *DirectoryBlock) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *DirectoryBlock) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}

func (e *DirectoryBlock) JSONBuffer(b *bytes.Buffer) error {
	return primitives.EncodeJSONToBuffer(e, b)
}

func (e *DirectoryBlock) String() string {
	var out primitives.Buffer
	fmt.Println("Justin aaaaaaaaaaaaaaaaaaaabbb String()")
	kmr, err := e.BuildKeyMerkleRoot()

	if err != nil {
		out.WriteString(fmt.Sprintf("%20s %v\n", "KeyMR:", err))
	} else {
		out.WriteString(fmt.Sprintf("%20s %v\n", "KeyMR:", kmr.String()))
	}

	fmt.Println("Justin aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa: String()")
	kmr, err = e.BuildBodyMR()
	if err != nil {
		out.WriteString(fmt.Sprintf("%20s %v\n", "BodyMR:", err))
	} else {
		out.WriteString(fmt.Sprintf("%20s %v\n", "BodyMR:", kmr.String()))
	}

	fh := e.GetFullHash()
	out.WriteString(fmt.Sprintf("%20s %v\n", "BodyMR:", fh.String()))

	out.WriteString(e.Header.String())
	out.WriteString("Entries: \n")
	for _, entry := range e.DBEntries {
		out.WriteString(entry.String())
	}

	return (string)(out.DeepCopyBytes())

}

func (b *DirectoryBlock) MarshalBinary() (data []byte, err error) {
	var buf primitives.Buffer

	b.Sort()

	fmt.Println("Justin aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa: MarshalBinary")
	b.BuildBodyMR()

	count := uint32(len(b.GetDBEntries()))
	b.GetHeader().SetBlockCount(count)

	data, err = b.GetHeader().MarshalBinary()
	if err != nil {
		return
	}
	buf.Write(data)

	for i := uint32(0); i < count; i++ {
		data, err = b.GetDBEntries()[i].MarshalBinary()
		if err != nil {
			return
		}
		buf.Write(data)
	}

	return buf.DeepCopyBytes(), err
}

func (b *DirectoryBlock) BuildBodyMR() (interfaces.IHash, error) {
	fmt.Println("kkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkkk")
	fmt.Println("Prev:", b.Header.GetPrevKeyMR().String())

	hashes := make([]interfaces.IHash, len(b.GetDBEntries()))
	for i, entry := range b.GetDBEntries() {
		fmt.Println("kk:", entry.String())
		data, err := entry.MarshalBinary()
		if err != nil {
			return nil, err
		}
		hashes[i] = primitives.Sha(data)
	}

	if len(hashes) == 0 {
		hashes = append(hashes, primitives.Sha(nil))
	}

	merkleTree := primitives.BuildMerkleTreeStore(hashes)
	merkleRoot := merkleTree[len(merkleTree)-1]

	b.GetHeader().SetBodyMR(merkleRoot)
	fmt.Println("MR:", merkleRoot.String())
	fmt.Println("")

	return merkleRoot, nil
}

func (b *DirectoryBlock) HeaderHash() (interfaces.IHash, error) {
	binaryEBHeader, err := b.GetHeader().MarshalBinary()
	if err != nil {
		return nil, err
	}
	return primitives.Sha(binaryEBHeader), nil
}

func (b *DirectoryBlock) BodyKeyMR() interfaces.IHash {
	fmt.Println("Justin aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa: BodyKeyMR")
	key, _ := b.BuildBodyMR()
	return key
}

func (b *DirectoryBlock) BuildKeyMerkleRoot() (keyMR interfaces.IHash, err error) {
	// Create the Entry Block Key Merkle Root from the hash of Header and the Body Merkle Root

	hashes := make([]interfaces.IHash, 0, 2)
	headerHash, err := b.HeaderHash()
	if err != nil {
		return nil, err
	}
	hashes = append(hashes, headerHash)
	fmt.Println("Justin aaaaaaaaaaaaaaaaaaaa BuildKeyMerkleRoot")
	hashes = append(hashes, b.BodyKeyMR())
	merkle := primitives.BuildMerkleTreeStore(hashes)
	keyMR = merkle[len(merkle)-1] // MerkleRoot is not marshalized in Dir Block

	b.KeyMR = keyMR

	b.GetFullHash() // Create the Full Hash when we create the keyMR

	return primitives.NewHash(keyMR.Bytes()), nil
}

func UnmarshalDBlock(data []byte) (interfaces.IDirectoryBlock, error) {
	dBlock := new(DirectoryBlock)
	dBlock.Header = NewDBlockHeader()
	dBlock.DBHash = primitives.NewZeroHash()
	dBlock.KeyMR = primitives.NewZeroHash()
	err := dBlock.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}
	return dBlock, nil
}

func (b *DirectoryBlock) UnmarshalBinaryData(data []byte) (newData []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error unmarshalling Directory Block: %v", r)
		}
	}()

	newData = data

	var fbh interfaces.IDirectoryBlockHeader = new(DBlockHeader)

	newData, err = fbh.UnmarshalBinaryData(newData)
	if err != nil {
		return
	}
	b.SetHeader(fbh)

	count := b.GetHeader().GetBlockCount()
	entries := make([]interfaces.IDBEntry, count)
	for i := uint32(0); i < count; i++ {
		entries[i] = new(DBEntry)
		newData, err = entries[i].UnmarshalBinaryData(newData)
		if err != nil {
			return
		}
	}

	fmt.Println("ATUM")
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
	fmt.Println("Justin cccccccccccccc: UnmarshalBinary")
	for aa, ent := range entries {
		fmt.Println(aa, "::ijij::", ent.String())
	}
	err = b.SetDBEntries(entries)
	if err != nil {
		return
	}

	return
}

func (h *DirectoryBlock) GetTimestamp() interfaces.Timestamp {
	return h.GetHeader().GetTimestamp()
}

func (b *DirectoryBlock) UnmarshalBinary(data []byte) (err error) {
	_, err = b.UnmarshalBinaryData(data)
	return
}

func (b *DirectoryBlock) GetHash() interfaces.IHash {
	return b.GetFullHash()
}

func (b *DirectoryBlock) GetFullHash() interfaces.IHash {
	binaryDblock, err := b.MarshalBinary()
	if err != nil {
		return nil
	}
	b.DBHash = primitives.Sha(binaryDblock)
	return b.DBHash
}

func (b *DirectoryBlock) AddEntry(chainID interfaces.IHash, keyMR interfaces.IHash) error {
	var dbentry interfaces.IDBEntry
	dbentry = new(DBEntry)
	dbentry.SetChainID(chainID)
	dbentry.SetKeyMR(keyMR)

	if b.DBEntries == nil {
		b.DBEntries = []interfaces.IDBEntry{}
	}

	fmt.Println("Just added: ", keyMR.String())
	return b.SetDBEntries(append(b.DBEntries, dbentry))
}

/*********************************************************************
 * Support
 *********************************************************************/

func NewDirectoryBlock(dbheight uint32, prev *DirectoryBlock) interfaces.IDirectoryBlock {
	newdb := new(DirectoryBlock)

	newdb.Header = new(DBlockHeader)
	newdb.Header.SetVersion(constants.VERSION_0)
	newdb.Header.SetPrevFullHash(primitives.NewZeroHash())
	newdb.Header.SetPrevKeyMR(primitives.NewZeroHash())

	if prev != nil {
		newdb.GetHeader().SetPrevFullHash(prev.GetFullHash())
		newdb.GetHeader().SetPrevKeyMR(prev.GetKeyMR())
	}

	newdb.GetHeader().SetDBHeight(dbheight)
	newdb.SetDBEntries(make([]interfaces.IDBEntry, 0))

	abChainID := primitives.NewHash(constants.ADMIN_CHAINID)
	newABHash := primitives.NewZeroHash()

	ecChainID := primitives.NewHash(constants.EC_CHAINID)
	newECHash := primitives.NewZeroHash()

	fbChainID := primitives.NewHash(constants.FACTOID_CHAINID)
	newFBHash := primitives.NewZeroHash()

	newdb.AddEntry(abChainID, newABHash)
	newdb.AddEntry(ecChainID, newECHash)
	newdb.AddEntry(fbChainID, newFBHash)

	fmt.Println("NEWA: ", abChainID.String(), newABHash.String())
	fmt.Println("NEWE:", ecChainID.String(), newECHash.String())
	fmt.Println("NEWF:", fbChainID.String(), newFBHash.String())

	return newdb
}
