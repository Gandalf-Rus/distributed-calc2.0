package loginuser

type request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type response struct {
	JwtToken string `json:"jwt_token"`
}
