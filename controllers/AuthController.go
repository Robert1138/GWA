package controllers

import (
	"encoding/json"
	"fmt"
	"goapp1/models"
	u "goapp1/util"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/csrf"
	"golang.org/x/crypto/bcrypt"
)

// SigningKey is the key or decoding secure cookies
var SigningKey = []byte(os.Getenv("token_secret"))

type userLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// JwtToken is the claims for inside out jwt token
type JwtToken struct {
	Exp    int64
	UserID uint
	jwt.StandardClaims
}

// Login will take credentials, verify and send a jwt in the response, TODO use secure cookies
func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user userLogin
		// Try to decode the request body into the struct. If there is an error,
		// respond to the client with the error message and a 400 status code.
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// when user credentials are identified we can set token and send cookie with jwt
		// boolean, User returned
		isValid, userCred := models.ValidUser(user.Username)
		if isValid {
			fmt.Println(userCred.UserID)
			http.SetCookie(w, PrepareCookie(CreateToken(w, userCred), "jwt"))
		} else {
			w.WriteHeader(http.StatusUnprocessableEntity) // send a 422 -- user they sent does not exist
			w.Write([]byte("422 - Unprocessable Entity - check what you're sending - user might not exist"))
		}
	}
}

// Register will
func Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser models.User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		newUser.Password, err = hashPassword(newUser.Password)
		if err != nil {
			fmt.Println("password couldnt be encrypted")
		}
		userAdded, err := models.AddUser(newUser)

		if err != nil {
			fmt.Println(err)
			if err == u.ErrUserExists {
				fmt.Println("user exists")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("409 - user with this username already exists"))
			} else if err == u.ErrPasswordInvalidFormat {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("401 - password too short"))
			}
		}
		http.SetCookie(w, PrepareCookie(CreateToken(w, newUser), "jwt"))
		fmt.Println(userAdded)
	}
}

// CsrfPostTest tests post with csrf protection
func CsrfPostTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("POST WENT THROUGH WITH VALID CSRF")

	}
}

// GetMessage returns a simple response -- placeholder to get csrf token
func GetMessage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u.Response(w, u.Message("success", "it worked"))
	}
}

// Hello lets you say hello to Obi-Wan Kenobi
func Hello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode("Hello there")
	}
}

func hashPassword(password string) (string, error) {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(passBytes), err
}

func checkPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}
	return true
}

// jWTFromCookie decodes the cookie and returns the jwt
func jWTFromCookie(r *http.Request) (string, error) {

	reqCookie, err := r.Cookie("jwt")
	if err != nil {
		return "", err
	}

	val := make(map[string]string)
	err = GetSecureCookie().Decode("jwt", reqCookie.Value, &val)
	if err == nil {
		fmt.Printf("JWT from decoded cookie %v", val["jwt"])
		return val["jwt"], nil
	}

	fmt.Println(err.Error())
	return "", err

}

// JwtMiddleware checks requests that require auth
func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		noAuthRequired := []string{"/hello", "/Login", "/Register"}
		requestURL := r.URL.Path
		requestMethod := r.Method
		fmt.Println(requestMethod)

		for i := range noAuthRequired {
			if noAuthRequired[i] == requestURL {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenStr, err1 := jWTFromCookie(r)
		if err1 != nil {
			fmt.Println(err1)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - cookie does not exists or error decoding value"))
			return
		}

		fmt.Println("jwt from cookie in middleware")
		fmt.Println(tokenStr)

		if tokenStr == "" {
			fmt.Println("no jwt")
			response := u.Message("fail", "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}

		// Check the token, do some error handling and then process the token
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

		if ok && token.Valid {
			// get claims here
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

// notes json tags in structs ---- field/member needs to be exported for encoding/json library to access it
//json.Unmarshal(rBody, &userLoginEvent)
//var userInfo []string
//u.Response(w, map[string]interface{}{"jwt": CreateToken(arr)})
//fmt.Println(w.Header)
// notes -- generating a random key has been problematic -- SOLVED problem was using old cookie in tests of the api route
// Login flow ----- Login with /Login route ------ will get a jwt token
// When GET request is made client will receive csrf_cookie and CSRF token that it must store
