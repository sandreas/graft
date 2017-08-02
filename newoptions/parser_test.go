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

	expect.True(subject.hasFlag(FLAG_ONE))
	expect.False(subject.hasFlag(FLAG_TWO))
	expect.False(subject.hasFlag(FLAG_THREE))
	expect.True(subject.hasFlag(FLAG_FOUR))
}


func TestParseMultipleBitFlagParams(t *testing.T) {
	expect := assert.New(t)

	subject := NewBitFlagParser(FLAG_ONE, FLAG_FOUR)

	expect.True(subject.hasFlag(FLAG_ONE))
	expect.False(subject.hasFlag(FLAG_TWO))
	expect.False(subject.hasFlag(FLAG_THREE))
	expect.True(subject.hasFlag(FLAG_FOUR))
}

func TestParseNoFlags(t *testing.T) {
	expect := assert.New(t)

	subject := NewBitFlagParser()

	expect.False(subject.hasFlag(FLAG_ONE))
	expect.False(subject.hasFlag(FLAG_TWO))
	expect.False(subject.hasFlag(FLAG_THREE))
	expect.False(subject.hasFlag(FLAG_FOUR))
}
