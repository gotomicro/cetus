package maya

type Entity interface {
	Id() int64
	SetId(id int64)
}
type BaseEntity struct {
	id int64
}

func (b *BaseEntity) Id() int64 {
	return b.id
}
func (b *BaseEntity) SetId(id int64) {
	b.id = id
}

type myEntity struct {
	BaseEntity
}

func Insert[E Entity](e *E) {
}
