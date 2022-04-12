package config

type JwtConfig struct {
	PrivateKey string `mapstructure:"private_key"`
	ExpiresAt  int    `mapstructure:"expires_at"`
	Issuer     string `mapstructure:"issuer"`
	Subject    string `mapstructure:"subject"`
	Audience   string `mapstructure:"audience"`
}
