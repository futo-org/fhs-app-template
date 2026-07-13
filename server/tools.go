//go:build tools

// This file exists only to pin tooling dependencies in go.mod so that
// `go mod tidy` does not remove them. It is never compiled into the app.
//
// atlas (see atlas.hcl) runs `go run ariga.io/atlas-provider-bun`, whose CLI
// imports every bun dialect. They must stay at the same version as the core
// `bun` module or they fail to compile (schema.Formatter skew).
package tools

import (
	_ "ariga.io/atlas-provider-bun/bunschema"
	_ "github.com/uptrace/bun/dialect/mssqldialect"
	_ "github.com/uptrace/bun/dialect/mysqldialect"
	_ "github.com/uptrace/bun/dialect/oracledialect"
	_ "github.com/uptrace/bun/dialect/pgdialect"
)
