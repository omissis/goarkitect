package file_test

import (
	"goarkitect/internal/arch/file"
	fs "goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_It_Checks_A_Set_Of_Files_Exist(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		filenames      []string
		wantViolations []rule.Violation
	}{
		{
			desc: "check that a set of files exists when it's actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project/Dockerfile"),
				filepath.Join(basePath, "test/project/Makefile"),
			},
			wantViolations: nil,
		},
		{
			desc: "check that a set of files exists when it's not there",
			filenames: []string{
				filepath.Join(basePath, "test/project/Foofile"),
				filepath.Join(basePath, "test/project/Barfile"),
			},
			wantViolations: []rule.Violation{
				rule.NewViolation("file 'Foofile' does not exist"),
				rule.NewViolation("file 'Barfile' does not exist"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.Set(tC.filenames...).
				Should(fs.Exist()).
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

func Test_It_Checks_A_Set_Of_Files_Do_Not_Exist(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		filenames      []string
		wantViolations []rule.Violation
	}{
		{
			desc: "check that a set of files does not exist when the files are not actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project/Foofile"),
				filepath.Join(basePath, "test/project/Barfile"),
			},
			wantViolations: nil,
		},
		{
			desc: "check that a set of files does not exist when the files are actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project/Dockerfile"),
				filepath.Join(basePath, "test/project/Makefile"),
			},
			wantViolations: []rule.Violation{
				rule.NewViolation("file 'Dockerfile' does exist"),
				rule.NewViolation("file 'Makefile' does exist"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.Set(tC.filenames...).
				Should(fs.NotExist()).
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

func Test_It_Checks_A_Set_Of_Files_Names_Matches_A_Regexp(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		filenames      []string
		regexp         string
		wantViolations []rule.Violation
	}{
		{
			desc: "check that a set of files' names match a regex when it's actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project/Dockerfile"),
				filepath.Join(basePath, "test/project/Makefile"),
			},
			regexp:         "[a-zA-Z0-9]+",
			wantViolations: nil,
		},
		{
			desc: "check that a set of files' names match a regex when it's not actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project/Foofile"),
				filepath.Join(basePath, "test/project/Barfile"),
			},
			regexp:         "[a-zA-Z0-9]+",
			wantViolations: nil,
		},
		{
			desc: "check that a set of files' names do not match a regex when it's actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project/Dockerfile"),
				filepath.Join(basePath, "test/project/Makefile"),
			},
			regexp: "[0-9]+",
			wantViolations: []rule.Violation{
				rule.NewViolation("file's name 'Dockerfile' does not match regex '[0-9]+'"),
				rule.NewViolation("file's name 'Makefile' does not match regex '[0-9]+'"),
			},
		},
		{
			desc: "check that a set of files' names do not match a regex when it's not actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project/Foofile"),
				filepath.Join(basePath, "test/project/Barfile"),
			},
			regexp: "[0-9]+",
			wantViolations: []rule.Violation{
				rule.NewViolation("file's name 'Foofile' does not match regex '[0-9]+'"),
				rule.NewViolation("file's name 'Barfile' does not match regex '[0-9]+'"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.Set(tC.filenames...).
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

func Test_It_Checks_A_Set_Of_Files_Names_Matches_A_Glob_Pattern(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		desc           string
		filenames      []string
		glob           string
		wantViolations []rule.Violation
	}{
		{
			desc: "check that a set of files' names match a glob pattern when it's actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project3/baz.go"),
				filepath.Join(basePath, "test/project3/quux.go"),
			},
			glob:           "*/*/*.go",
			wantViolations: nil,
		},
		{
			desc: "check that a set of files' names match a glob pattern when it's not actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project3/baz.ts"),
				filepath.Join(basePath, "test/project3/quux.ts"),
			},
			glob:           "*/*/*.ts",
			wantViolations: nil,
		},
		{
			desc: "check that a set of files' names do not match a glob pattern when it's actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project3/baz.go"),
				filepath.Join(basePath, "test/project3/quux.go"),
			},
			glob: "*/*/*.ts",
			wantViolations: []rule.Violation{
				rule.NewViolation("file's path 'baz.go' does not match glob pattern '*/*/*.ts'"),
				rule.NewViolation("file's path 'quux.go' does not match glob pattern '*/*/*.ts'"),
			},
		},
		{
			desc: "check that a set of files' names do not match a glob pattern when it's not actually there",
			filenames: []string{
				filepath.Join(basePath, "test/project3/baz.ts"),
				filepath.Join(basePath, "test/project3/quux.ts"),
			},
			glob: "*/*/*.go",
			wantViolations: []rule.Violation{
				rule.NewViolation("file's path 'baz.ts' does not match glob pattern '*/*/*.go'"),
				rule.NewViolation("file's path 'quux.ts' does not match glob pattern '*/*/*.go'"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.Set(tC.filenames...).
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
