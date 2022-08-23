package logx

import (
	"fmt"
	"log"
	"time"

	"github.com/omissis/goarkitect/internal/jsonx"
)

var format string = "text"

func SetFormat(f string) {
	if f != "text" && f != "json" {
		log.Fatal(
			fmt.Errorf("format '%s' is not valid, accepted options are: json, text", f),
		)
	}

	format = f
}

func Fatal(v error) {
	switch format {
	case "text":
		log.Fatal(v)
	case "json":
		log.Fatal(
			jsonx.Marshal(
				map[string]any{
					"time":  time.Now().Format(time.RFC3339),
					"level": "ERROR",
					"msg":   v.Error(),
				},
			),
		)
	default:
		log.Fatalf("invalid format value: '%s'", format)
	}
}
