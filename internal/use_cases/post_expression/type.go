package postexpression

type request struct {
	ExitId     string `json:"exit_id"`
	Expression string `json:"expression"`
}
