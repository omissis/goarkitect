package file_test

import (
	"goarkitect/internal/arch/file"
	fe "goarkitect/internal/arch/file/except"
	fs "goarkitect/internal/arch/file/should"
	ft "goarkitect/internal/arch/file/that"
	"os"
	"path/filepath"
	"testing"
)

func Test_It_Checks_A_File_Exists(t *testing.T) {
	filename := "./test/project/Makefile"

	vs := file.One(filename).
		Should(fs.Exist()).
		Because("it is needed to encapsulate common project operations")

	if vs != nil {
		for _, v := range vs {
			t.Errorf("%v", v)
		}
	}
}

func Test_It_Checks_Multiple_Files_Exists_When_They_Do(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	vs := file.All().
		That(ft.AreInFolder(filepath.Join(basePath, "test/project"), false)).
		Should(fs.Exist()).
		AndShould(fs.EndWith("file")).
		Because("they are needed for developing and deploying the project")

	if vs != nil {
		for _, v := range vs {
			t.Errorf("%v", v)
		}
	}
}

func Test_It_Checks_Multiple_Files_Exists_When_They_Dont(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	vs := file.All().
		That(ft.AreInFolder(filepath.Join(basePath, "test/project"), false)).
		Should(fs.Exist()).
		AndShould(fs.EndWith(".yaml")).
		Because("they are needed for developing and deploying the project")

	if vs == nil {
		t.Errorf("Expected violations, got nil")
	}

	if len(vs) != 2 {
		t.Errorf("Expected 2 violations, got %d", len(vs))
	}

	if vs[0].String() != "file's name 'Dockerfile' does not end with '.yaml'" {
		t.Errorf("Expected violation for Dockerfile, got: '%s'", vs[0].String())
	}

	if vs[1].String() != "file's name 'Makefile' does not end with '.yaml'" {
		t.Errorf("Expected violation for Makefile, got: '%s'", vs[1].String())
	}
}

func Test_It_Checks_Multiple_Files_Exists_When_They_Do_With_Exception(t *testing.T) {
	basePath, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	vs := file.All().
		That(ft.AreInFolder(filepath.Join(basePath, "test/project"), false)).
		Should(fs.Exist()).
		AndShould(fs.StartWith("Docker")).
		Except(fe.This("Makefile"), fe.This("Something")).
		Because("they are needed for developing and deploying the project")

	if vs != nil {
		for _, v := range vs {
			t.Errorf("%v", v)
		}
	}
}
