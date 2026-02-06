package webserver

import (
	"fmt"
	"io"

	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func newTpl(name string) *template.Template {
	tplFuncs := template.FuncMap{
		"dbtype": dbtype,
	}
	return template.New(name).Funcs(tplFuncs).Delims("<%{", "}%>")
}

// SetContentByTplFile 从模板文件中设置内容
//
// Example1:
//
//	f, err := os.OpenFile(targetFilepath, os.O_CREATE|os.O_WRONLY, 0o755)
//	SetContentByTplFile(tplFilepath, f, data)
//	f.Close()
//
// Example2:
//
//	var bf bytes.Buffer
//	var data = map[string]interface{}{"name": "Tom"}
//	SetContentByTplFile(tplFilepath, &bf, data)
//	SetContentByTplFile(tplFilepath, os.Stdout, data)
func SetContentByTplFile(tplFilepath string, wr io.Writer, data interface{}) error {
	// t, err := template.ParseFiles(tplFilepath)
	t, err := parseFiles(tplFilepath)
	if err != nil {
		return err
	}
	return t.Execute(wr, data)
}

func parseFiles(filenames ...string) (*template.Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("template: no files named in call to ParseFiles")
	}

	var t *template.Template
	for _, filename := range filenames {
		name, b, err := readFileOS(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		var tmpl *template.Template
		if t == nil {
			t = newTpl(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name) // .Funcs(tplFuncs)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func readFileOS(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}

func dbtype(t string) string {
	result := "VARCHAR(255)"
	switch strings.ToUpper(t) {
	case "STRING":
		result = "VARCHAR(255)"
	case "INT", "SMALLINT", "BIGINT", "FLOAT", "TEXT":
		result = strings.ToUpper(t)
	}
	return result
}
