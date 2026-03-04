package static

import "embed"

//go:embed terraform/*.tf
//go:embed terraform/modules/*/*.tf
//go:embed terraform/scripts/*.sh
//go:embed terraform/*.hcl
var Embeds embed.FS
