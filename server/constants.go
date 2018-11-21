package server

import (
  "hash"
  "crypto/sha256"
)

var (KeyHashAlgo func() hash.Hash  = sha256.New)

const(
    HttpAddr string = ":8080"
    HttpsAddr string = ":9090"
    SqlUser string = "root"
    SqlPass string = ""
    SqlDb string = "cloudreader"
    UsernameMaxLength int = 16
    SaltLength int = 16
    KeyHashLength int = 16
    KeyHashIterations int = 250000
    BookHashLength int = 16
    PathMaxLength int = 100

)
