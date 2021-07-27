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
	Types           map[string]bool
}

type structTemplateData struct {
	TypName   string
	SrcPkg    string
	TargetPkg string
	Fields    []string
}

func main() {
	d := &data{}
	flag.StringVar(&d.Src, "src", "", "Source package path")
	flag.StringVar(&d.Target, "target", "", "Target package path")
	flag.StringVar(&d.SrcPkg, "src-pkg", "", "Source package name")
	flag.StringVar(&d.TargetPkg, "target-pkg", "", "Target package name")
	flag.StringVar(&d.TargetRelative, "target-relative", "", "Relative target package path")
	flag.StringVar(&d.Out, "out", "", "Output file name")
	flag.StringVar(&d.OutPkg, "out-pkg", "", "Output package name")
	flag.StringVar(&d.TypesListString, "types", "", "The type to write migration functions for")
	flag.Parse()

	typesList := strings.Split(d.TypesListString, ",")
	d.Types = make(map[string]bool)
	for _, typItem := range typesList {
		d.Types[typItem] = true
	}

	parseStructs(d)

	parseTransientStructs(d)

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

func parseStructs(d *data) {
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

func parseTransientStructs(d *data) {
	// Add any dependency structs also to the list to generate.
	for structName := range d.Types {
		item, ok := structs[structName]
		if !ok {
			panic(structName)
		}
		for _, field := range item.Fields.List {
			fieldType, ok := field.Type.(*ast.StarExpr)
			if !ok {
				continue
			}
			if isUnexportedField(field.Names[0].Name) {
				continue
			}
			if strings.Contains(field.Names[0].Name, "Time") {
				continue
			}
			d.Types[fmt.Sprintf("%s", fieldType.X)] = true
		}
	}
}

func migrateStruct(srcPkg, targetPkg, typName string) string {
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
	tpl, err := template.New("struct").Funcs(template.FuncMap{
		"capitalize": func(str string) string {
			return strings.Title(str)
		},
		"migrateStruct": migrateStruct,
		"handleField":   handleField,
	}).Parse(structTemplate)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBufferString("")
	tpl.Execute(buf, structTemplateData{
		TypName:   typName,
		SrcPkg:    srcPkg,
		TargetPkg: targetPkg,
		Fields:    fields,
	})
	return buf.String()
}

func handleField(srcPkg, targetPkg, fieldName string) string {
	fmt.Sprintln("hitting handle field", srcPkg, targetPkg, fieldName)
	_, isStruct := structs[fieldName]
	if isStruct {
		return fmt.Sprintf(
			"%sTo%s%s(src.%s)",
			capitalize(srcPkg),
			capitalize(targetPkg),
			fieldName,
			fieldName,
		)
	}
	return fmt.Sprintf("src.%s", fieldName)
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

func capitalize(str string) string {
	return strings.Title(str)
}

var structTemplate = `&{{.TargetPkg}}.{{.TypName}}{
	{{range .Fields}}
		{{.}}: {{handleField .SrcPkg .TargetPkg .}},{{end}}
	}`

var topLevelTemplate = `package {{.OutPkg}}

import (
	{{.SrcPkg}} "{{.Src}}"
	{{.TargetPkg}} "{{.Target}}"
)
{{ $data := . }}
{{range $item, $b := .Types}}
// {{capitalize $data.SrcPkg}}To{{capitalize $data.TargetPkg}}{{$item}} --
func {{capitalize $data.SrcPkg}}To{{capitalize $data.TargetPkg}}{{$item}}(src *{{$data.SrcPkg}}.{{$item}}) *{{$data.TargetPkg}}.{{$item}} {
	if src == nil {
		return &{{$data.TargetPkg}}.{{$item}}{}
	}
	return {{migrateStruct $data.SrcPkg $data.TargetPkg $item}}
}

// {{capitalize $data.TargetPkg}}To{{capitalize $data.SrcPkg}}{{$item}} --
func {{capitalize $data.TargetPkg}}To{{capitalize $data.SrcPkg}}{{$item}}(src *{{$data.TargetPkg}}.{{$item}}) *{{$data.SrcPkg}}.{{$item}} {
	if src == nil {
		return &{{$data.SrcPkg}}.{{$item}}{}
	}
	return {{migrateStruct $data.TargetPkg $data.SrcPkg $item}}
}
{{end}}
`
