package utils

import (
	"math/rand"
	"strconv"
	"time"
)

const (
	LowerLetters      = "abcdefghijklmnopqrstuvwxyz"
	LetterLowersBytes = "0123456789abcdefghijklmnopqrstuvwxyz"
	LetterBytes       = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	BaseBytes         = "0123456789abcdef"
)

func RandStringBytes(n int, s string) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = s[rand.Intn(len(s))]
	}
	return string(b)
}

func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn((max+1)-min) + min
}

func RandomPossibility(n int) bool {
	r := RandomInt(0, 100)
	return r < n
}

func GenerateUsername() (username, firstname, lastname string) {
	a, b, c := generateStupidName()
	firstname = a
	lastname = c
	if RandomPossibility(40) {
		username = a + b + c
	} else if RandomPossibility(30) {
		username = a + b + c + strconv.Itoa(RandomInt(1, 99))
	} else if RandomPossibility(22) {
		username = a + b + c + strconv.Itoa(RandomInt(2020, 2020))
	} else if RandomPossibility(22) {
		username = c + b + strconv.Itoa(RandomInt(1, 2020)) + a
	} else if RandomPossibility(34) {
		username = a + strconv.Itoa(RandomInt(1, 2020)) + b + c
	} else {
		username = a + b + c + strconv.Itoa(RandomInt(111, 1111))
	}
	return
}
