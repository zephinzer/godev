// +build ignore

package main

import (
	"io/ioutil"
	"os"
	"path"
	"text/template"
	"time"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	dockerfile, err := ioutil.ReadFile(path.Join(cwd, "/data/generate/Dockerfile"))
	if err != nil {
		panic(err)
	}
	makefile, err := ioutil.ReadFile(path.Join(cwd, "/data/generate/Makefile"))
	if err != nil {
		panic(err)
	}
	dotGitignore, err := ioutil.ReadFile(path.Join(cwd, "/data/generate/.gitignore"))
	if err != nil {
		panic(err)
	}
	dotDockerignore, err := ioutil.ReadFile(path.Join(cwd, "/data/generate/.dockerignore"))
	if err != nil {
		panic(err)
	}
	f, err := os.Create("./data.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	packageTemplate.Execute(f, struct {
		Timestamp       time.Time
		Dockerfile      string
		Makefile        string
		DotGitignore    string
		DotDockerignore string
	}{
		Timestamp:       time.Now(),
		Dockerfile:      string(dockerfile),
		Makefile:        string(makefile),
		DotDockerignore: string(dotDockerignore),
		DotGitignore:    string(dotGitignore),
	})
}

var packageTemplate = template.Must(template.New("test").Parse(`
// DO NOT EDIT - CREATED BY GO:GENERATE AT
//   {{.Timestamp}}
// FILE GENERATED USING ~/app/data/generate.go

package main

var DataDockerfile = ` + "`" + `
{{.Dockerfile}}
` + "`" + `
var DataMakefile = ` + "`" + `
{{.Makefile}}
` + "`" + `
var DataDotGitignore = ` + "`" + `
{{.DotGitignore}}
` + "`" + `
var DataDotDockerignore = ` + "`" + `
{{.DotDockerignore}}
` + "`" + `
`))
