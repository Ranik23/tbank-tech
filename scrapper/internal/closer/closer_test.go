package closer

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestCloser_Add(t *testing.T) {
	logger := slog.Default()

	c := NewCloser(logger)

	closeFunc := func() error {
		return nil
	}

	c.Add(closeFunc)

	assert.Len(t, c.closers, 1, "Closer should contain 1 function")
}


func TestCloser_Close_Success(t *testing.T) {
	
	logger := slog.Default()

	c := NewCloser(logger)

	closeFunc := func() error {
		return nil
	}

	c.Add(closeFunc)

	err := c.Close()

	assert.NoError(t, err, "Close should succeed without errors")
}

func TestCloser_Close_Failure(t *testing.T) {
	
	logger := slog.Default()

	c := NewCloser(logger)

	closeFunc := func() error {
		return errors.New("close failed")
	}

	c.Add(closeFunc)

	err := c.Close()

	assert.Error(t, err, "Close should return error when close function fails")
}