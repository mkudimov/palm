package tokengenerator

import (
	"math/rand"
	"os"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyz"

func Generate(tokenLength int, tokenNumber int64, out string) error {
	tokens := make([]string, tokenNumber)
	for i := int64(0); i < tokenNumber; i++ {
		tokens[i] = generateToken(tokenLength)
	}
	return dump(tokens, out)
}

func generateToken(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func dump(tokens []string, out string) error {
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, t := range tokens {
		_, err := f.WriteString(t + "\n")
		if err != nil {
			return err
		}
	}
	return nil
}
