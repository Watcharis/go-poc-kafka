package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

func RandomTxnRef() string {
	min := 10000000
	max := 99999999
	randomNumber := rand.Intn(max-min) + min

	return fmt.Sprintf("%d%d", randomNumber, time.Now().UnixNano())
}

// generate defualt uuid ex.8855b6f3-19eb-4790-88d9-a464a18c90aa
func GenerateUUID() string {
	uid := uuid.NewString()
	return uid
}

// generate modified uuid ex.X083FC7E7E25E4EBBB829E53BDCC2BBB9
func GenerateUUIDModified(prefix string) string {
	uid := uuid.NewString()
	return prefix + strings.ToUpper((strings.ReplaceAll(uid, "-", "")))
}

const letterAlphaNumeric = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandAlphaNumeric(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterAlphaNumeric[rand.Intn(len(letterAlphaNumeric))]
	}
	return string(b)
}

const letterNumeric = "0123456789"

func RandNumeric(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterNumeric[rand.Intn(len(letterNumeric))]
	}
	return string(b)
}
