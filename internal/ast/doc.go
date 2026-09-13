// Package ast defines the abstract syntax tree node types produced by the
// parser from Radic source code.
//
// Every node implements one of three marker interfaces depending on its
// role in the tree: Expr for expressions, Stmt for statements, and Decl
// for top-level declarations.
package ast
