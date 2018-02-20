package auth

import (
	"context"
	"crypto/rsa"
	"io/ioutil"
	"strings"
	"time"
)

const (
	// defaultTTL это время токена поумолчанию - 1 час,
	defaultTTL = 3600
)

// Auth имеет метод GRPC преобзавателя
type Auth struct {
	AuthConfig

	// публичные и приватные ключи
	PubKey  *rsa.PublicKey
	PrivKey *rsa.PrivateKey
}

// AuthConfig конфигурации для Auth
type AuthConfig struct {
	PubKeyPath  string
	PrivKeyPath string

	// время для токена
	TTL int64

	// домен на который дается токен
	Domain string
}

func (a *Auth) LoadKeys() error {
	var (
		err error
		key []byte
	)

	// загружаем приватный ключ
	key, err = ioutil.ReadFile(a.PrivKeyPath)
	if err != nil {
		return err
	}

	a.PrivKey, err = jwt.ParseRSAPrivateKeyFromPEM(key)
	if err != nil {
		return err
	}

	// загружаем публичный ключ
	key, err = ioutil.ReadFile(a.PubKeyPath)
	if err != nil {
		return err
	}

	a.PubKey, err = jwt.ParseRSAPublicKeyFromPEM(key)
	if err != nil {
		return err
	}

	return nil
}

// NewAuth принимает конфигурации и создает обьект с методом преборазателя для GRPC
func NewAuth(c AuthConfig) (*Auth, error) {
	auth := &Auth{
		AuthConfig: c,
	}

	err := auth.LoadKeys()
	if err != nil {
		return nil, err
	}

	if auth.TTL == 0 {
		auth.TTL = defaultTTL
	}

	return auth, nil
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
func (a *Auth) Authorize(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
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
