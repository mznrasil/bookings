package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var tests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", http.MethodGet, []postData{}, http.StatusOK},
	{"about", "/about", http.MethodGet, []postData{}, http.StatusOK},
	{"contact", "/contact", http.MethodGet, []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", http.MethodGet, []postData{}, http.StatusOK},
	{"majors-suite", "/majors-suite", http.MethodGet, []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", http.MethodGet, []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary", http.MethodGet, []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", http.MethodGet, []postData{}, http.StatusOK},
	{"post-search-availability", "/search-availability", http.MethodPost, []postData{
		{"start", "2000-08-10"},
		{"end", "2000-08-20"},
	}, http.StatusOK},
	{"post-search-availability-json", "/search-availability-json", http.MethodPost, []postData{
		{"start", "2000-08-10"},
		{"end", "2000-08-20"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", http.MethodPost, []postData{
		{"first_name", "Rasil"},
		{"last_name", "Maharjan"},
		{"email", "mzn.rasil@gmail.com"},
		{"phone", "9841224466"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	mux := getRoutes()
	server := httptest.NewServer(mux)
	defer server.Close()

	for _, test := range tests {
		if test.method == http.MethodGet {
			resp, err := server.Client().Get(server.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("For %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		} else if test.method == http.MethodPost {
			values := url.Values{}
			for _, param := range test.params {
				values.Add(param.key, param.value)
			}

			resp, err := server.Client().PostForm(server.URL+test.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatusCode {
				t.Errorf("For %s, expected %d but got %d", test.name, test.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
