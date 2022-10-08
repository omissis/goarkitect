package santhosh_test

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/omissis/goarkitect/internal/schema/santhosh"
)

func TestLoadSchema(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		desc       string
		basePath   string
		wantErr    bool
		wantErrMsg string
	}{
		{
			desc:       "wrong basepath",
			basePath:   filepath.Join(os.TempDir(), strconv.Itoa(rand.Int()), "not-existiing"),
			wantErr:    true,
			wantErrMsg: "no such file or directory",
		},
		{
			desc:       "broken schema",
			basePath:   "testdata/loader-broken-schema",
			wantErr:    true,
			wantErrMsg: "failed to add resource to json schema compiler",
		},
		{
			desc:       "bad schema",
			basePath:   "testdata/loader-bad-schema",
			wantErr:    true,
			wantErrMsg: "failed to compile json schema",
		},
		{
			desc:     "good schema",
			basePath: "testdata/loader-good-schema",
		},
	}
	for _, tC := range testCases {
		tC := tC

		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			_, err := santhosh.LoadSchema(tC.basePath)
			if (err != nil) != tC.wantErr {
				t.Errorf("LoadSchema() error = %v, wantErr %v", err, tC.wantErr)

				return
			}

			if tC.wantErr && err != nil && !strings.Contains(err.Error(), tC.wantErrMsg) {
				t.Errorf("LoadSchema() gotErrMsg = %s, wantErrMsg %s", err, tC.wantErrMsg)
			}
		})
	}
}
