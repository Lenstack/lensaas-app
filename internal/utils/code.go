package utils

import (
	"math/rand"
)

func NewCode() int64 {
	rand.Float64()
	max := 9999999
	min := 1000000
	code := min + rand.Intn(max-min)
	return int64(code)
}
