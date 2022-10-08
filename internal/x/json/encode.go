package json

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Marshal(v ...any) (string, error) {
	ret := make([]string, len(v))

	for i, vv := range v {
		b, err := json.Marshal(vv)
		if err != nil {
			return "", fmt.Errorf("cannot marshal value '%+v': %w", vv, err)
		}

		ret[i] = string(b)
	}

	return strings.Join(ret, " "), nil
}
