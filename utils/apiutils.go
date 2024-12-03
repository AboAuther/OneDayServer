package utils

import (
	"net/url"
	"sort"
	"strings"
	"time"
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
