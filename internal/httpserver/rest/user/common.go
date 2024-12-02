package user

import (
	"time"
)

var nintyDayMillsAgo = time.Hour.Milliseconds() * 24 * 90

const defaultLimit = 500
