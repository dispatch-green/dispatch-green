package server

import "embed"

// content holds our static web server content.
//go:embed web/pages/* web/widgets/*
var Templates embed.FS

//go:embed css/*
var Static embed.FS
