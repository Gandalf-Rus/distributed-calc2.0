package postexpression

type request struct {
	userId     int64
	ExitId     string `json:"exit_id"`
	Expression string `json:"expression"`
}
