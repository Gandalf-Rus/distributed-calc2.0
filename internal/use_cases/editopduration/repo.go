package editopduration

type repo interface {
	GetTokens() ([]string, error)
}

type config interface {
	ChangeOpDuration(operator string, duration int) error
}
