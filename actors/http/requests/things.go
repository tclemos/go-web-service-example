package requests

type CreateThing struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type GetThing struct {
	Code string `json:"code"`
}
