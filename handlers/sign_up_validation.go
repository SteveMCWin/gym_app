package handlers

import (
	"crypto/rand"
	"time"
	"log"
	"encoding/binary"
)

var signupTokens map[int]string

func init() {
	signupTokens = make(map[int]string)
}

func CreateToken(user string, valid_duration time.Duration) int {
	var token_val int

	for {
		token_val = generateCode()
		if _, exists := signupTokens[token_val]; exists == false {
			break
		}
	}

	signupTokens[token_val] = user

	timer := time.NewTimer(valid_duration)

	go func(val int) {
		<-timer.C
		if _, exists := signupTokens[val]; exists == true {
			delete(signupTokens, val)
		}
	}(token_val)

	return token_val
}

func generateCode() int {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return int(binary.BigEndian.Uint32(b)%1000000)
}
