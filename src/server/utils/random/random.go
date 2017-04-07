package random

import (
	"math/rand"
	"time"
)

func RandomNum(len int)[]byte{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := []byte{}
	for i:=0; i<len; i++ {
		num = append(num, byte(r.Intn(10)))
	}
	return num
}
