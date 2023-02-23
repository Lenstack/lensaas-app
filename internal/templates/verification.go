package templates

type VerificationCode struct {
	Name  string
	Email string
	Code  string
}

type VerificationEmail struct {
	Name  string
	Email string
	Token string
}
