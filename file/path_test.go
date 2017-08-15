package file

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestNewPath(t *testing.T) {
	expect := assert.New(t)

	p := NewPath("../data/fixtures")
	expect.Equal("../data/fixtures", p.String())

	p = NewPath("/tmp/test")
	expect.Equal("/tmp/test", p.String())

	p = NewPath("./data/fixtures")
	expect.Equal("data/fixtures", p.String())

	p = NewPath("./data/fixtures/")
	expect.Equal("data/fixtures", p.String())
}
