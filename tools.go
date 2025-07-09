// +build tools

// This file imports packages that are required for running tools such as migrations,
// but aren't directly used in the application code. This is only necessary for go mod to retain them.

package tools

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

