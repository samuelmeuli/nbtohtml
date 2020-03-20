module github.com/samuelmeuli/nbtohtml

go 1.14

require (
	github.com/alecthomas/chroma v0.7.1
	github.com/alecthomas/kong v0.2.4
	github.com/buildkite/terminal-to-html v3.2.0+incompatible
	github.com/dlclark/regexp2 v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/stretchr/testify v1.5.1
	golang.org/x/sys v0.0.0-20200317113312-5766fd39f98d // indirect
	gopkg.in/russross/blackfriday.v2 v2.0.1
)

replace gopkg.in/russross/blackfriday.v2 v2.0.1 => github.com/russross/blackfriday/v2 v2.0.1
