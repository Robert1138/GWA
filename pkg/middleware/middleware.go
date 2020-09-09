package middleware

import (
	"fmt"
	auth "goapp1/pkg/auth"
	"goapp1/util"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/csrf"
)

// var SigningKey = []byte(os.Getenv("token_secret"))
/*
type JwtToken struct {
	Exp    int64
	UserID uint
	jwt.StandardClaims
}
*/

// JwtMiddleware checks requests that require auth
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuthRequired := []string{"/hello", "/login", "/register", "/pic", "/time", "/item"}
		requestURL := r.URL.Path
		requestMethod := r.Method
		fmt.Println(requestMethod)

		for i := range noAuthRequired {
			if noAuthRequired[i] == requestURL {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenStr, err1 := auth.JWTFromCookie(r)
		// missing cookie or issue decoding it
		if err1 != nil {
			fmt.Println(err1)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - cookie does not exists or error decoding value"))
			return
		}
		// the cookie did not contain a jwt
		if tokenStr == "" {
			fmt.Println("no jwt")
			response := util.Message("fail", "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			util.Response(w, response)
			return
		}

		// Check the token, do some error handling and then process the token
		claimsTk := &auth.JwtToken{}
		token, err := jwt.ParseWithClaims(tokenStr, claimsTk, func(token *jwt.Token) (interface{}, error) {
			return auth.SigningKey, nil
		})
		// indicated token wasnt created correctly or it was modified
		if err != nil {
			fmt.Println("malformed token")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - malformed token"))
			return
		}

		claims, ok := token.Claims.(*auth.JwtToken)

		if ok && token.Valid {
			// get claims here -- in this case its the userId and do something like pass it on with the request
			fmt.Println("user id from valid jwt")
			fmt.Println(claims.UserID)

		} else {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 -Invalid token"))
			return
		}
		// fmt.Printf("%+v %+v", claims.UserID, claims.Exp)
		next.ServeHTTP(w, r)
	})
}

// CsrfTokenMiddleware adds the X-CSRF-Token to get requests
func CsrfTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		next.ServeHTTP(w, r)
	})
}
func RouteLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
