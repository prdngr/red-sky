package static

import "embed"

//go:embed .initialized
//go:embed terraform/*.tf
//go:embed terraform/*.hcl
var Embeds embed.FS
