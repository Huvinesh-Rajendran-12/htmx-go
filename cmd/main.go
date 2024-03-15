package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
    templates *template.Template
}

type Contact struct {
    Name string
    Email string
}

func NewContact(name string, email string) *Contact {
    return &Contact{
        Name: name,
        Email: email,
    }
}

type Data struct {
    Contacts []*Contact
    Count int
}

func (data *Data) hasEmail(email string) bool {
    for _, contact := range data.Contacts {
        if contact.Email == email {
            return true
        }
    }
    return false
}

func newData() *Data {
    return &Data{
        Contacts: []*Contact{
            NewContact("suji", "sujithra"),
            NewContact("huvinesh", "huvinesh"),
            NewContact("thevitra", "thevitra"),
        },
        Count: 3,
    }
}

func (data *Data) addData(contact *Contact) {
    data.Contacts = append(data.Contacts, contact)
    data.Count++
}

type Count struct {
    Count int
}

type FormData struct {
    Values map[string]string
    Errors map[string]string
}

func newFormData() *FormData {
    return &FormData{
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

type Page struct {
    Data *Data
    FormData *FormData
}

func newPage() *Page{
    return &Page{
        Data: newData(),
        FormData: newFormData(),
    }
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
    return &Templates{
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}

func main() {
    e := echo.New()
    e.Use(middleware.Logger())
    page := newPage()
    e.Renderer = NewTemplate()
    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", page)
    })
    
    e.POST("/contacts", func(c echo.Context) error {
        name := c.FormValue("name")
        email := c.FormValue("email")
        if page.Data.hasEmail(email) {
            fmt.Print("email checking...")
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email
            formData.Errors["email"] = "Email already exists!"
            return c.Render(422, "form", formData)
        }
        contact := NewContact(name, email)
        page.Data.addData(contact)
        e.Logger.Info(page)
        c.Render(200, "form", newFormData())
        return c.Render(200, "oob-contact",contact)
    })
    e.Logger.Fatal(e.Start(":42069"))
}
