package server

import (
	"golang.org/x/crypto/pbkdf2"
	"net/http"
	"strconv"
)

func HashAndSaltPassword(password, salt []byte) []byte {
	return pbkdf2.Key(password, salt, KeyHashIterations, KeyHashLength, KeyHashAlgo)
}

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
