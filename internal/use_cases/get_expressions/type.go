package getexpressions

type request struct {
	UserId int
}

type response struct {
	Expressions []jsonExpression `json:"expressions"`
}

type jsonExpression struct {
	Body    string   `json:"body"`
	Result  *float64 `json:"result"`
	Status  string   `json:"status"`
	Message string   `json:"message"`
}
