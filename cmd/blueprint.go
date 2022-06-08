package cmd

import "embed"

//go:embed templates/*
var FileTemplates embed.FS
