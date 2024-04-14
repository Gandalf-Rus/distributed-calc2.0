package registrateuser

type request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
