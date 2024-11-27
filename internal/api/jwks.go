package api

import (
	"errors"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"maple/internal/conf"
)

type MapleJWKS struct {
	KeyFunc *keyfunc.Keyfunc
}

func (j MapleJWKS) GetKeyFuncForParse() func(token *jwt.Token) (any, error) {
	return (*j.KeyFunc).Keyfunc
}

func NewMapleJWKS(config *conf.MapleConfigurations) (*MapleJWKS, error) {
	if config.Auth.KeysEndpoint == "" {
		return nil, errors.New("JWKs cache is not created since KEYS_ENDPOINT is not set")
	}

	keyFn, err := keyfunc.NewDefault([]string{config.Auth.KeysEndpoint})
	if err != nil {
		return nil, err
	}

	return &MapleJWKS{
		KeyFunc: &keyFn,
	}, nil
}
