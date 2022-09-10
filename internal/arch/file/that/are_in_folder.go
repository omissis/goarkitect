package that

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/omissis/goarkitect/internal/arch/file"
	"github.com/omissis/goarkitect/internal/arch/rule"
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
	errors    []error
}

func (e *AreInFolderExpression) GetErrors() []error {
	return e.errors
}

func (e *AreInFolderExpression) Evaluate(rb rule.Builder) {
	var (
		files []string
		err   error
	)

	frb, ok := rb.(*file.RuleBuilder)
	if !ok {
		e.errors = append(e.errors, file.ErrInvalidRuleBuilder)

		return
	}

	if e.recursive {
		files, err = e.getFilesRecursive(e.folder)
	} else {
		files, err = e.getFiles(e.folder)
	}

	if err != nil {
		frb.AddError(err)
	}

	frb.SetFiles(files)
}

func (e *AreInFolderExpression) getFilesRecursive(folder string) ([]string, error) {
	var filenames []string
	if err := filepath.Walk(folder, e.visit(&filenames)); err != nil {
		return nil, fmt.Errorf("error walking the path '%s': %w", folder, err)
	}

	return filenames, nil
}

func (e *AreInFolderExpression) visit(files *[]string) filepath.WalkFunc {
	return func(path string, file os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !file.IsDir() {
			*files = append(*files, path)
		}

		return nil
	}
}

func (e *AreInFolderExpression) getFiles(folder string) ([]string, error) {
	files, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("error getting files in folder '%s': %w", folder, err)
	}

	filePaths := make([]string, 0)

	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, filepath.Join(folder, file.Name()))
		}
	}

	return filePaths, nil
}
