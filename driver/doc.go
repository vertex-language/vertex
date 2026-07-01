// Package driver implements the Vertex compiler pipeline: resolving a
// project's dependency graph via pkg.Load, lowering every module it
// contains through Vertex IR, Machine IR, and native code generation, and
// linking the result into an object file or executable.
//
// A project's own language runtime, if it has one, is not special-cased
// anywhere in this package: it is whatever module-shaped dependency the
// project's own vs.mod happens to declare, resolved and compiled through
// exactly the same pkg.Graph / pipeline.Unit path as every other
// dependency. There is no -packages-dir flag, no reserved "runtime/"
// directory, and no separate runtime-object parameter threaded through
// the linker — that was a stopgap from before pkg.Load existed to resolve
// dependencies at all, and it's gone now that it does. Vertex packages
// are just packages.
package driver