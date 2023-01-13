package config

import (
	"time"

	"github.com/spf13/viper"
)

// Vipper package using mapstructure
type Config struct {
	DBUri                 string        `mapstructure:"MONGODB_LOCAL_URI"`
	RedisUri              string        `mapstructure:"REDIS_URL"`
	Port                  string        `mapstructure:"PORT"`
	AccessTokenPrivateKey string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey  string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	PrivBuf               []byte        `mapstructure:"-"`
	PubBuf                []byte        `mapstructure:"-"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge     int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge    int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
