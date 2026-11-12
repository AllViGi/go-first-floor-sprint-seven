package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBadRequest(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	type Request struct {
		url  string
		body string
		code int
	}

	requests := []Request{
		{url: "/cafe", body: "unknown city", code: http.StatusBadRequest},
		{url: "/cafe?city=omsk", body: "unknown city", code: http.StatusBadRequest},
		{url: "/cafe?city=tula&count=na", body: "incorrect count", code: http.StatusBadRequest},
	}
	for _, req := range requests {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, req.url, nil)

		handler.ServeHTTP(response, request)

		assert.Equal(t, req.code, response.Code)
		assert.Equal(t, req.body, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	type Count struct {
		numberOfCafe int
		url          string
	}
	counts := []Count{
		{numberOfCafe: 0, url: "/cafe?city=moscow&count=0"},
		{numberOfCafe: 1, url: "/cafe?city=moscow&count=1"},
		{numberOfCafe: 2, url: "/cafe?city=moscow&count=2"},
		{numberOfCafe: 100, url: "/cafe?city=moscow&count=100"},
	}

	for _, count := range counts {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, count.url, nil)

		handler.ServeHTTP(response, request)

		require.Equal(t, http.StatusOK, response.Code)
		result := response.Body.String()
		var ln int
		if result != "" {
			ln = len(strings.Split(result, ","))
		}
		assert.Equal(t, min(count.numberOfCafe, len(cafeList["moscow"])), ln)
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	type request struct {
		search    string
		wantCount int
		url       string
	}
	requests := []request{
		{search: "фасоль", wantCount: 0, url: "/cafe?city=moscow&search=фасоль"},
		{search: "кофе", wantCount: 2, url: "/cafe?city=moscow&search=кофе"},
		{search: "вилка", wantCount: 1, url: "/cafe?city=moscow&search=вилка"},
	}
	for _, r := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, r.url, nil)

		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)
		result := response.Body.String()
		var ln int
		if result != "" {
			cafeSlc := strings.Split(result, ",")
			ln = len(cafeSlc)
			for _, cafe := range cafeSlc {
				assert.True(t, strings.Contains(strings.ToLower(cafe), strings.ToLower(r.search)))
			}

		}
		assert.Equal(t, r.wantCount, ln)

	}
}
