package render

import (
	"net/http"
	"testing"

	"github.com/mznrasil/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	td := models.TemplateData{}

	req, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(req.Context(), "flash", "123")

	result := AddDefaultData(&td, req)

	if result.Flash != "123" {
		t.Error("Flash value of 123 not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Errorf("Error creating teamplate cache: %v", err)
	}
	appConfig.TemplateCache = tc
	appConfig.UseCache = false

	req, err := getSession()
	if err != nil {
		t.Error("Cannot get the session", err)
	}
	var ww myWriter
	err = Template(&ww, req, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("Error writing template to browser")
	}

	err = Template(&ww, req, "non-existent.page.tmpl", &models.TemplateData{})
	if err == nil {
		t.Error("Rendered template that does not exist")
	}
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(appConfig)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error("Error creatint template cache")
	}
}

func getSession() (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		return nil, err
	}

	ctx := req.Context()
	ctx, _ = session.Load(ctx, req.Header.Get("X-Session"))
	req = req.WithContext(ctx)

	return req, nil
}
