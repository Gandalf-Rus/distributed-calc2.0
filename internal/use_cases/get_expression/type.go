package getexpression

type request struct {
	UserId int
	ExitId string `json:"exit_id"`
}

type response struct {
	Body    string   `json:"body"`
	Result  *float64 `json:"result"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
}
