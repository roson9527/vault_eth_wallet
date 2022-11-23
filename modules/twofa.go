package modules

type Authenticator struct {
	Secret string
	Expire int
	Code   uint32
}
