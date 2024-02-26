package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExpiresIn int    `mapstructure:"JWT_EXPIRESIN"`
	TokenAuth    *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)
	return cfg, nil
}
