package static

import "embed"

//go:embed terraform/*.tf
//go:embed terraform/*.hcl
var Terraform embed.FS
