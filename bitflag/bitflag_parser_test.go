package bitflag_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/bitflag"
)
const (
	FLAG_ONE   bitflag.Flag = 1 << iota
	FLAG_TWO
	FLAG_THREE
	FLAG_FOUR
)


func TestParseSingleOrConcatBitFlagParam(t *testing.T) {
	expect := assert.New(t)

	subject := bitflag.NewParser(FLAG_ONE|FLAG_FOUR)

	expect.True(subject.HasFlag(FLAG_ONE))
	expect.False(subject.HasFlag(FLAG_TWO))
	expect.False(subject.HasFlag(FLAG_THREE))
	expect.True(subject.HasFlag(FLAG_FOUR))
}


func TestParseMultipleBitFlagParams(t *testing.T) {
	expect := assert.New(t)

	subject := bitflag.NewParser(FLAG_ONE, FLAG_FOUR)

	expect.True(subject.HasFlag(FLAG_ONE))
	expect.False(subject.HasFlag(FLAG_TWO))
	expect.False(subject.HasFlag(FLAG_THREE))
	expect.True(subject.HasFlag(FLAG_FOUR))
}

func TestParseNoFlags(t *testing.T) {
	expect := assert.New(t)

	subject := bitflag.NewParser()

	expect.False(subject.HasFlag(FLAG_ONE))
	expect.False(subject.HasFlag(FLAG_TWO))
	expect.False(subject.HasFlag(FLAG_THREE))
	expect.False(subject.HasFlag(FLAG_FOUR))
}


func TestSetFlag(t *testing.T) {
	expect := assert.New(t)

	subject := bitflag.NewParser(FLAG_ONE, FLAG_FOUR)

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
