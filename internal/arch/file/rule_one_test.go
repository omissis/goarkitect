package file_test

import (
	"goarkitect/internal/arch/file"
	fs "goarkitect/internal/arch/file/should"
	"goarkitect/internal/arch/rule"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_It_Checks_A_File_Exists(t *testing.T) {
	testCases := []struct {
		desc           string
		filename       string
		wantViolations []rule.Violation
	}{
		{
			desc:           "check that a file exists when it's actually there",
			filename:       "./test/project/Makefile",
			wantViolations: nil,
		},
		{
			desc:     "check that a file exists when it's not there",
			filename: "./test/project/NotExistingFile",
			wantViolations: []rule.Violation{
				rule.NewViolation("file 'NotExistingFile' does not exist"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.One(tC.filename).
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

func Test_It_Checks_A_File_Does_Not_Exist(t *testing.T) {
	testCases := []struct {
		desc           string
		filename       string
		wantViolations []rule.Violation
	}{
		{
			desc:           "check that a file does not exist",
			filename:       "./test/project/NotExistingFile",
			wantViolations: nil,
		},
		{
			desc:     "check that a file does not exist when it's actually there",
			filename: "./test/project/Makefile",
			wantViolations: []rule.Violation{
				rule.NewViolation("file 'Makefile' does exist"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.One(tC.filename).
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

func Test_It_Checks_A_File_Name_Matches_A_Regexp(t *testing.T) {
	testCases := []struct {
		desc           string
		filename       string
		regexp         string
		wantViolations []rule.Violation
	}{
		{
			desc:           "check that a file name matches a regex when it's actually there",
			filename:       "./test/project/Makefile",
			regexp:         "[a-zA-Z0-9]+",
			wantViolations: nil,
		},
		{
			desc:           "check that a file name matches a regex when it's not actually there",
			filename:       "./test/project/NotExistingFile",
			regexp:         "[a-zA-Z0-9]+",
			wantViolations: nil,
		},
		{
			desc:     "check that a file name does not match a regex when it's actually there",
			filename: "./test/project/Makefile",
			regexp:   "[0-9]+",
			wantViolations: []rule.Violation{
				rule.NewViolation("file's name 'Makefile' does not match regex '[0-9]+'"),
			},
		},
		{
			desc:     "check that a file name does not match a regex when it's not actually there",
			filename: "./test/project/NotExistingFile",
			regexp:   "[0-9]+",
			wantViolations: []rule.Violation{
				rule.NewViolation("file's name 'NotExistingFile' does not match regex '[0-9]+'"),
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			vs, errs := file.One(tC.filename).
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
