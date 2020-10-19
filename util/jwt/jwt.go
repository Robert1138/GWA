package jwt

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	u "github.com/Robert1138/GWA/pkg/user"
	_ "github.com/Robert1138/GWA/util/log"
	"github.com/joho/godotenv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/securecookie"
)

// pointer to the struct that encodes and decodes
var sCookie *securecookie.SecureCookie

// SigningKey is for the jwt
var SigningKey = []byte(os.Getenv("token_secret"))

// CustomToken represents the contents of the token
type CustomToken struct {
	UserID uint
	jwt.StandardClaims
}

func init() {
	godotenv.Load("..\\src\\github.com\\Robert1138\\GWA\\.env") // must load env since init() runs before main
	var hashKey = []byte(os.Getenv("hashkey"))
	var blockKey = []byte(os.Getenv("blockkey"))
	sCookie = securecookie.New(hashKey, blockKey)
	fmt.Println("sCookie has been initialized")
	// fmt.Println(base64.StdEncoding.EncodeToString([]byte(securecookie.GenerateRandomKey(32))))
	// fmt.Println(base64.StdEncoding.EncodeToString([]byte(securecookie.GenerateRandomKey(32))))
}

// CreateToken creates a token with user auth data (UserID, role, etc) and returns it as a string
func CreateToken(w http.ResponseWriter, userInfo u.User) string {
	expTime := time.Now().Add(time.Minute * 10).Unix()
	claims := &CustomToken{UserID: userInfo.UserID, StandardClaims: jwt.StandardClaims{ExpiresAt: expTime}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SigningKey)

	if err != nil { // TODO if token fails return error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("token creation failed, try again"))
	}

	fmt.Println(tokenString)
	return tokenString
}

// GetSecureCookie will return global SecureCookie variable
func GetSecureCookie() *securecookie.SecureCookie {
	return sCookie
}

// PrepareCookie will prep a secure cookie with a payload
// params allow you to add a jwt plus its value as well as a csrf and its value
// Options are jwt and _gorilla_csrf
func PrepareCookie(token string, tokenName string) *http.Cookie {
	val := map[string]string{
		tokenName: token,
	}

	var cookie *http.Cookie
	if encodedCookie, err := sCookie.Encode(tokenName, val); err == nil {
		tempCookie := &http.Cookie{
			Name:     tokenName,
			Value:    encodedCookie,
			Path:     "/",
			HttpOnly: true,
		}
		cookie = tempCookie
	} else {
		fmt.Println(err) // TODO: log this - will require loggger to be passed in or access directly from log pkg
	}
	return cookie
}

// ParseFromCookie decodes the cookie and returns the jwt
func ParseFromCookie(r *http.Request) (string, error) {

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

// ParseUserID takes a JWT token string and returns the UserID claim
func ParseUserID(tokenStr string) (int, error) {
	claimsTK := &CustomToken{}
	token, err := jwt.ParseWithClaims(tokenStr, claimsTK, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})

	if err != nil {
		return 0, errors.New("malformed token")
	}

	claims, ok := token.Claims.(*CustomToken)
	if ok && token.Valid {
		return int(claims.UserID), nil
	}
	return 0, nil

}
