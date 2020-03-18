package main

import (
	"bytes"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/samuelmeuli/nbtohtml"
)

type convertCmd struct {
	Path string `kong:"arg,name='path',help='Jupyter Notebook file to convert.',type='existingfile'"`
}

var (
	// Populated by GoReleaser
	version = "?"

	description = `
nbtohtml is a library for converting Jupyter Notebook files to HTML.
`
	cli struct {
		Convert convertCmd       `kong:"cmd,help='Convert Jupyter Notebook file to HTML.'"`
		Version kong.VersionFlag `kong:"cmd,help='Show version.'"`
	}
)

func (r *convertCmd) Run() error {
	var notebookPath = r.Path

	// Convert notebook file to HTML and print result
	notebookHTML := new(bytes.Buffer)
	err := nbtohtml.ConvertFile(notebookHTML, notebookPath)
	if err != nil {
		return err
	}

	fmt.Println(notebookHTML)
	return nil
}

func main() {
	ctx := kong.Parse(&cli, kong.Description(description), kong.Vars{
		"version": version,
	})
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
