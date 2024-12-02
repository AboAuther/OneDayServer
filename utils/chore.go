package utils

import (
	"strings"

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
