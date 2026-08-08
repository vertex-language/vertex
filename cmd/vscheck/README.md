# vscheck

`vscheck` is a minimal, no-emit front-end driver for [Vertex](https://github.com/vertex-language/vertex). It scans and parses each `*.vs` file given on the command line and reports whatever diagnostics the scanner/parser produced.

It doesn't resolve names, type-check, or generate code — it only runs source through `scanner` → `parser` → `ast` and logs the result. Think `go vet` or `tsc --noEmit`, scoped to the front end only.

## Install

```
GOPROXY=direct go install github.com/vertex-language/vertex/cmd/vscheck@latest
```

## Usage

```
vscheck file.vs [file.vs ...]
```

Flags:

| Flag         | Description                                             |
|--------------|----------------------------------------------------------|
| `-dump`      | Dump the parsed AST for each file (`ast.Fdump`)          |
| `-comments`  | Retain comments while parsing (`parser.ParseComments`)   |

Examples:

```
vscheck main.vs other.vs
vscheck -dump -comments main.vs
```

## Exit status

`vscheck` exits `0` if every file was read and parsed (diagnostics from the parser itself don't count as failure — a partial parse is still a successful run of the tool). It exits `1` if any file couldn't be read, and `2` if no files were given.

## What "ok" means here

A parse can report diagnostics and still produce a tree — `ParseFile` always returns a non-nil tree, even for broken input, because a partial parse is still a usable one. `vscheck` logs every diagnostic it gets back but doesn't treat their presence as a tool failure; it's a reporting tool, not a gate.
```

And the corresponding rename in the source — `cmd/vscheck/main.go` — with the two lines that reference the binary name updated:

```go
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [flags] file.vs [file.vs ...]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
```

That one already uses `os.Args[0]`, so it doesn't hardcode the name — no other changes needed there. Just move the file to `cmd/vscheck/main.go` so the import path matches what the README installs.