package domain

var (
	StatusWrapped = ThingStatus("Wrapped")
)

type ThingCode string
type ThingStatus string

func (t ThingCode) String() string {
	return string(t)
}

type Thing struct {
	Code   ThingCode
	Name   string
	Status ThingStatus
}
