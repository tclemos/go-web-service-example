package domain

type ThingCode string

func (t ThingCode) String() string {
	return string(t)
}

type Thing struct {
	Code ThingCode `json:"code"`
	Name string    `json:"name"`
}
