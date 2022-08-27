package jsonx

import (
	"encoding/json"
	"log"
	"strings"
)

func Marshal(v ...any) string {
	ret := make([]string, len(v))
	for i, vv := range v {
		b, err := json.Marshal(vv)
		if err != nil {
			log.Fatal(err)
		}

		ret[i] = string(b)
	}

	return strings.Join(ret, " ")
}
