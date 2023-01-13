package services

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"github.com/TranQuocToan1996/redislearn/config"
	"github.com/golang-jwt/jwt"
)

var (
	JwtObj *jwtProvider
)

type jwtProvider struct {
	config    config.Config
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

type UserClaimData struct {
	UID       string    `json:"uid"`
	LoginTime time.Time `json:"loginTime"`
}

type UserClaim struct {
	*jwt.StandardClaims
	User UserClaimData `json:"user"`
}

func NewJWT(cfg config.Config) error {
	JwtObj = &jwtProvider{config: cfg}

	blockPriv, _ := pem.Decode(cfg.PrivBuf)

	privKey, err := x509.ParsePKCS1PrivateKey(blockPriv.Bytes)
	if err != nil {
		return err
	}

	JwtObj.signKey = privKey

	// blockPub, _ := pem.Decode(cfg.PubBuf)

	// pubInterface, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
	// if err != nil {
	// 	return err
	// }

	// JwtObj.verifyKey = pubInterface.(*rsa.PublicKey)
	JwtObj.verifyKey = &privKey.PublicKey
	return nil
}

func (j *jwtProvider) CreateToken(uid string) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims = &UserClaim{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.config.AccessTokenExpiresIn * time.Minute).Unix(),
		},
		UserClaimData{UID: uid, LoginTime: time.Now()},
	}

	return t.SignedString(j.signKey)
}

func (j *jwtProvider) CreateRefreshToken(uid string) (string, error) {
	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims = &UserClaim{
		&jwt.StandardClaims{
			ExpiresAt: time.Now().Add(j.config.RefreshTokenExpiresIn * time.Minute).Unix(),
		},
		UserClaimData{UID: uid, LoginTime: time.Now()},
	}

	return t.SignedString(j.signKey)
}

func (j *jwtProvider) ValidateToken(token string) (interface{}, error) {
	tokenParse, err := jwt.ParseWithClaims(token, &UserClaim{}, func(t *jwt.Token) (interface{}, error) {
		return j.verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	if tokenParse == nil {
		return nil, errors.New("cant parse token")
	}

	return tokenParse.Claims.(*UserClaim), nil
}
