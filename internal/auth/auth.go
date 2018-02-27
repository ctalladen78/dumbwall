package auth

import (
	"crypto/rsa"
	"io/ioutil"
	"path"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	defaultTTL = 3600
)

// Auth имеет метод GRPC преобзавателя
type Auth struct {
	ttl int64

	domain string

	pubKey  *rsa.PublicKey
	privKey *rsa.PrivateKey
}

func New(certsDir, domain string, ttl int64) (*Auth, error) {
	key, err := ioutil.ReadFile(path.Join(certsDir, "private.key"))
	if err != nil {
		return nil, err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	key, err = ioutil.ReadFile(path.Join(certsDir, "public.key"))
	if err != nil {
		return nil, err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return nil, err
	}

	return &Auth{
		ttl:     ttl,
		domain:  domain,
		privKey: privKey,
		pubKey:  pubKey,
	}, nil
}

func (a *Auth) defaultMapClaims() jwt.MapClaims {
	mapClaims := jwt.MapClaims{}

	mapClaims["iat"] = time.Now().Unix()
	mapClaims["exp"] = time.Now().Add(time.Second * time.Duration(a.ttl)).Unix()

	return mapClaims
}

func (a *Auth) CreateJWTToken(claims map[string]string) (string, error) {
	mc := a.defaultMapClaims()

	for k, v := range claims {
		mc[k] = v
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, mc)

	return t.SignedString(a.privKey)
}

func (a *Auth) Validate(token string) (map[string]interface{}, error) {
	mapClaims := jwt.MapClaims{}

	t, err := jwt.ParseWithClaims(token, &mapClaims, func(token *jwt.Token) (interface{}, error) {
		return a.pubKey, nil
	})

	if err != nil {
		return mapClaims, err
	}

	if !t.Valid {
		return mapClaims, err
	}

	return mapClaims, nil
}
