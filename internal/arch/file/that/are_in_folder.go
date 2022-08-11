package that

import (
	"goarkitect/internal/arch/file"
	"goarkitect/internal/arch/rule"
	"io/ioutil"
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
	var files []string
	var err error

	if e.recursive {
		files, err = e.getFilesRecursive(e.folder)
	} else {
		files, err = e.getFiles(e.folder)
	}

	if err != nil {
		rb.(*file.RuleBuilder).AddError(err)
	}

	rb.(*file.RuleBuilder).SetFiles(files)
}

func (e *AreInFolderExpression) getFilesRecursive(folder string) ([]string, error) {
	var filenames []string
	if err := filepath.Walk(folder, e.visit(&filenames)); err != nil {
		return nil, err
	}

	return filenames, nil
}

func (e *AreInFolderExpression) visit(files *[]string) filepath.WalkFunc {
	return func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if file.IsDir() {
			return nil
		}

		*files = append(*files, path)

		return nil
	}
}

func (e *AreInFolderExpression) getFiles(folder string) ([]string, error) {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	filePaths := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(folder, file.Name()))
		}
	}

	return filePaths, nil
}
