package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	urlValues := url.Values{}
	form := New(urlValues)

	var expectedType = "*forms.Form"
	actualType := reflect.TypeOf(form).String()

	if actualType != expectedType {
		t.Errorf("Expected %s, but got %s", expectedType, actualType)
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/some-url", nil)
	form := New(r.PostForm)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("Form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("Form shows invalid when all the required fields are passed")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/some-url", nil)
	form := New(r.PostForm)
	if form.Has("a") {
		t.Error("Form does not have the field but shows that it does")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	r.Form = postedData
	if !form.Has("a") {
		t.Error("Has the designated field, but shows that it doesn't")
	}
}

func TestForm_MinLength(t *testing.T) {
	values := url.Values{
		"name": {"Rasil"},
	}
	form := New(values)
	hasMinLength := form.MinLength("name", 3)
	if !hasMinLength {
		t.Error("Satisfied minimum length validation but still throws error")
	}

	values = url.Values{
		"name": {"R"},
	}
	form = New(values)
	hasMinLength = form.MinLength("name", 4)
	if hasMinLength {
		t.Error("Should throw min length error but passed")
	}
}

func TestForm_Valid(t *testing.T) {
	errs := errors{}
	if len(errs) != 0 {
		t.Error("Form doesn't have errors but still form is invalid")
	}

	errs = errors{}
	errs.Add("field", "Field is required")
	if len(errs) == 0 {
		t.Error("Form contains errors but is shown as valid")
	}
}

func TestForm_IsEmail(t *testing.T) {
	values := url.Values{
		"email": {"rasil@gmail.com"},
	}
	form := New(values)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("Email is valid but throws error")
	}

	values = url.Values{
		"email": {"rasil@gmail"},
	}
	form = New(values)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("Email is invalid but test passes")
	}
}
