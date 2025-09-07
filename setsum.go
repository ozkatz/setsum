package setsum

import (
	"encoding/binary"
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

const (
	// SETSUM_BYTES is the number of bytes in the digest of both the hash used by setsum and the output
	// of setsum.
	SETSUM_BYTES uint = 32

	// SETSUM_BYTES_PER_COLUMN is the number of bytes per column.
	// This should evenly divide the number of bytes.  This number is
	// implicitly wound through the code in its use of uint32 to store columns as it's the number of bytes
	// used to store a uint32.
	SETSUM_BYTES_PER_COLUMN uint = 4

	// SETSUM_COLUMNS is the number of columns in the logical/internal representation of the setsum.
	SETSUM_COLUMNS uint = SETSUM_BYTES / SETSUM_BYTES_PER_COLUMN
)

// SETSUM_PRIMES is an array of primes used to construct a field of different size and transformations
// for each column.
var SETSUM_PRIMES = [SETSUM_COLUMNS]uint32{
	4294967291, 4294967279, 4294967231, 4294967197,
	4294967189, 4294967161, 4294967143, 4294967111,
}

func AddState(lhs, rhs []uint32) [SETSUM_COLUMNS]uint32 {
	var ret [SETSUM_COLUMNS]uint32
	for i := 0; i < int(SETSUM_COLUMNS); i++ {
		lc := uint64(lhs[i])
		rc := uint64(rhs[i])
		sum := lc + rc
		p := uint64(SETSUM_PRIMES[i])
		if sum >= p {
			sum -= p
		}
		ret[i] = uint32(sum)
	}
	return ret
}

func InvertState(state [SETSUM_COLUMNS]uint32) [SETSUM_COLUMNS]uint32 {
	var ret [SETSUM_COLUMNS]uint32
	for i := 0; i < int(SETSUM_COLUMNS); i++ {
		ret[i] = SETSUM_PRIMES[i] - state[i]
	}
	return ret
}

func itemVectoredToState(item [][]byte) [SETSUM_COLUMNS]uint32 {
	h := sha3.New256()
	for _, piece := range item {
		_, _ = h.Write(piece)
	}
	hashBytes := h.Sum(nil)
	return hashToState(hashBytes)
}

func hashToState(hashBytes []byte) [SETSUM_COLUMNS]uint32 {
	var ret [SETSUM_COLUMNS]uint32
	for i := 0; i < int(SETSUM_COLUMNS); i++ {
		offset := i * int(SETSUM_BYTES_PER_COLUMN)
		val := binary.LittleEndian.Uint32(hashBytes[offset : offset+int(SETSUM_BYTES_PER_COLUMN)])
		v := uint64(val)
		p := uint64(SETSUM_PRIMES[i])
		if v >= p {
			v = v % p
		}
		ret[i] = uint32(v)
	}
	return ret
}

type Setsum struct {
	state [SETSUM_COLUMNS]uint32
}

func Default() *Setsum {
	return &Setsum{
		state: [SETSUM_COLUMNS]uint32{0},
	}
}

func (s *Setsum) Insert(item []byte) {
	s.InsertVector([][]byte{item})
}

func (s *Setsum) InsertVector(item [][]byte) {
	itemState := itemVectoredToState(item)
	newState := AddState(s.state[:], itemState[:])
	s.state = newState
}

func (s *Setsum) Remove(item []byte) {
	s.RemoveVector([][]byte{item})
}

func (s *Setsum) RemoveVector(item [][]byte) {
	itemState := itemVectoredToState(item)
	inv := InvertState(itemState)
	newState := AddState(s.state[:], inv[:])
	s.state = newState
}

func (s *Setsum) Digest() [SETSUM_BYTES]byte {
	var itemHash [SETSUM_BYTES]byte
	for col := 0; col < int(SETSUM_COLUMNS); col++ {
		idx := col * int(SETSUM_BYTES_PER_COLUMN)
		binary.LittleEndian.PutUint32(itemHash[idx:idx+int(SETSUM_BYTES_PER_COLUMN)], s.state[col])
	}
	return itemHash
}

func (s *Setsum) Merge(other *Setsum) *Setsum {
	state := AddState(s.state[:], other.state[:])
	return &Setsum{state: state}
}

func (s *Setsum) Subtract(other *Setsum) *Setsum {
	rhsState := InvertState(other.state)
	state := AddState(s.state[:], rhsState[:])
	return &Setsum{state: state}
}

func (s *Setsum) HexDigest() string {
	d := s.Digest()
	return hex.EncodeToString(d[:])
}

func FromDigest(digest []byte) *Setsum {
	if len(digest) != int(SETSUM_BYTES) {
		return nil
	}
	var state [SETSUM_COLUMNS]uint32
	for col := 0; col < int(SETSUM_COLUMNS); col++ {
		idx := col * int(SETSUM_BYTES_PER_COLUMN)
		state[col] = binary.LittleEndian.Uint32(digest[idx : idx+int(SETSUM_BYTES_PER_COLUMN)])
	}
	return &Setsum{state: state}
}

func FromHexDigest(hexDigest string) *Setsum {
	if len(hexDigest) != int(SETSUM_BYTES)*2 {
		return nil
	}
	b, err := hex.DecodeString(hexDigest)
	if err != nil || len(b) != int(SETSUM_BYTES) {
		return nil
	}
	return FromDigest(b)
}

func (s *Setsum) String() string {
	return s.HexDigest()
}

func (s *Setsum) Equals(other *Setsum) bool {
	return s.Digest() == other.Digest() && s.state == other.state
}
