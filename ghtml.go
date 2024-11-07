// Copyright 2024 Enren Shen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package ghtml

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

var (
	engine *HTMLEngine
)

// Engine is the framework's instance.
// Create an instance of Engine, by using New().
type HTMLEngine struct {
	directory string
	extension string
	reload    bool
	left      string
	right     string
	layout    string
	funcMap   template.FuncMap
	template  *template.Template
	rmu       sync.RWMutex
}

// New returns a new blank Engine instance without any middleware attached.
// By default, the configuration is:
// - reload:    false
// - left:      {{
// - right:     }}
// - layout:    ""
// - funcMap:   nil
func New(directory, extension string) *HTMLEngine {
	engine = &HTMLEngine{
		directory: directory,
		extension: extension,
		reload:    false,
		left:      "{{",
		right:     "}}",
		layout:    "",
		funcMap:   make(template.FuncMap),
	}
	return engine
}

// Delim Setting Go syntax identifier
func (r *HTMLEngine) Delim(left, right string) *HTMLEngine {
	r.left, r.right = left, right
	return r
}

// Layout Set default layout after instantiation
func (r *HTMLEngine) Layout(name string) *HTMLEngine {
	r.layout = name
	return r
}

// Reload
// If Reload is set to true, it means that the disk file is reloaded every time and can be used for debugging environments.
// If Reload is set to home and only initialized once, it can be used in production environments.
func (r *HTMLEngine) Reload(reload bool) *HTMLEngine {
	r.reload = reload
	return r
}

// AddFunc Custom functions
func (r *HTMLEngine) AddFunc(name string, function any) *HTMLEngine {
	r.rmu.Lock()
	r.funcMap[name] = function
	r.rmu.Unlock()
	return r
}

// Instance Gin interface call
func (r *HTMLEngine) Instance(name string, data any) render.Render {
	return r.render(name, data, r.layout)
}

func (r *HTMLEngine) render(name string, data any, layout string) render.Render {

	if err := r.loadTemplate(); err != nil {
		_, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("Error in %s:%d%s\n", file, line, err)
		} else {
			log.Println(err)
		}

	}

	if layout != "" {
		r.loadFunc(name, data)
		name = layout
	}

	return render.HTML{
		Template: r.template,
		Name:     name,
		Data:     data,
	}
}

func (r *HTMLEngine) loadTemplate() error {
	if !r.reload && r.template != nil {
		return nil
	}

	r.rmu.Lock()
	defer r.rmu.Unlock()
	r.template = template.New(r.directory)
	r.template.Delims(r.left, r.right)

	err := filepath.Walk(r.directory, func(path string, info fs.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			rel, err := filepath.Rel(r.directory, path)
			if err != nil {
				return err
			}
			ext := filepath.Ext(path)
			if ext == r.extension {
				buf, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				standard := filepath.ToSlash(rel)
				tpl := r.template.New(standard)

				_, err = tpl.Funcs(r.emptyFunc()).Funcs(r.funcMap).Parse(string(buf))
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// Define empty functions to prevent loading errors
func (r *HTMLEngine) emptyFunc() template.FuncMap {
	return template.FuncMap{
		"content": func() string {
			return ""
		},
		"section": func() string {
			return ""
		},
		"render": func() string {
			return ""
		},
	}
}

func (r *HTMLEngine) loadFunc(name string, data any) {
	fc := template.FuncMap{
		"content": func() template.HTML {
			buf := new(bytes.Buffer)
			r.template.ExecuteTemplate(buf, name, data)
			return template.HTML(buf.String())
		},
		"section": func(sectionName string) template.HTML {
			replaceName := strings.Replace(name, r.extension, "", -1)
			fullSectionName := fmt.Sprintf("%s-%s", replaceName, sectionName)
			buf := new(bytes.Buffer)
			r.template.ExecuteTemplate(buf, fullSectionName, data)
			return template.HTML(buf.String())
		},
		"render": func(fullSectionName string) template.HTML {
			buf := new(bytes.Buffer)
			r.template.ExecuteTemplate(buf, name, data)
			return template.HTML(buf.String())
		},
	}

	for k, v := range r.funcMap {
		fc[k] = v
	}

	if tpl := r.template.Lookup(name); tpl != nil {
		tpl.Funcs(fc)
	}
}

// SetLayout  Setting layout
func SetLayout(ctx *gin.Context, name string) {
	ctx.Set("layout", name)
}

// Render Customizable response status code
func Render(ctx *gin.Context, code int, name string, data any) {
	layout, exists := ctx.Get("layout")
	if !exists {
		ctx.HTML(code, name, data)
	} else {
		ctx.Render(code, engine.render(name, data, layout.(string)))
	}
}

// View Default response status code 200
func View(ctx *gin.Context, name string, data any) {
	Render(ctx, 200, name, data)
}
