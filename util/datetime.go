package util

import "time"

const (
	DATETIME_FORMAT_PRICE_DUE_IN_DB     = "2006-01-02 00:00:00"
	DATETIME_FORMAT_PRICE_DUE_IN_REDIS  = "02012006"
	DATETIME_FORMAT_DATE_IN_RESPONSE_EN = "2 January 2006"
	DATETIME_FORMAT_DATE_IN_SLIP        = "2 January 2006 15:04:05"
	DATETIME_FORMAT_PRICE_DUE           = "2006-01-02"
)

var Now = func() time.Time {
	return time.Now()
}
