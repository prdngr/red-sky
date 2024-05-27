package static

import "embed"

//go:embed config.yaml
//go:embed terraform/*.tf
//go:embed terraform/*.hcl
var Embeds embed.FS
