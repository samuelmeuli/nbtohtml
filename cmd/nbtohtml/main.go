package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/samuelmeuli/nbtohtml"
	"os"
)

type convertCmd struct {
	Path string `arg name:"path" help:"Jupyter Notebook file to convert" type:"path"`
}

func (r *convertCmd) Run() error {
	var notebookPath = r.Path

	// Make sure notebook file exists
	if _, err := os.Stat(notebookPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", notebookPath)
		}
		return fmt.Errorf("error while checking whether file %s exists", notebookPath)
	}

	// Convert notebook file to HTML and print result
	html, err := nbtohtml.ConvertFileToHTML(notebookPath)
	if err != nil {
		return err
	}

	fmt.Println(html)
	return nil
}

var cli struct {
	Debug bool `help:"Enable debug mode"`

	Convert convertCmd `cmd help:"Convert Jupyter Notebook file to HTML"`
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
