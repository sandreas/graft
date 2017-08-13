package sftpd

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFileTreeIntegrate(t *testing.T) {
	expect := assert.New(t)
	fileTree := NewFileTree("data/fixtures")
	expect.NotNil(fileTree)

}