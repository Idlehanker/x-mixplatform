package util

import (
	"crypto/sha1"
	"fmt"
	"log"
)

// Config is variable contains App Configuration JSON content
// var Config Configuration
var logger *log.Logger

// P is println function
func P(a ...interface{}) {
	fmt.Println(a)
}

func display() {

}

// Encrypt hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
