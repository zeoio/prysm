package main

import (
	"flag"
	"os"
	"strings"
	"text/template"
)

type data struct {
	Src       string
	SrcPkg    string
	Target    string
	TargetPkg string
	Out       string
	OutPkg    string
	Type      string
}

func main() {
	var d data
	flag.StringVar(&d.Src, "src", "", "Source package path")
	flag.StringVar(&d.Target, "target", "", "Target package path")
	flag.StringVar(&d.SrcPkg, "src-pkg", "", "Source package name")
	flag.StringVar(&d.TargetPkg, "target-pkg", "", "Target package name")
	flag.StringVar(&d.Out, "out", "", "Output file name")
	flag.StringVar(&d.Out, "out-pkg", "", "Output package name")
	flag.StringVar(&d.Type, "type", "", "The type to write migration functions for")
	flag.Parse()
	f, err := os.Create(d.Out)
	if err != nil {
		panic(err)
	}
	tpl, err := template.New("migration").Funcs(template.FuncMap{
		"capitalize": func(str string) string {
			return strings.Title(str)
		},
	}).Parse(queueTemplate)
	if err != nil {
		panic(err)
	}
	tpl.Execute(f, d)
}

var queueTemplate = `package {{.OutPkg}}

import (
	{{.SrcPkg}} "{{.Src}}"
	{{.TargetPkg}} "{{.Target}}"
)

func {{capitalize .SrcPkg}}To{{capitalize .TargetPkg}}{{.Type}}(src *{{.SrcPkg}}.{{.Type}}) *{{.TargetPkg}}.{{.Type}} {
	return &{{.TargetPkg}}.{{.Type}}{}
}

func {{capitalize .TargetPkg}}To{{capitalize .SrcPkg}}{{.Type}}(src *{{.TargetPkg}}.{{.Type}}) *{{.SrcPkg}}.{{.Type}} {
	return &{{.SrcPkg}}.{{.Type}}{}
}
`
