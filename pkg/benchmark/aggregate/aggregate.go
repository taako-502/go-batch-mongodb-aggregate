package aggregate

import "github.com/taako-502/go-batch-mongodb-aggregate/pkg/infrastructure"

type Aggregate struct {
	infrastructure *infrastructure.Infrastructure
}

func NewAggregate(infrastructure *infrastructure.Infrastructure) *Aggregate {
	return &Aggregate{
		infrastructure: infrastructure,
	}
}
