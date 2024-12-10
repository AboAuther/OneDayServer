package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	logger "github.com/sirupsen/logrus"
)

const HexPrefix = "0x"

func GenerateUUID() string {
	uid, err := uuid.NewUUID()
	if err != nil {
		logger.Errorf("new uuid failed, err: %s", err)
		return uuid.NewString()
	}
	return uid.String()
}

func FormatHexString(input string) string {
	if strings.HasPrefix(input, HexPrefix) {
		return input
	} else {
		return HexPrefix + input
	}
}

func CleanHexString(input string) string {
	if !strings.HasPrefix(input, HexPrefix) {
		return input
	} else {
		return strings.TrimPrefix(input, HexPrefix)
	}
}

func GenerateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	min := int64(1)
	for i := 1; i < length; i++ {
		min *= 10
	}
	max := min*10 - 1
	return fmt.Sprintf("%0*d", length, rand.Int63n(max-min+1)+min)
}
