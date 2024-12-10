package utils

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var RecvWindow = time.Minute.Milliseconds()

const (
	OneDayApiKey                  = "ONE-WAY-API-KEY"
	OneDayPassphrase              = "ONE-DAY-PASSPHRASE"
	OneDayTimestamp               = "ONE-DAY-TIMESTAMP"
	OneDaySignature               = "ONE-DAY-SIGNATURE"
	OneDayUID                     = "ONE-DAY-UID"
	UserInContext                 = "userInContext"
	OneDayAPIIPRateLimitRemaining = "ONE-DAY-API-IP-RATE-LIMIT-REMAINING"
	UserApiPrefix                 = "/api/v1/oneDay/user"
	OneDayAuthorization           = "Authorization"
)

func SortQueryString(rawQuery string) string {
	if rawQuery == "" {
		return ""
	}

	values, err := url.ParseQuery(rawQuery)
	if err != nil {
		return ""
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var sortedQuery []string
	for _, key := range keys {
		sortedValues := values[key]
		sort.Strings(sortedValues)
		for _, value := range sortedValues {
			sortedQuery = append(sortedQuery, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}

	return strings.Join(sortedQuery, "&")
}

func ParseJWT(tokenString string, secretKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
