package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyNameOk(t *testing.T) {
	expected := "request_id"

	ctx := context.Background()
	ctx = context.WithValue(ctx, ctxKeyName, expected)

	actual := GetKeyName(ctx)
	assert.Equal(t, expected, actual)
}

func TestGetKeyNameIDEmpty(t *testing.T) {
	ctx := context.Background()

	actual := GetKeyName(ctx)
	assert.Equal(t, "", actual)
}
