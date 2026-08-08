package sweep

import (
	"context"
	"errors"

	"github.com/openai/openai-go/v3"
	lop "github.com/samber/lo/parallel"
)

type Sweepable interface {
	Delete(ctx context.Context) error
}

func Sweep(ctx context.Context, sweepables []Sweepable) error {
	return errors.Join(lop.Map(sweepables, func(sweepable Sweepable, _ int) error {
		return sweepable.Delete(ctx)
	})...)
}

type SweeperFn func(ctx context.Context, client *openai.Client) ([]Sweepable, error)
