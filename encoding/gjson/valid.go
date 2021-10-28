package gjson

import "github.com/tidwall/gjson"

func IsValid(json string) bool {
	return gjson.Valid(json)
}
