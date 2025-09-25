package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	max := int64(1)
	for range length {
		max *= 10
	}

	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		return "", err
	}

	otp := fmt.Sprintf("%0*d", length, n.Int64())
	return otp, nil
}
