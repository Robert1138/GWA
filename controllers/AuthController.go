package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"goapp1/models"
	u "goapp1/util"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SigningKey = []byte(os.Getenv("token_secret"))

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims for inside out jwt token
type JwtToken struct {
	Exp    int64
	UserID uint
	jwt.StandardClaims
}

// GetMessage returns a simple response
func GetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u.Response(w, u.Message("success", "it worked"))
	}
}

// Login will take credentials, verify and send a jwt in the response, TODO use secure cookies
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO check credentials THEN return response
		// Will hold user info from login post
		var user userLogin
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Do something with the userLogin struct
		fmt.Println("User: " + user.Username + " " + user.Password)

		isValid, userCred := models.ValidUser(user.Username)
		// when user credentials are identified we can set token and send cookie with jwt
		if isValid {
			fmt.Println(userCred.UserID)
			w.Header().Set("X-CSRF-Token", "bleh")
			http.SetCookie(w, PrepareCookie(CreateToken(w, userCred), "jwt"))
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity) // send a 422 -- user they sent does not exist
			w.Write([]byte("422 - Unprocessable Entity - check what you're sending - user might not exist"))
		}
	}
}

// jWTFromCookie decodes the cookie and returns the jwt
func Test() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqCookie, err := r.Cookie("jwt")
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 cookie does not exist"))
		}

		if reqCookie == nil {
			fmt.Println("reqCookie is nil")
		} else {
			val := make(map[string]string)
			err = GetSecureCookie().Decode("jwt", reqCookie.Value, &val)
			if err == nil {
				fmt.Printf("JWT from decoded cookie %v", val["jwt"])
			} else {
				fmt.Println(err.Error())
			}
		}
		u.Response(w, u.Message("success", "Test()"))
	}
}

// jWTFromCookie decodes the cookie and returns the jwt
func jWTFromCookie(r *http.Request) (string, error) {

	reqCookie, err := r.Cookie("jwt")
	if err != nil {
		return "", errors.New("cookie does not exist")
	}

	val := make(map[string]string)
	err = GetSecureCookie().Decode("jwt", reqCookie.Value, &val)
	if err == nil {
		fmt.Printf("JWT from decoded cookie %v", val["jwt"])
		return val["jwt"], nil
	} else {
		fmt.Println(err.Error())
		return "", err
	}

}

// JwtMiddleware checks requests that require auth
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuthRequired := []string{"/api/Login", "/api/hello"}
		requestURL := r.URL.Path
		requestMethod := r.Method
		fmt.Println(requestMethod)

		for i := range noAuthRequired {
			if noAuthRequired[i] == requestURL {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")
		fmt.Println("JWT from client request", tokenHeader)
		if tokenHeader == "" {
			fmt.Println("no jwt")
			response := u.Message("fail", "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}

		newToken, err1 := jWTFromCookie(r)
		if err1 != nil {
			fmt.Println(err1)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - cookie does not exists or error decoding value"))
		} else {
			fmt.Println("jwt from cookie in middleware")
			fmt.Println(newToken)
		}
		// Check the token, do some error handling and then process the token
		tokenStr := tokenHeader[7:len(tokenHeader)]
		claimsTk := &JwtToken{}
		token, err := jwt.ParseWithClaims(tokenStr, claimsTk, func(token *jwt.Token) (interface{}, error) {
			return SigningKey, nil
		})

		if err != nil {
			fmt.Println("malformed token")
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 - malformed token"))
			return
		}

		claims, ok := token.Claims.(*JwtToken)
		fmt.Println("user id from jwt")
		fmt.Println(claims.UserID)

		if ok && token.Valid {
			// get claims here
			fmt.Println("user id from jwt")
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

// CreateToken call this to create the token using secret and return the string
// User auth data includes role, username,
func CreateToken(w http.ResponseWriter, userInfo models.User) string {
	expTime := time.Now().Add(time.Hour * 24).Unix()
	claims := &JwtToken{Exp: expTime, UserID: userInfo.UserID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SigningKey)
	// need to figure out a better way to handle when this function ails to create a token
	// at the moment it returns a 500
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("token creation failed, try again"))
	}

	fmt.Println(tokenString)
	return tokenString
}

/* old create token ---- works

func CreateToken(userInfo models.User) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["UserID"] = userInfo.UserID
	claims["Exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(SigningKey)
	//w.Write([]byte(tokenString))
	fmt.Println(tokenString)
	return tokenString
}

*/
// CreateCookie will put the jwt inside a secure cookie and set it
func CreateCookie() string {
	return ""
}

// notes json tags in structs ---- field/member needs to be exported for encoding/json library to access it
//json.Unmarshal(rBody, &userLoginEvent)
//var userInfo []string
//u.Response(w, map[string]interface{}{"jwt": CreateToken(arr)})
//fmt.Println(w.Header)
// notes -- generating a random key has been problematic
