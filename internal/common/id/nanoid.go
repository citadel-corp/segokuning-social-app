package id

import (
	gonanoid "github.com/matoous/go-nanoid"
)

const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateStringID(n int) string {
	return gonanoid.MustGenerate(chars, n)
}
