package maya

// public interface OptionalSPI {}

type Selectable interface {
	aggr()
}
type Column string

func (c Column) aggr() {}

type Aggregate struct{}

func (c Aggregate) aggr() {}

// func (s *Selector) Select(cols ...Selectable) *Selector {
// 	panic("implement me")
// }
