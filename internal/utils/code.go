package utils

import (
	"math/rand"
	"strconv"
)

func NewCode() string {
	rand.Float64()
	max := 9999999
	min := 1000000
	code := min + rand.Intn(max-min)
	return strconv.FormatInt(int64(code), 10)
}
