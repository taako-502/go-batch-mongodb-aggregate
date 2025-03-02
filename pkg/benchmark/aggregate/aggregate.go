package aggregate

import "github.com/taako-502/go-batch-mongodb-aggregate/pkg/infrastructure"

type Aggregate struct {
	Infrastructure *infrastructure.Infrastructure
}

func NewAggregate(infrastructure *infrastructure.Infrastructure) *Aggregate {
	return &Aggregate{
		Infrastructure: infrastructure,
	}
}
