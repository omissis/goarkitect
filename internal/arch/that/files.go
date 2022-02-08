package that

import (
	"goarkitect/internal/arch"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func AreInFolder(folder string, recursive bool) *AreInFolderExpression {
	return &AreInFolderExpression{
		folder:    folder,
		recursive: recursive,
	}
}

type AreInFolderExpression struct {
	folder    string
	recursive bool
}

func (e *AreInFolderExpression) Evaluate(rule *arch.RuleSubjectBuilder) {
	if rule.Kind() != arch.FilesRuleKind {
		panic("AreInFolderExpression can only be used with files rules.")
	}

	if e.recursive {
		rule.Files = e.getFilesRecursive(e.folder)
	} else {
		rule.Files = e.getFiles(e.folder)
	}
}

func (e *AreInFolderExpression) getFilesRecursive(folder string) []string {
	var filenames []string
	err := filepath.Walk(folder, e.visit(&filenames))
	if err != nil {
		log.Println(err)
		return nil
	}
	return filenames
}

func (e *AreInFolderExpression) visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		*files = append(*files, path)
		return nil
	}
}

func (e *AreInFolderExpression) getFiles(folder string) []string {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Println(err)
		return nil
	}

	filenames := make([]string, len(files))
	for i, file := range files {
		filenames[i] = file.Name()
	}

	return filenames
}

func EndsWith(s string) *EndsWithExpression {
	return &EndsWithExpression{
		suffix: s,
	}
}

type EndsWithExpression struct {
	suffix string
}

func (e EndsWithExpression) Evaluate(rule *arch.RuleSubjectBuilder) {
	if rule.Kind() != arch.FilesRuleKind {
		panic("AreInFolderExpression can only be used with files rules.")
	}

	files := make([]string, 0)
	for _, f := range rule.Files {
		if strings.HasSuffix(f, e.suffix) {
			files = append(files, f)
		}
	}
	rule.Files = files
}
