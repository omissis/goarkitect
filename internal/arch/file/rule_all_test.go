package file_test

import (
	"goarkitect/internal/arch/file"
	fs "goarkitect/internal/arch/file/should"
	ft "goarkitect/internal/arch/file/that"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_It_Checks_All_Files_In_A_Folder_Start_With(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		folder         string
		wantViolations []rule.Violation
	}{
		{
			desc:           "check that all files in a folder start with the given prefix when they do",
			folder:         filepath.Join(basePath, "test/project2"),
			wantViolations: nil,
		},
		{
			desc:   "check that all files in a folder start with the given prefix when they don't",
			folder: filepath.Join(basePath, "test/project"),
			wantViolations: []rule.Violation{
				rule.NewViolation("file's name 'Makefile' does not start with 'Docker'"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.All().
				That(ft.AreInFolder(tC.folder, false)).
				Should(fs.StartWith("Docker")).
				Because("testing reasons")

			if !cmp.Equal(vs, tC.wantViolations, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("Expected %v, got %v", tC.wantViolations, vs)
			}

			if errs != nil {
				t.Errorf("Expected errs to be nil, got: %+v", errs)
			}
		})
	}
}

func Test_It_Checks_All_Files_In_A_Folder_End_With(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		folder         string
		wantViolations []rule.Violation
	}{
		{
			desc:           "check that all files in a folder end with the given suffix when they do",
			folder:         filepath.Join(basePath, "test/project"),
			wantViolations: nil,
		},
		{
			desc:   "check that all files in a folder end with the given suffix when they don't",
			folder: filepath.Join(basePath, "test/config"),
			wantViolations: []rule.Violation{
				rule.NewViolation("file's name 'base.yml' does not end with 'file'"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.All().
				That(ft.AreInFolder(tC.folder, false)).
				Should(fs.EndWith("file")).
				Because("testing reasons")

			if !cmp.Equal(vs, tC.wantViolations, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("Expected %v, got %v", tC.wantViolations, vs)
			}

			if errs != nil {
				t.Errorf("Expected errs to be nil, got: %+v", errs)
			}
		})
	}
}

func Test_It_Checks_All_Files_Names_In_A_Folder_Match_A_Regexp(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		folder         string
		regexp         string
		wantViolations []rule.Violation
	}{
		{
			desc:           "check that all files' names in a folder match a regex",
			folder:         filepath.Join(basePath, "test/project"),
			regexp:         "[a-zA-Z0-9]+",
			wantViolations: nil,
		},
		{
			desc:   "check that all files' names in a folder do not match a regex",
			folder: filepath.Join(basePath, "test/project"),
			regexp: "[0-9]+",
			wantViolations: []rule.Violation{
				rule.NewViolation("file's name 'Dockerfile' does not match regex '[0-9]+'"),
				rule.NewViolation("file's name 'Makefile' does not match regex '[0-9]+'"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.All().
				That(ft.AreInFolder(tC.folder, false)).
				Should(fs.MatchRegex(tC.regexp)).
				Because("testing reasons")

			if !cmp.Equal(vs, tC.wantViolations, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("Expected violations %v, got %v", tC.wantViolations, vs)
			}

			if errs != nil {
				t.Errorf("Expected errs to be nil, got: %+v", errs)
			}
		})
	}
}

func Test_It_Checks_All_Files_Paths_In_A_Folder_Match_A_Glob_Pattern(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		folder         string
		glob           string
		wantViolations []rule.Violation
	}{
		{
			desc:           "check that all files' names in a folder match a glob pattern",
			folder:         filepath.Join(basePath, "test/project3"),
			glob:           "*/*/*.go",
			wantViolations: nil,
		},
		{
			desc:   "check that all files' names in a folder do not match a glob pattern",
			folder: filepath.Join(basePath, "test/project3"),
			glob:   "*/*/*.ts",
			wantViolations: []rule.Violation{
				rule.NewViolation("file's path 'baz.go' does not match glob pattern '*/*/*.ts'"),
				rule.NewViolation("file's path 'quux.go' does not match glob pattern '*/*/*.ts'"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.All().
				That(ft.AreInFolder(tC.folder, false)).
				Should(fs.MatchGlob(tC.glob, basePath)).
				Because("testing reasons")

			if !cmp.Equal(vs, tC.wantViolations, cmp.AllowUnexported(rule.Violation{}), cmpopts.EquateEmpty()) {
				t.Errorf("Expected violations %v, got %v", tC.wantViolations, vs)
			}

			if errs != nil {
				t.Errorf("Expected errs to be nil, got: %+v", errs)
			}
		})
	}
}
