package server

import (
	"crypto/sha256"
	"hash"
)

var (
	KeyHashAlgo func() hash.Hash = sha256.New
)

const (
	HttpAddr          string = ":8080"
	HttpsAddr         string = ":9090"
	SqlUser           string = "root"
	SqlPass           string = ""
	SqlDb             string = "cloudreader"
	UsernameMaxLength int    = 16
	SaltLength        int    = 16
	KeyHashLength     int    = 16
	KeyHashIterations int    = 250000
	BookHashLength    int    = 32
	PathMaxLength     int    = 100
)
