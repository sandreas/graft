package sftpd_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/sftpd"
)

func TestFileTreeIntegrate(t *testing.T) {
	expect := assert.New(t)
	fileTree := sftpd.NewFileTree("data/fixtures")
	expect.NotNil(fileTree)

}