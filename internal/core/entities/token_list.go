package entities

type TokenList struct {
	Type       string
	Token      string
	UserId     string
	Blocked    bool
	Expiration int64
}
