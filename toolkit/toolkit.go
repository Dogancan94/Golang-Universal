package toolkit

import (
	"crypto/rand"
	"fmt"
)

const randomStringSource = "abcdefghijklmnoprstuvxyzABCDEFGHIJKLMNOPRSTUVXYZ0123456789_+"

type Tools struct{}

func (tool *Tools) CreateRandomString(number int) string {
	s, r := make([]rune, number), []rune(randomStringSource)
	for i := range s {
		prime, err := rand.Prime(rand.Reader, len(r))
		if err != nil {
			fmt.Println(err)
		}
		x, y := prime.Uint64(), uint64(len(r))
		s[i] = r[x%y]
	}

	return string(s)
}
