package log

import (
	"errors"
	"fmt"
	"log"
	"time"

	jsonx "github.com/omissis/goarkitect/internal/x/json"
)

var (
	//nolint:gochecknoglobals // it's global to keep logx simple and close to stdlib's
	format = "text"

	ErrUnknownOutputFormat = errors.New("unknown output format, supported ones are: json, text")
)

func SetFormat(f string) {
	if f != "text" && f != "json" {
		log.Fatal(
			fmt.Errorf("'%s' :%w", f, ErrUnknownOutputFormat),
		)
	}

	format = f
}

func Fatal(v error) {
	switch format {
	case "text":
		log.Fatal(v)

	case "json":
		jv, jerr := jsonx.Marshal(
			map[string]any{
				"time":  time.Now().Format(time.RFC3339),
				"level": "ERROR",
				"msg":   v.Error(),
			},
		)
		if jerr != nil {
			log.Fatal(jerr)
		}

		log.Fatal(jv)

	default:
		log.Fatalf("invalid format value: '%s'", format)
	}
}
