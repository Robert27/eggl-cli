package kill

import (
	"context"
	"fmt"
)

type Finder interface {
	FindListeners(ctx context.Context, port int) ([]Process, error)
}

var errUnsupportedPlatform = fmt.Errorf("finding listeners on this platform is not supported")

func DefaultFinder() Finder {
	return platformFinder()
}

type unsupportedFinder struct{}

func (unsupportedFinder) FindListeners(ctx context.Context, port int) ([]Process, error) {
	return nil, errUnsupportedPlatform
}
