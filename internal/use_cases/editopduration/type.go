package editopduration

type request struct {
	Operator string `json:"operator"`
	Delay    int    `json:"delay"`
}
