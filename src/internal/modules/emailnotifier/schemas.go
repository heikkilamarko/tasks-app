package emailnotifier

import "embed"

//go:embed schemas/*.json
var SchemasFS embed.FS
