package auth

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"strings"
	"time"
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
	mapClaims["exp"] = time.Now().Add(time.Second * time.Duration(a.TTL)).Unix()

	return mapClaims
}

func (a *Auth) CreateJWTToken(claims map[string]string) (string, error) {
	mc := a.defaultMapClaims()

	for k, v := range claims {
		mc[k] = v
	}

	t := jwt.NewWithClaims(jwt.SigningMethodRS256, mc)

	return t.SignedString(a.PrivKey)
}

// Authorize это преобразаватель который проверяет токен и валидирует его.
func (a *Auth) Authorize() (interface{}, error) {
	var newCtx context.Context
	var err error

	md, ok := metadata.FromIncomingContext(ctx)

	bearerToken, ok := md["authorization"]
	if len(bearerToken) == 0 || !ok {
		newCtx = context.WithValue(ctx, "authorized", false)
		return handler(newCtx, req)
	}

	mapClaims := jwt.MapClaims{}

	t := strings.TrimPrefix(bearerToken[0], "Bearer ")

	token, err := jwt.ParseWithClaims(t, &mapClaims, func(token *jwt.Token) (interface{}, error) {
		return a.PubKey, nil
	})

	if err != nil {
		return nil, grpc.Errorf(codes.Unauthenticated, "failed to parse token %v", err)
	}

	if !token.Valid {
		return nil, grpc.Errorf(codes.Unauthenticated, "invalid token")
	}

	newCtx = context.WithValue(ctx, "claims", mapClaims)

	return handler(newCtx, req)
}
