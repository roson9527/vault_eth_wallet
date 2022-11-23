package modules

// SocialID is a struct that contains the information of a discord / twitter / telegram user
type SocialID struct {
	User        string         `json:"user" mapstructure:"user"`                                       // User is the user id of the discord / twitter / telegram user
	Password    string         `json:"password,omitempty" mapstructure:"password,omitempty"`           // Password is the password of the discord / twitter / telegram user
	Mobile      string         `json:"mobile,omitempty" mapstructure:"mobile,omitempty"`               // Mobile is the mobile of the discord / twitter / telegram user
	TwoFASecret string         `json:"two_fa_secret,omitempty" mapstructure:"two_fa_secret,omitempty"` // TwoFASecret is the two-factor authentication secret of the discord / twitter / telegram user
	GAuth       *Authenticator `json:"g_auth,omitempty" mapstructure:"g_auth,omitempty"`               // GAuth is the Google authenticator of the discord / twitter / telegram user
	UpdateTime  int64          `json:"update_time" mapstructure:"update_time"`                         // UpdateTime is the update time of the socialId
	NameSpaces  []string       `json:"namespaces,omitempty" mapstructure:"namespaces,omitempty"`       // NameSpaces is the namespaces of the socialId
	App         string         `json:"app" mapstructure:"app,omitempty"`                               // App is the app of the socialId
}
