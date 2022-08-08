package that

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

func (e *AreInFolderExpression) Evaluate(rb rule.Builder) {
	frb := rb.(*file.RuleBuilder)

	if e.recursive {
		frb.SetFiles(e.getFilesRecursive(e.folder))
	} else {
		frb.SetFiles(e.getFiles(e.folder))
	}
}

func (e *AreInFolderExpression) getFilesRecursive(folder string) []string {
	var filenames []string
	if err := filepath.Walk(folder, e.visit(&filenames)); err != nil {
		log.Println(err)
		return nil
	}

	return filenames
}

func (e *AreInFolderExpression) visit(files *[]string) filepath.WalkFunc {
	return func(path string, _ os.FileInfo, err error) error {
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
		panic(err)
	}

	filePaths := make([]string, len(files))
	for i, file := range files {
		filePaths[i] = filepath.Join(folder, file.Name())
	}

	return filePaths
}
