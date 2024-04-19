package utils

import (
	"math/rand"
	"time"
)

var MyRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomInt(min, max int64) int64 {
	return min + MyRand.Int63n(max-min+1)
}

func RandomsString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[MyRand.Intn(len(letters))]
	}
	return string(b)
}

// RandomOwner generate a random owner name
func RandomOwner() string {
	if MyRand.Intn(2) == 0 {
		return maleNames[MyRand.Intn(len(maleNames))]
	} else {
		return femaleNames[MyRand.Intn(len(femaleNames))]
	}
}

// RandomMoney generate a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency generate a random currency code
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[MyRand.Intn(n)]
}
