package qqcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEntityExpired(t *testing.T) {
	e := entity{
		value:        nil,
		expiredAfter: time.Now().UTC().Add(-1 * time.Second).UnixNano(),
	}
	require.True(t, e.isExpired())
}

func TestEntityPermanent(t *testing.T) {
	e := entity{
		value:        nil,
		expiredAfter: 0,
	}
	require.False(t, e.isExpired())
}
