package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"os"
)

type context struct {
	Debug bool
}

type convertCmd struct {
	IncludeCodeCSS     bool   `help:"Whether styles for syntax highlighting in code should be included in the HTML" default:"true"`
	IncludeNotebookCSS bool   `help:"Whether styles for the Jupyter Notebook (e.g. cells and Markdown) should be included in the HTML" default:"true"`
	CodeLightStyle     string `help:"Chroma style to use in light mode" default:"github" enum:"abap,algol,algol_nu,arduino,autumn,borland,bw,colorful,dracula,emacs,friendly,fruity,github,igor,lovelace,manni,monokai,monokailight,murphy,native,paraiso-dark,paraiso-light,pastie,perldoc,pygments,rainbow_dash,rrt,solarized-dark,solarized-dark256,solarized-light,swapoff,tango,trac,vim,vs,xcode"`
	CodeDarkStyle      string `help:"Chroma style to use in dark mode" default:"dracula" enum:"abap,algol,algol_nu,arduino,autumn,borland,bw,colorful,dracula,emacs,friendly,fruity,github,igor,lovelace,manni,monokai,monokailight,murphy,native,paraiso-dark,paraiso-light,pastie,perldoc,pygments,rainbow_dash,rrt,solarized-dark,solarized-dark256,solarized-light,swapoff,tango,trac,vim,vs,xcode"`

	Path string `arg name:"path" help:"Jupyter Notebook file to convert" type:"path"`
}

func (r *convertCmd) Run(ctx *context) error {
	// Make sure notebook file exists
	if _, err := os.Stat(r.Path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file %s does not exist", r.Path)
		}
		return fmt.Errorf("error while checking whether file %s exists", r.Path)
	}

	// TODO: Convert
	fmt.Println("Convert", r.Path)
	return nil
}

var cli struct {
	Debug bool `help:"Enable debug mode"`

	Convert convertCmd `cmd help:"Convert Jupyter Notebook file to HTML"`
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command
	err := ctx.Run(&context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}
