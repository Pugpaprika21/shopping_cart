package utils

import (
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func UintFromString(strInput string) uint {
	data, err := strconv.ParseUint(strInput, 10, 0)
	if err != nil {
		return 0
	}
	return uint(data)
}

func MoneyFormat(s float64, useDecimal bool) string {
	formatted := strconv.FormatFloat(s, 'f', 2, 64)
	parts := strings.Split(formatted, ".")
	integral := parts[0]

	decimal := ""
	if useDecimal {
		decimal = "." + parts[1]
	}

	integralWithCommas := ""
	for i, rune := range integral {
		if i != 0 && (len(integral)-i)%3 == 0 {
			integralWithCommas += ","
		}
		integralWithCommas += string(rune)
	}

	return integralWithCommas + decimal
}

func GenerateString() string {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)

	var builder strings.Builder
	for i := 0; i < 6; i++ {
		builder.WriteByte(charset[random.Intn(len(charset))])
	}
	return builder.String()
}
