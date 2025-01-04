package main

import "embed"

//go:embed files/*
var files embed.FS

//indexData, _ := files.ReadFile("index.html")
