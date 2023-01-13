package config

import (
	"os"
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
	ARGON2IDKeyLength     uint32        `mapstructure:"ARGON2ID_KEY_LENGTH"`
	Origin                string        `mapstructure:"CLIENT_ORIGIN"`
	EmailFrom             string        `mapstructure:"EMAIL_FROM"`
	SMTPHost              string        `mapstructure:"SMTP_HOST"`
	SMTPPass              string        `mapstructure:"SMTP_PASS"`
	SMTPPort              int           `mapstructure:"SMTP_PORT"`
	SMTPUser              string        `mapstructure:"SMTP_USER"`
	GrpcServerAddress     string        `mapstructure:"GRPC_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("config")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	loadKeyBuf(&config)
	return
}

func loadKeyBuf(cfg *Config) error {
	priBuf, err := os.ReadFile(cfg.AccessTokenPrivateKey)
	if err != nil {
		return err
	}

	pubBuf, err := os.ReadFile(cfg.AccessTokenPublicKey)
	if err != nil {
		return err
	}

	cfg.PrivBuf = priBuf
	cfg.PubBuf = pubBuf

	return nil
}
