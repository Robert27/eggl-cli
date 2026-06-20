package kill

import "context"

type Killer interface {
	Kill(ctx context.Context, pid int, force bool) error
}

func DefaultKiller() Killer {
	return platformKiller()
}
