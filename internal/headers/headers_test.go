package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nAHeader:     good   \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	// assert.Equal(t, "good", headers["aheader"])
	assert.Equal(t, "good", headers.Get("aheader"))
	assert.Equal(t, 47, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: Invalid character
	headers = NewHeaders()
	data = []byte("H@ost: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)

	//Test: Multiple values for same field-name
	headers = NewHeaders()
	data = []byte("machine: switch\r\nmachine: 3ds\r\n\r\n")
	n, done, err = headers.Parse(data)
	assert.Equal(t, "switch, 3ds", headers.Get("machine"))
	require.NoError(t, err)
	assert.Equal(t, 33, n)
	assert.True(t, done)
}
