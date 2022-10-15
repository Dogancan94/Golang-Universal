package helpers

import (
	"log"
	"math/rand"
	"time"
)

func InfoLog(message interface{}) {
	log.Println(message)
}

func CreateRandomNumber(n int) int {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(n)
	return value
}
