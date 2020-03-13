package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/samuelmeuli/nbtohtml"
)

type convertCmd struct {
	Path string `arg name:"path" help:"Jupyter Notebook file to convert." type:"existingfile"`
}

var (
	// Populated by GoReleaser
	version = "?"
	commit  = "?"
	date    = "?"

	description = `
nbtohtml is a library for converting Jupyter Notebook files to HTML.
`
	cli struct {
		Convert convertCmd       `cmd help:"Convert Jupyter Notebook file to HTML."`
		Version kong.VersionFlag `help:"Show version."`
	}
)

func (r *convertCmd) Run() error {
	var notebookPath = r.Path

	// Convert notebook file to HTML and print result
	html, err := nbtohtml.ConvertFileToHTML(notebookPath)
	if err != nil {
		return err
	}

	fmt.Println(html)
	return nil
}

func main() {
	ctx := kong.Parse(&cli, kong.Description(description), kong.Vars{
		"version": fmt.Sprintf("%s-%s-%s", version, commit, date),
	})
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
