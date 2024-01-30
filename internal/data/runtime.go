package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {
	jsonVal := fmt.Sprintf("%d mins", r)

	quotedJSONValue := strconv.Quote(jsonVal)
	return []byte(quotedJSONValue), nil
}
