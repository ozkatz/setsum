package setsum_test

import (
	"testing"

	"github.com/ozkatz/setsum"
	"gotest.tools/v3/assert"
)

var SEVEN_VALUES = [32]byte{
	197, 179, 253, 77, 1, 242, 184, 4, 15, 84, 171, 116, 18, 202, 83, 187, 252, 153, 14, 39,
	42, 64, 173, 209, 196, 206, 186, 107, 47, 228, 114, 213,
}

func TestConstants(t *testing.T) {
	assert.Equal(t, setsum.SETSUM_BYTES, setsum.SETSUM_BYTES_PER_COLUMN*setsum.SETSUM_COLUMNS)
}

func TestAddState(t *testing.T) {
	lhs := [setsum.SETSUM_COLUMNS]uint32{1, 2, 3, 4, 5, 6, 7, 8}
	rhs := [setsum.SETSUM_COLUMNS]uint32{2, 4, 6, 8, 10, 12, 14, 16}
	expected := [setsum.SETSUM_COLUMNS]uint32{3, 6, 9, 12, 15, 18, 21, 24}
	returned := setsum.AddState(lhs[:], rhs[:])
	assert.Equal(t, expected, returned)
}

func TestAddStateExactlyPrimes(t *testing.T) {
	lhs := [setsum.SETSUM_COLUMNS]uint32{
		3146800025, 1792545563, 417324692, 3444237760, 2812742746, 1608771649, 1661742866,
		3220115897,
	}
	rhs := [setsum.SETSUM_COLUMNS]uint32{
		1148167266, 2502421716, 3877642539, 850729437, 1482224443, 2686195512, 2633224277,
		1074851214,
	}
	expected := [setsum.SETSUM_COLUMNS]uint32{0}
	returned := setsum.AddState(lhs[:], rhs[:])
	assert.Equal(t, expected, returned)
}

func TestInvertStateDescending(t *testing.T) {
	stateIn := [setsum.SETSUM_COLUMNS]uint32{
		0xffffeeee, 0xddddcccc, 0xbbbbaaaa, 0x99998888, 0x77776666, 0x66665555, 0x44443333,
		0x22221111,
	}
	expected := [setsum.SETSUM_COLUMNS]uint32{
		4365, 572666659, 1145328917, 1717991189, 2290653487, 2576984612, 3149646900, 3722309174,
	}
	returned := setsum.InvertState(stateIn)
	assert.Equal(t, expected, returned)
}

func TestInsert7(t *testing.T) {

	t.Run("Sorted", func(t *testing.T) {
		s := setsum.Default()
		s.Insert([]byte("this is the first value"))
		s.Insert([]byte("this is the second value"))
		s.Insert([]byte("this is the third value"))
		s.Insert([]byte("this is the fourth value"))
		s.Insert([]byte("this is the fifth value"))
		s.Insert([]byte("this is the sixth value"))
		s.Insert([]byte("this is the seventh value"))
		digest := s.Digest()
		assert.Equal(t, SEVEN_VALUES, digest)
	})

	t.Run("Reversed", func(t *testing.T) {
		s := setsum.Default()
		s.Insert([]byte("this is the seventh value"))
		s.Insert([]byte("this is the sixth value"))
		s.Insert([]byte("this is the fifth value"))
		s.Insert([]byte("this is the fourth value"))
		s.Insert([]byte("this is the third value"))
		s.Insert([]byte("this is the second value"))
		s.Insert([]byte("this is the first value"))
		digest := s.Digest()
		assert.Equal(t, SEVEN_VALUES, digest)
	})

	t.Run("Random", func(t *testing.T) {
		s := setsum.Default()
		s.Insert([]byte("this is the fifth value"))
		s.Insert([]byte("this is the fourth value"))
		s.Insert([]byte("this is the third value"))
		s.Insert([]byte("this is the sixth value"))
		s.Insert([]byte("this is the seventh value"))
		s.Insert([]byte("this is the second value"))
		s.Insert([]byte("this is the first value"))
		digest := s.Digest()
		assert.Equal(t, SEVEN_VALUES, digest)
	})

	t.Run("MergeTwoSets", func(t *testing.T) {
		s1 := setsum.Default()
		s1.Insert([]byte("this is the first value"))
		s1.Insert([]byte("this is the second value"))
		s1.Insert([]byte("this is the third value"))
		s1.Insert([]byte("this is the fourth value"))

		s2 := setsum.Default()
		s2.Insert([]byte("this is the fifth value"))
		s2.Insert([]byte("this is the sixth value"))
		s2.Insert([]byte("this is the seventh value"))

		s12 := s1.Merge(s2)
		digest := s12.Digest()
		assert.Equal(t, SEVEN_VALUES, digest)
	})
}

