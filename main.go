// Command goaddimport adds an import spec to a Go source file.
// It reads the input file from standard in, and writes the result
// to standard out.
package main

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/imports"
)

type importer struct {
	ip  string
	in  io.Reader
	out io.Writer
}

func main() {
	if len(os.Args) != 2 || os.Args[1] == "help" || os.Args[1] == "-h" {
		fmt.Fprintf(os.Stderr, "Usage: %s import_path\nInput is read from stdin and output is written to stdout.\n", os.Args[0])
		os.Exit(1)
	}
	imp := importer{
		ip:  os.Args[1],
		in:  os.Stdin,
		out: os.Stdout,
	}
	if err := imp.Run(); err != nil {
		log.Fatal(err)
	}
}

func (i importer) Run() error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "file.go", i.in, parser.ParseComments)
	if err != nil {
		return err
	}

	astutil.AddImport(fset, f, i.ip)

	// Invoke goimports code in "format only" mode to get proper import sorting.
	// There's a little bit of duplicated work here, but this is preferrable to
	// pulling in unexported code from goimports.
	buf := &bytes.Buffer{}
	err = printer.Fprint(buf, fset, f)
	if err != nil {
		return err
	}

	out, err := imports.Process("file.go", buf.Bytes(),
		&imports.Options{FormatOnly: true, Comments: true, TabIndent: true, TabWidth: 8})
	if err != nil {
		return err
	}

	_, err = i.out.Write(out)
	return err
}
