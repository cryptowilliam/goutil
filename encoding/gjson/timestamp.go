package gjson

import (
	"fmt"
	"time"
)

// Json index friendly timestamp.
// Why choose format time as nanosecond int64, but not human readable string?
// time.RFC3339Nano is the only format which can format time with nanosecond accuracy,
// but it has timezone information in it, when you need create unique index for timestamp in json/jsonb database,
// same epoch timestamps in different timezone may be considered as different string/index.
type JifTimestamp time.Time

func (t JifTimestamp) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("%d", time.Time(t).Nanosecond())
	return []byte(stamp), nil
}
