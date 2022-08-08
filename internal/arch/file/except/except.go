package except

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
	"path/filepath"
)

type Expression struct {
	value string
}

func (e Expression) Evaluate(rb rule.Builder) {
	frb := rb.(*file.RuleBuilder)

	nf := make([]string, 0)
	for _, f := range frb.GetFiles() {
		if filepath.Base(f) != e.value {
			nf = append(nf, f)
		}
	}

	frb.SetFiles(nf)

	return
}
