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
	BcryptCost            int           `mapstructure:"BCYPT_COST"`
	ARGON2IDMemory        uint32        `mapstructure:"ARGON2ID_MEMORY"`
	ARGON2IDIteration     uint32        `mapstructure:"ARGON2ID_ITERATION"`
	ARGON2IDParallelsism  uint8         `mapstructure:"ARGON2ID_PARALLELISM"`
	ARGON2IDSaltLength    uint32        `mapstructure:"ARGON2ID_SALT_LENGTH"`
	ARGON2IDKeyLength     uint32        `mapstructure:"ARGON2ID_KEY_LENGTH=32"`
}

func LoadConfig(path string) (config Config, err error) {
	// TODO: Read key
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
