package postexpression

type request struct {
	userId     int
	ExitId     string `json:"exit_id"`
	Expression string `json:"expression"`
}
