package apputils_test

import (
	"errors"
	"testing"
	"net"
	"github.com/stretchr/testify/assert"
	"github.com/sandreas/graft/apputils"
)

func mockErrorCallback(protocol, address string) (net.Conn, error) {
	return nil, errors.New("Mock error: " + protocol + ", " + address)
}

func TestGetOutboundIpAsStringFallback(t *testing.T) {
	expect := assert.New(t)

	fallback := "local"
	got, err := apputils.GetOutboundIpAsString(fallback, mockErrorCallback)
	expect.Equal(fallback, got)
	expect.Error(err)
}

