package modules

type Authenticator struct {
	Secret string `json:"-" mapstructure:"-"`
	Expire int    `json:"expire" mapstructure:"expire"`
	Code   uint32 `json:"code" mapstructure:"code"`
}
