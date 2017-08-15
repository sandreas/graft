package file

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewPath(t *testing.T) {
	expect := assert.New(t)

	p := NewPath("../data/test.txt")
	expect.Equal(".."+separator+"data"+separator+"test.txt", p.String())
	expect.True(p.IsFile())
	expect.False(p.IsDir())
	expect.False(p.IsAbs())


	p = NewPath("/tmp/test")
	expect.Equal(""+separator+"tmp"+separator+"test", p.String())
	expect.True(p.IsFile())
	expect.False(p.IsDir())
	expect.True(p.IsAbs())


	p = NewPath("./data/test.txt")
	expect.Equal("data"+separator+"test.txt", p.String())
	expect.True(p.IsFile())
	expect.False(p.IsDir())
	expect.False(p.IsAbs())

	p = NewPath("./data/fixtures/")
	expect.Equal("data"+separator+"fixtures"+separator, p.String())
	expect.False(p.IsFile())
	expect.True(p.IsDir())
	expect.False(p.IsAbs())

	p = NewPath(".\\data//fixtures\\")
	expect.Equal("data"+separator+"fixtures"+separator+"", p.String())
	expect.False(p.IsFile())
	expect.True(p.IsDir())
	expect.False(p.IsAbs())

	//p := NewPath("\\\\tmp\\test\\uncpath.txt")
	//expect.Equal(separator+separator+"tmp"+separator+"test"+separator+"uncpath.txt", p.String())
	//expect.False(p.IsFile())
	//expect.True(p.IsDir())
	//expect.True(p.IsAbs())
}
