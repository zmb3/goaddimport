package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestInvalidImportPath(t *testing.T) {

}

func TestAddExistingImport(t *testing.T) {
	tests := []struct {
		testCase, importPath string
	}{
		{"3", "fmt"},
		{"5", "github.com/user1/pkg1"},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			orig := &bytes.Buffer{}
			outBuf := &bytes.Buffer{}
			f, err := os.Open(filepath.Join("testdata", test.testCase))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			i := importer{
				in:  io.TeeReader(f, orig),
				out: outBuf,
				ip:  test.importPath,
			}
			err = i.Run()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(orig.Bytes(), outBuf.Bytes()) {
				t.Errorf("test %s output does not match original:\n\n%s", test.testCase, outBuf.String())
			}
		})
	}
}

func TestAddImport(t *testing.T) {
	tests := []struct {
		name, description string
	}{
		{"1", "no imports"},
		{"2", "no imports (comments in file)"},
		{"3", "single stdlib import"},
		{"4", "multiple stdlib imports"},
		{"5", "multiple third party"},
		{"6", "variety"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s: %s", test.name, test.description), func(t *testing.T) {
			buf := &bytes.Buffer{}
			f, err := os.Open(filepath.Join("testdata", test.name))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			i := importer{
				in:  f,
				out: buf,
				ip:  "github.com/someuser/somepkg",
			}
			err = i.Run()
			if err != nil {
				t.Fatal(err)
			}
			golden, err := ioutil.ReadFile(filepath.Join("testdata", test.name+".golden"))
			if err != nil {
				t.Fatalf("Can open golden file %s: %v", test.name, err)
			}
			if !bytes.Equal(golden, buf.Bytes()) {
				t.Errorf("test %s output does not match golden:\n\n%s", test.name, buf.String())
			}
		})
	}
}
