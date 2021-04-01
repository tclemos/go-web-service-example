package domain

import "github.com/gofrs/uuid"

var (
	StatusWrapped = ThingStatus("Wrapped")
)

type ThingID int
type ThingCode uuid.UUID
type ThingStatus string

func (t ThingCode) String() string {
	return uuid.UUID(t).String()
}

type Thing struct {
	ID     ThingID
	Code   ThingCode
	Name   string
	Status ThingStatus
}

type ThingCreated struct {
	Thing Thing
}
