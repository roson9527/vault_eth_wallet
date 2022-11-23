package modules

// SocialID is a struct that contains the information of a discord / twitter / telegram user
type SocialID struct {
	User        string   `json:"user" mapstructure:"user"`                                       // PrivateKey is the private key of the discord
	Password    string   `json:"password,omitempty" mapstructure:"password,omitempty"`           // PublicKey is the public key of the wallet
	Mobile      string   `json:"mobile,omitempty" mapstructure:"mobile,omitempty"`               // Address is the address of the wallet
	TwoFASecret string   `json:"two_fa_secret,omitempty" mapstructure:"two_fa_secret,omitempty"` // Address is the address of the wallet
	UpdateTime  int64    `json:"update_time" mapstructure:"update_time"`                         // key pair update time
	NameSpaces  []string `json:"namespaces,omitempty" mapstructure:"namespaces,omitempty"`       // 用于项目区分
	App         string   `json:"app" mapstructure:"app,omitempty"`                               // 用于类型区分
}
