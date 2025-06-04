package ui

import "embed"

//go:embed static
var StaticFS embed.FS

//go:embed static/robots.txt
var RobotsTXT []byte
