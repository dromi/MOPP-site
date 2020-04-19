package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtRefreshTime = 30 * time.Second
	jwtLifetime    = 5 * time.Minute
	jwtSecret      = "my_secret_key"
	tokenName      = "mopp_token"
)

var errTokenExpired = errors.New("Token expired")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func HashPassword(password *string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), 8)
	if err != nil {
		panic(err)
	}
	*password = string(hash)
}

func ArePasswordsMatching(passwordHashed, passwordPlain *string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(*passwordHashed), []byte(*passwordPlain))
	return err == nil
}

func CreateJWTCookie(name string) (*http.Cookie, error) {
	expirationTime := time.Now().Add(jwtLifetime)
	claims := &Claims{
		Username: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:    tokenName,
		Value:   tokenString,
		Expires: expirationTime,
	}, nil
}

func updateJWTCookie(claims *Claims) (*http.Cookie, error) {
	expirationTime := time.Now().Add(jwtLifetime)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}
	return &http.Cookie{
		Name:    tokenName,
		Value:   tokenString,
		Expires: expirationTime,
	}, nil
}

func getJWTCookie(r *http.Request) (*Claims, error) {
	c, err := r.Cookie(tokenName)
	if err != nil {
		return nil, err
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, errTokenExpired
	}
	return claims, nil
}

func authHandler(fn func(http.ResponseWriter, *http.Request, *MetaData)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims, err := getJWTCookie(r)
		if err != nil {
			if err == http.ErrNoCookie {
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}
			if err == jwt.ErrSignatureInvalid {
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}
			if err == errTokenExpired {
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < jwtRefreshTime {
			cookie, err := updateJWTCookie(claims)
			if err != nil {
				panic(err)
			}
			http.SetCookie(w, cookie)
		}
		fn(w, r, &MetaData{Username: claims.Username})
	}
}
