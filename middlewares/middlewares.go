package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

// var jwtKey = os.Getenv("JWT_secret") TODO: Uncomment and add secret to .env file
var jwtKey = []byte("my_secret")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var Authorization = Use(isLogged, isAuthorized, refreshToken)
var Logged = Use(isLogged, refreshToken)

/* Middleware definition */
type Middleware func(http.Handler) http.Handler

func (mw Middleware) ThenFunc(finalPage func(http.ResponseWriter, *http.Request)) http.Handler {
	return mw(http.HandlerFunc(finalPage))
}

func Use(mw ...Middleware) Middleware {
	return func(finalPage http.Handler) http.Handler {
		for i := len(mw) - 1; i >= 0; i-- {
			finalPage = mw[i](finalPage)
		}
		return finalPage
	}
}

func isLogged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("Token")
		if err != nil {
			context.Set(r, "userLogged", false)
		} else {
			claims, token, err := parseCookie(c)
			if err != nil || !token.Valid {
				context.Set(r, "userLogged", false)
			} else {
				context.Set(r, "userLogged", true)
				context.Set(r, "claims", claims)
			}
		}
		next.ServeHTTP(w, r)
	})
}

func isAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userLogged, ok := context.Get(r, "userLogged").(bool)
		if userLogged == false || !ok {
			fmt.Println("Access denied")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func refreshToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := context.Get(r, "claims").(*Claims)
		if ok {
			if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 5*time.Minute {
				SetJwtCookie(w, claims.Username)
				fmt.Println("Refreshed Token")
			} else {
				fmt.Println("Not Refreshed")
			}
		}
		next.ServeHTTP(w, r)
	})
}

/* Utils */
func parseCookie(c *http.Cookie) (*Claims, *jwt.Token, error) {
	tokenString := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, nil, err
	}
	return claims, token, nil
}

func SetJwtCookie(w http.ResponseWriter, username string) {
	validToken, expirationTime, err := GenerateJWT(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "Token",
		Value:   validToken,
		Expires: expirationTime,
	})
}

func GenerateJWT(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Println("Something went wrong: %s", err)
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}
