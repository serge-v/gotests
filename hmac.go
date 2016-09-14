package main

import (
	"crypto/sha512"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
)

func main() {

	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}

	dk := pbkdf2.Key([]byte("some password"), salt, 4096, 32, sha512.New)

	fmt.Println("hash:", dk)
	fmt.Println("hash:", hex.EncodeToString(dk))
	fmt.Println("salt:", hex.EncodeToString(salt))
}
