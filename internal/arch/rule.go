package arch

import "goarkitect/internal/arch/slice"

type RuleKind int

const (
	FilesRuleKind RuleKind = iota
)

func NewRule() *RuleSubjectBuilder {
	return &RuleSubjectBuilder{}
}

type RuleSubjectBuilder struct {
	kind  RuleKind
	Files []string
	thats []That
	excepts []string
	shoulds []Should
	because string
}

func (rb *RuleSubjectBuilder) init() {
	rb.kind = FilesRuleKind
	rb.Files = []string{}
	rb.thats = []That{}
	rb.excepts = []string{}
	rb.shoulds = []Should{}
	rb.because = ""
}

func (rb *RuleSubjectBuilder) Kind() RuleKind {
	return rb.kind
}

func (rb *RuleSubjectBuilder) AllFiles() *RuleSubjectBuilder {
	rb.init()

	rb.kind = FilesRuleKind

	return rb
}

func (rb *RuleSubjectBuilder) That(t That) *RuleSubjectBuilder {
	rb.thats = []That{t}

	return rb
}

func (rb *RuleSubjectBuilder) AndThat(t That) *RuleSubjectBuilder {
	rb.thats = append(rb.thats, t)

	return rb
}

func (rb *RuleSubjectBuilder) Except(s ...string) *RuleSubjectBuilder {
	rb.excepts = s

	return rb
}

func (rb *RuleSubjectBuilder) Should(e Should) *RuleSubjectBuilder {
	rb.shoulds = []Should{e}

	return rb
}

func (rb *RuleSubjectBuilder) AndShould(e Should) *RuleSubjectBuilder {
	rb.shoulds = append(rb.shoulds, e)

	return rb
}

func (rb *RuleSubjectBuilder) Because(b string) bool {
	rb.because = b

	for _, that := range rb.thats {
		that.Evaluate(rb)
	}

	for _, except := range rb.excepts {
		for j, file := range rb.Files {
			if file == except {
				slice.Remove(rb.Files, j)
			}
		}
	}

	for _, should := range rb.shoulds {
		violations := should.Evaluate(rb)
		if len(violations) > 0 {
			return false
		}
	}

	return true
}
