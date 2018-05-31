package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"golang.org/x/crypto/bcrypt"
)

const (
	fakeuser    = "chrisd"
	fakepass    = "hotdog23"
	ctxNameKey  = "name"
	ctxRolesKey = "roles"
)

type User struct {
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	Roles     []string `json:"-"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
}

var users = map[string]User{
	fakeuser: User{
		Username:  fakeuser,
		Firstname: "Chris",
		Lastname:  "Dyer",
		Roles: []string{
			"superuser",
		},
	},
}

var mySigningKey = []byte("secret")
var fakeuserHash []byte

func init() {
	// Simulating Hash stored in user table for credentials.
	fakeuserHash, _ = bcrypt.GenerateFromPassword([]byte(fakepass), bcrypt.DefaultCost)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/token", PostForToken)

	authenticatedGroup := r.Group(nil)
	authenticatedGroup.Use(Authenticate())
	authenticatedGroup.Get("/greeting", func(w http.ResponseWriter, r *http.Request) {
		rolesContext := r.Context().Value(ctxRolesKey)
		roleString, _ := rolesContext.(string)
		roles := strings.FieldsFunc(roleString, func(c rune) bool {
			return ',' == c
		})

		var isAdmin bool
		for _, role := range roles {
			fmt.Println(role)
			if role == "admin" {
				isAdmin = true
			}
		}

		if isAdmin {
			w.Write([]byte(fmt.Sprintf("Hello %v!!",
				r.Context().Value(ctxNameKey))))
			return
		}

		http.Error(w, "", http.StatusForbidden)
	})

	http.ListenAndServe(":8080", r)
}

func PostForToken(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	if populatedUser, ok := users[user.Username]; ok {
		if err := bcrypt.CompareHashAndPassword(fakeuserHash, []byte(user.Password)); err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)

		// Simplifying implementation by using static data.
		claims["name"] = fmt.Sprintf("%v %v", populatedUser.Firstname,
			populatedUser.Lastname)
		claims["roles"] = strings.Join(populatedUser.Roles, ",")
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		tokenString, _ := token.SignedString(mySigningKey)

		w.Write([]byte(tokenString))
		return
	}

	http.Error(w, "", http.StatusUnauthorized)
}

func Authenticate() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authorization := r.Header.Get("Authorization")
			trimmedAuth := strings.Fields(authorization)

			// Trim out Bearer from Authorization Header
			if authorization == "" || len(trimmedAuth) == 0 {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(trimmedAuth[1], claims,
				func(token *jwt.Token) (interface{}, error) {
					return mySigningKey, nil
				})
			if err != nil {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxNameKey, claims["name"])
			ctx = context.WithValue(ctx, ctxRolesKey, claims["roles"])
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}