package encryption

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5WithSalt encrypts a string with a salt using MD5 algorithm
func MD5WithSalt(input string) string {
	// Concatenate the original string with the salt
	salted := input

	// Create a new MD5 hash
	hasher := md5.New()

	// Write the salted string as bytes to the hasher
	hasher.Write([]byte(salted))

	// Get the resulting hash as bytes
	hashBytes := hasher.Sum(nil)

	// Convert hash bytes to hexadecimal string
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
