package web

import "embed"

//go:embed index.html ideas.html css/* js/* img/*
var FS embed.FS
