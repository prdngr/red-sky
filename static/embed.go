package static

import "embed"

//go:embed terraform/*.tf
//go:embed terraform/*.hcl
var Embeds embed.FS
