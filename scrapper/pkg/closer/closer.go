package closer

import (
	"fmt"
	"log/slog"
	"sync"
)

type CloseFunction func() error

type Closer struct {
	mut     sync.Mutex
	closers []CloseFunction
	logger  *slog.Logger
}

func NewCloser(logger *slog.Logger) *Closer {
	return &Closer{
		logger: logger,
	}
}

func (c *Closer) Add(function CloseFunction) {
	c.mut.Lock()
	defer c.mut.Unlock()
	c.closers = append(c.closers, function)
}

func (c *Closer) Close() error {
	c.mut.Lock()
	defer c.mut.Unlock()

	var closeErrs []error
	for _, closeFunc := range c.closers {
		if err := closeFunc(); err != nil {
			c.logger.Error("Failed to close resource", slog.String("error", err.Error()))
			closeErrs = append(closeErrs, err)
		}
	}
	if len(closeErrs) > 0 {
		return fmt.Errorf("encountered errors while closing resources: %v", closeErrs)
	}

	return nil
}
