package arch_test

import (
	"goarkitect/internal/arch"
	"goarkitect/internal/arch/should"
	"goarkitect/internal/arch/that"
	"testing"
)

//func TestRuleFileExist(t *testing.T) {
//	filename := "test/project/Makefile"
//
//	r := arch.NewFileRule(filename).Should(arch.NewExists())
//
//	if r.Validate() != true {
//		t.Errorf("expected file %s to exist, it did not", filename)
//	}
//}

func TestRuleFilesExist(t *testing.T) {
	filenames := []string{"test/project/Makefile", "test/project/Dockerfile"}

	r := arch.NewRule().AllFiles().
		That(that.AreInFolder("test/project", false)).
		AndThat(that.EndsWith("file")).
		Except("Makefile").
		Should(should.Exist()).
		Because("some reason")

	if r != true {
		t.Errorf("expected Files %s to exist, they did not", filenames)
	}
}
