package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"privateterraformregistry/internal/moduleprotocol/servicediscovery"
	"testing"
)

func TestGetServiceDiscovery(t *testing.T) {
	r := gin.Default()
	r.POST("/.well-known/terraform.json", getServiceDiscovery)

	req, _ := http.NewRequest("POST", "/.well-known/terraform.json", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("w.Code = %d; want 200", w.Code)
	}

	got := servicediscovery.ServiceDiscovery{}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Error(err)
	}

	want := servicediscovery.ServiceDiscovery{
		ModulePath: "/modules/v1/",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("wrong result\n%s", diff)
	}
}
