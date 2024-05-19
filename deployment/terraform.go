package deployment

import "embed"

//go:embed all:*.tf
//go:embed all:*.hcl
var Templates embed.FS
