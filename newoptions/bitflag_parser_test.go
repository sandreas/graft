package newoptions

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
const (
	FLAG_ONE BitFlag = 1 << iota
	FLAG_TWO
	FLAG_THREE
	FLAG_FOUR
)


func TestParseSingleOrConcatBitFlagParam(t *testing.T) {
	expect := assert.New(t)

	subject := NewBitFlagParser(FLAG_ONE|FLAG_FOUR)

	expect.True(subject.HasFlag(FLAG_ONE))
	expect.False(subject.HasFlag(FLAG_TWO))
	expect.False(subject.HasFlag(FLAG_THREE))
	expect.True(subject.HasFlag(FLAG_FOUR))
}


func TestParseMultipleBitFlagParams(t *testing.T) {
	expect := assert.New(t)

	subject := NewBitFlagParser(FLAG_ONE, FLAG_FOUR)

	expect.True(subject.HasFlag(FLAG_ONE))
	expect.False(subject.HasFlag(FLAG_TWO))
	expect.False(subject.HasFlag(FLAG_THREE))
	expect.True(subject.HasFlag(FLAG_FOUR))
}

func TestParseNoFlags(t *testing.T) {
	expect := assert.New(t)

	subject := NewBitFlagParser()

	expect.False(subject.HasFlag(FLAG_ONE))
	expect.False(subject.HasFlag(FLAG_TWO))
	expect.False(subject.HasFlag(FLAG_THREE))
	expect.False(subject.HasFlag(FLAG_FOUR))
}


func TestSetFlag(t *testing.T) {
	expect := assert.New(t)

	subject := NewBitFlagParser(FLAG_ONE, FLAG_FOUR)

	expect.True(subject.HasFlag(FLAG_ONE))
	expect.False(subject.HasFlag(FLAG_TWO))
	expect.False(subject.HasFlag(FLAG_THREE))
	expect.True(subject.HasFlag(FLAG_FOUR))

	subject.SetFlag(FLAG_TWO)
	expect.True(subject.HasFlag(FLAG_ONE))
	expect.True(subject.HasFlag(FLAG_TWO))
	expect.False(subject.HasFlag(FLAG_THREE))
	expect.True(subject.HasFlag(FLAG_FOUR))

	subject.SetFlag(FLAG_THREE)
	expect.True(subject.HasFlag(FLAG_ONE))
	expect.True(subject.HasFlag(FLAG_TWO))
	expect.True(subject.HasFlag(FLAG_THREE))
	expect.True(subject.HasFlag(FLAG_FOUR))
}
