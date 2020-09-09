package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// SigningKey is the key for token creation should be changed frequently and prob should go unexported
var SigningKey = []byte(os.Getenv("token_secret"))

// JwtToken is the struct for storing token information probably doesnt need to be exported either
type JwtToken struct {
	Exp    int64
	UserID uint
	jwt.StandardClaims
}

// CreateToken call this to create the token using secret and return the string
// User auth data includes role, username,
func CreateToken(w http.ResponseWriter, userInfo User) string {
	expTime := time.Now().Add(time.Hour * 24).Unix()
	claims := &JwtToken{Exp: expTime, UserID: userInfo.UserID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(SigningKey)
	// need to figure out a better way to handle when this function fails to create a token
	// at the moment it returns a 500
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("token creation failed, try again"))
	}

	fmt.Println(tokenString)
	return tokenString
}

// JWTFromCookie decodes the cookie and returns the jwt
func JWTFromCookie(r *http.Request) (string, error) {

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
