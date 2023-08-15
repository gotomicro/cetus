package xpongo

import (
	"net/http"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin/render"
)

type (
	Pongo2 struct {
		TemplateSet *pongo2.TemplateSet
	}

	Render struct {
		Template *pongo2.Template
		Name     string
		Data     interface{}
	}

	Context pongo2.Context
)

var htmlContentType = []string{"text/html; charset=utf-8"}

func New(opts ...Option) *Pongo2 {
	var cfg config
	for _, opt := range opts {
		opt(&cfg)
	}
	p := &Pongo2{}
	var loader pongo2.TemplateLoader
	if cfg.fs != nil {
		loader = pongo2.NewFSLoader(cfg.fs)
	} else {
		loader = pongo2.MustNewLocalFileSystemLoader(cfg.path)
	}
	p.TemplateSet = pongo2.NewSet("ginpongo", loader)
	p.TemplateSet.Debug = cfg.debug
	for k, v := range cfg.globalContext {
		p.TemplateSet.Globals[k] = v
	}
	return p
}

func (p Pongo2) Instance(name string, data interface{}) render.Render {
	tpl := pongo2.Must(p.TemplateSet.FromCache(name))
	return &Render{
		Template: tpl,
		Name:     name,
		Data:     data,
	}
}

func (r Render) Render(w http.ResponseWriter) error {
	c := pongo2.Context(r.Data.(Context))
	return r.Template.ExecuteWriter(c, w)
}

func (r Render) WriteContentType(w http.ResponseWriter) {
	r.writeContentType(w, htmlContentType)
}

func (r Render) writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
