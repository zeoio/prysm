package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/prysmaticlabs/prysm/shared/fileutil"
)

var structs = map[string]*ast.StructType{}

type data struct {
	Src             string
	SrcPkg          string
	Target          string
	TargetPkg       string
	TargetRelative  string
	Out             string
	OutPkg          string
	TypesListString string
	Types           []string
}

type structTemplateData struct {
	TypName   string
	TargetPkg string
	Fields    []string
}

func main() {
	var d data
	flag.StringVar(&d.Src, "src", "", "Source package path")
	flag.StringVar(&d.Target, "target", "", "Target package path")
	flag.StringVar(&d.SrcPkg, "src-pkg", "", "Source package name")
	flag.StringVar(&d.TargetPkg, "target-pkg", "", "Target package name")
	flag.StringVar(&d.TargetRelative, "target-relative", "", "Relative target package path")
	flag.StringVar(&d.Out, "out", "", "Output file name")
	flag.StringVar(&d.OutPkg, "out-pkg", "", "Output package name")
	flag.StringVar(&d.TypesListString, "types", "", "The type to write migration functions for")
	flag.Parse()

	d.Types = strings.Split(d.TypesListString, ",")

	parseStructs(d)

	f, err := os.Create(d.Out)
	if err != nil {
		panic(err)
	}
	tpl, err := template.New("migration").Funcs(template.FuncMap{
		"capitalize": func(str string) string {
			return strings.Title(str)
		},
		"migrateStruct": migrateStruct,
	}).Parse(topLevelTemplate)
	if err != nil {
		panic(err)
	}
	tpl.Execute(f, d)
}

func parseStructs(d data) {
	fset := token.NewFileSet()
	pkgPath, err := fileutil.ExpandPath(d.TargetRelative)
	if err != nil {
		panic(err)
	}
	packages, err := parser.ParseDir(fset, pkgPath, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	for _, pkg := range packages {
		for _, f := range pkg.Files {
			for _, decl := range f.Decls {
				fn, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}
				// Needs to be a type declaration.
				if fn.Tok.String() != "type" {
					continue
				}
				// Needs to be a type specification.
				sp, ok := fn.Specs[0].(*ast.TypeSpec)
				if !ok {
					continue
				}
				// Needs to be a struct type.
				structTyp, ok := sp.Type.(*ast.StructType)
				if !ok {
					continue
				}
				structs[sp.Name.String()] = structTyp
			}
		}
	}
}

func migrateStruct(targetPkg string, typName string) string {
	structObj, ok := structs[typName]
	if !ok {
		panic(fmt.Sprintf("Struct with name %s not found", typName))
	}
	fields := make([]string, 0)
	for _, field := range structObj.Fields.List {
		name := field.Names[0].Name
		if isUnexportedField(name) {
			continue
		}
		fields = append(fields, name)
	}
	fmt.Println(fields)
	tpl, err := template.New("struct").Funcs(template.FuncMap{
		"capitalize": func(str string) string {
			return strings.Title(str)
		},
		"migrateStruct": migrateStruct,
	}).Parse(structTemplate)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBufferString("")
	tpl.Execute(buf, structTemplateData{
		TypName:   typName,
		TargetPkg: targetPkg,
		Fields:    fields,
	})
	return buf.String()
}

func isUnexportedField(str string) bool {
	return unicode.IsLower(firstRune(str))
}

func firstRune(str string) (r rune) {
	for _, r = range str {
		return
	}
	return
}

var structTemplate = `&{{.TargetPkg}}.{{.TypName}}{
	{{range .Fields}}
	{{.}}: src.{{.}},{{end}}
}`

var topLevelTemplate = `package {{.OutPkg}}

import (
	{{.SrcPkg}} "{{.Src}}"
	{{.TargetPkg}} "{{.Target}}"
)
{{ $data := . }}
{{range .Types}}
// {{capitalize $data.SrcPkg}}To{{capitalize $data.TargetPkg}}{{.}} --
func {{capitalize $data.SrcPkg}}To{{capitalize $data.TargetPkg}}{{.}}(src *{{$data.SrcPkg}}.{{.}}) *{{$data.TargetPkg}}.{{.}} {
	if src == nil {
		return &{{$data.TargetPkg}}.{{.}}{}
	}
	return {{migrateStruct $data.TargetPkg .}}
}

// {{capitalize $data.TargetPkg}}To{{capitalize $data.SrcPkg}}{{.}} --
func {{capitalize $data.TargetPkg}}To{{capitalize $data.SrcPkg}}{{.}}(src *{{$data.TargetPkg}}.{{.}}) *{{$data.SrcPkg}}.{{.}} {
	if src == nil {
		return &{{$data.SrcPkg}}.{{.}}{}
	}
	return {{migrateStruct $data.SrcPkg .}}
}
{{end}}
`
