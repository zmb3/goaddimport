# goaddimport

A simple utility for editors to add imports to Go source files.

Several Go editors have commands for automatically adding import statements to Go files.
This allows users to quickly add an import without scrolling to the top of the file
and potentially losing context.

Implementing this logic in a language other than Go can be tricky:

If the source file doesn't yet have any imports, a single line import declaration
should be chosen:

```
import "fmt"
```

If the source file already contains imports, a multi-line import declaration is preferred:

```
import (
        "fmt"
        "io"
)
```

Additionally imports should be sorted, and standard library imports should be grouped
separately from third-party imports (a la `goimports`).

`goaddimport` is a small tool that handles these cases for you.

## Usage:

The tool accepts the input file on stdin and writes its output to stdout.
It takes a single argument - the import path that should be added.

`$ goaddimport io/ioutil < main.go`
