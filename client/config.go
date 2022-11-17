package client

type Config struct {
	Address    string `json:"address"`
	Token      string `json:"token"`
	SecretPath string `json:"secret_path"`
}
