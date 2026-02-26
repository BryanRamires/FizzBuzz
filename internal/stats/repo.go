package stats

import "context"

type Repository interface {
	Inc(context.Context, Key) error
	Top(context.Context) (Top, bool, error)
}
