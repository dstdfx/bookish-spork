package testutils

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// RespToJSON converts given type to JSON string.
func RespToJSON(t *testing.T, v interface{}) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	err := enc.Encode(v)
	assert.NoError(t, err)

	return buf.String()
}
