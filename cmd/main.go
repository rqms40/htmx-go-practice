package main

import (
	"html/template"
	"io"
	"sync/atomic"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

type Count struct {
	Count int64
}

func main() {
	e := echo.New()
	
	var count int64

	t := &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${method} ${uri} ${status}\n",
	}))

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", Count{Count: atomic.LoadInt64(&count)})
	})

	e.POST("/count", func(c echo.Context) error {
		newCount := atomic.AddInt64(&count, 1)
		return c.Render(200, "count", Count{Count: newCount})
	})

    e.Logger.Fatal(e.Start(":5000"))
}
