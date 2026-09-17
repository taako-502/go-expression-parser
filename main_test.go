package main

import (
	"encoding/json"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHTTP(t *testing.T) {
	h := handler()
	for _, tc := range []struct {
		expression string
		status     int
		result     float64
	}{
		{"1+2*3", 200, 7}, {"0", 200, 0}, {"(1+2)*3", 200, 9}, {"1/0", 400, 0}, {"1+", 400, 0},
	} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest("GET", "/api/evaluate?expression="+url.QueryEscape(tc.expression), nil))
		if rec.Code != tc.status {
			t.Fatalf("%s: status %d", tc.expression, rec.Code)
		}
		var out response
		if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
			t.Fatal(err)
		}
		if tc.status == 200 && (out.Result == nil || *out.Result != tc.result || out.AST == nil) {
			t.Fatalf("unexpected response: %+v", out)
		}
		if tc.status == 400 && (out.Error == "" || out.Result != nil) {
			t.Fatalf("missing error: %+v", out)
		}
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "<html lang=\"ja\">") {
		t.Fatal("HTML not served")
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/missing", nil))
	if rec.Code != 404 {
		t.Fatalf("unknown route: %d", rec.Code)
	}
}