func TestInsertRemove(t *testing.T) {
	s := setsum.Default()
	s.Insert([]byte("this is the first value"))
	s.Insert([]byte("this is the second value"))
	s.Insert([]byte("this is the third value"))
	s.Insert([]byte("this is the fourth value"))
	s.Insert([]byte("this is the fifth value"))
	s.Insert([]byte("this is the sixth value"))
	s.Insert([]byte("this is the seventh value"))
	s.Remove([]byte("this is the seventh value"))
	s.Remove([]byte("this is the sixth value"))
	s.Remove([]byte("this is the fifth value"))
	s.Remove([]byte("this is the fourth value"))
	s.Remove([]byte("this is the third value"))
	s.Remove([]byte("this is the second value"))
	s.Remove([]byte("this is the first value"))
	digest := s.Digest()
	assert.Equal(t, setsum.Default().Digest(), digest)
}

func TestRemoveTwoSets(t *testing.T) {
	s := setsum.Default()
	s.Insert([]byte("this is the first value"))
	s.Insert([]byte("this is the second value"))
	s.Insert([]byte("this is the third value"))
	s.Insert([]byte("this is the fourth value"))
	s.Insert([]byte("this is the fifth value"))
	s.Insert([]byte("this is the sixth value"))
	s.Insert([]byte("this is the seventh value"))

	s1 := setsum.Default()
	s1.Insert([]byte("this is the first value"))
	s1.Insert([]byte("this is the second value"))
	s1.Insert([]byte("this is the third value"))
	s1.Insert([]byte("this is the fourth value"))

	s2 := setsum.Default()
	s2.Insert([]byte("this is the fifth value"))
	s2.Insert([]byte("this is the sixth value"))
	s2.Insert([]byte("this is the seventh value"))

	s12 := s.Subtract(s1).Subtract(s2)
	digest := s12.Digest()
	assert.Equal(t, setsum.Default().Digest(), digest)
}

func TestFromDigest(t *testing.T) {
	s := setsum.Default()
	s.Insert([]byte("this is the first value"))
	s.Insert([]byte("this is the second value"))
	s.Insert([]byte("this is the third value"))
	s.Insert([]byte("this is the fourth value"))
	s.Insert([]byte("this is the fifth value"))
	s.Insert([]byte("this is the sixth value"))
	s.Insert([]byte("this is the seventh value"))
	assert.Assert(t, setsum.FromDigest(SEVEN_VALUES[:]).Equals(s))
}

func TestFromHexDigest(t *testing.T) {
	const SEVEN_HEX_VALUES = "c5b3fd4d01f2b8040f54ab7412ca53bbfc990e272a40add1c4ceba6b2fe472d5"
	assert.Assert(t, setsum.FromDigest(SEVEN_VALUES[:]).Equals(setsum.FromHexDigest(SEVEN_HEX_VALUES)))
}

func TestFromHexDigestInvalid(t *testing.T) {
	assert.Assert(t, setsum.FromHexDigest("invalid") == nil)
}

func TestFromDigestInvalid(t *testing.T) {
	assert.Assert(t, setsum.FromDigest([]byte("invalid")) == nil)
}

func TestFromDigestInvalidLength(t *testing.T) {
	assert.Assert(t, setsum.FromDigest([]byte("invalid")) == nil)
}

func BenchmarkInsert(b *testing.B) {
	s := setsum.Default()
	values := []string{
		"this is the first value",
		"this is the second value",
		"this is the third value",
		"this is the fourth value",
		"this is the fifth value",
		"this is the sixth value",
		"this is the seventh value",
	}
	for i := 0; i < b.N; i++ {
		s.Insert([]byte(values[i%len(values)]))
	}
}

func BenchmarkInsertVector(b *testing.B) {
	s := setsum.Default()
	values := []string{
		"this is the first value",
		"this is the second value",
		"this is the third value",
		"this is the fourth value",
		"this is the fifth value",
		"this is the sixth value",
		"this is the seventh value",
	}
	for i := 0; i < b.N; i++ {
		s.InsertVector([][]byte{[]byte(values[i%len(values)])})
	}
}
