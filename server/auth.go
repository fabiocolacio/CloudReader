package server

//auth.go is the module used to implement User login authentication
import (
	"golang.org/x/crypto/pbkdf2"
	"net/http"
	"strconv"
)


//HashAndSaltPassword take in a password string and array of bytes
//and return the hashed password concatenated with array of bytes
func HashAndSaltPassword(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, KeyHashIterations, KeyHashLength, KeyHashAlgo)
}

//VerifyUser verify that the user login session is still valid
//everytime it is called. It takes in an http request as an argument,
//and return the user id to verify if the user login session is valid,
//return an error otherwise.
func VerifyUser(req *http.Request) int {
	cookie, err := req.Cookie("session")
	if err == http.ErrNoCookie {
		return 0
	}

	uid, _ := strconv.Atoi(cookie.Value)
	if UserExists(uid) {
		return uid
	}

	return 0
}
