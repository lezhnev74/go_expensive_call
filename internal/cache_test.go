package internal

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	InitCache()

	counter := 0
	expensiveCall := func() (string, error) {
		counter++
		return fmt.Sprintf("%d", counter), nil
	}

	_, err := Cache("key", time.Second, expensiveCall)
	require.NoError(t, err)
	_, err = Cache("key", time.Second, expensiveCall)
	require.NoError(t, err)

	require.Equal(t, 1, counter)
}

func TestCacheEviction(t *testing.T) {
	InitCache()

	counter := 0
	expensiveCall := func() (string, error) {
		counter++
		return fmt.Sprintf("%d", counter), nil
	}

	_, err := Cache("key", time.Millisecond, expensiveCall)
	require.NoError(t, err)
	require.Equal(t, 1, counter)

	time.Sleep(time.Millisecond * 2) // this evicts the Cache

	_, err = Cache("key", time.Millisecond, expensiveCall)
	require.NoError(t, err)

	require.Equal(t, 2, counter)
}
