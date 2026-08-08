// cmd/vscheck/main.go
//
// vscheck is a minimal, no-emit front-end driver for Vertex: it scans and
// parses each *.vs file given on the command line and reports whatever
// diagnostics the scanner/parser produced. It never resolves names, never
// type-checks, and never generates code — it only exercises the front end
// (scanner -> parser -> ast) and logs the result.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vertex-language/vertex/ast"
	"github.com/vertex-language/vertex/parser"
	"github.com/vertex-language/vertex/token"
)

func main() {
	var (
		dump    = flag.Bool("dump", false, "dump the parsed AST for each file (ast.Fdump)")
		comment = flag.Bool("comments", false, "retain comments while parsing (parser.ParseComments)")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] file.vs [file.vs ...]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	log.SetFlags(0)

	var mode parser.Mode
	if *comment {
		mode |= parser.ParseComments
	}

	exit := 0
	for _, path := range args {
		if err := checkFile(path, mode, *dump); err != nil {
			log.Printf("%s: %v", path, err)
			exit = 1
		}
	}
	os.Exit(exit)
}

// checkFile reads, scans, and parses a single file, logging any diagnostics
// produced. It returns an error only for problems outside the language
// front end itself (e.g. the file couldn't be read) — parse diagnostics are
// logged, not returned, since a partial parse is still a successful run of
// the tool.
func checkFile(path string, mode parser.Mode, dump bool) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	// NOTE: token.NewFile's exact signature isn't spelled out in the
	// package docs I have — adjust this line if it differs.
	file := token.NewFile(path, src)

	tree, diags := parser.ParseFile(file, mode)
	// parser.ParseFile always returns a non-nil tree, even on a partial
	// parse, so tree.Release() below is always valid to defer.
	defer tree.Release()

	if len(diags) == 0 {
		log.Printf("%s: ok", path)
	} else {
		token.SortDiagnostics(diags)
		for _, d := range diags {
			pos := file.Position(d.Pos)
			log.Printf("%s:%d:%d: %s", path, pos.Line, pos.Col, d.Msg)
		}
		log.Printf("%s: %d diagnostic(s)", path, len(diags))
	}

	if dump {
		if err := ast.Fdump(os.Stdout, tree); err != nil {
			return fmt.Errorf("dump: %w", err)
		}
	}

	return nil
}