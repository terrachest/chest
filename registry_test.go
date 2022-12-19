package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"privateterraformregistry/internal/module"
	"privateterraformregistry/internal/moduleprotocol/listavailableversions"
	"privateterraformregistry/internal/moduleprotocol/servicediscovery"
	"privateterraformregistry/internal/modules"
	"testing"
)

func TestGetServiceDiscovery(t *testing.T) {
	r := gin.Default()
	r.GET("/.well-known/terraform.json", getServiceDiscovery)

	req, _ := http.NewRequest("GET", "/.well-known/terraform.json", nil)
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

func TestListAvailableServices(t *testing.T) {
	ms := modules.Modules{
		Modules: []module.Module{
			{"hashicorp", "consul", "aws", "0.1.0"},
			{"hashicorp", "consul", "aws", "1.0.0"},
			{"hashicorp", "consul", "aws", "1.1.0"},
		},
	}

	r := gin.Default()
	r.GET("/modules/v1/:namespace/:name/:system/versions", listAvailableVersions(&ms))

	req, _ := http.NewRequest("GET", "/modules/v1/hashicorp/consul/aws/versions", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Errorf("w.Code = %d; want 200", w.Code)
	}

	got := listavailableversions.AvailableVersions{}
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Errorf("")
	}

	want := listavailableversions.AvailableVersions{
		Modules: listavailableversions.Modules{
			Versions: []listavailableversions.Version{
				{"0.1.0"},
				{"1.0.0"},
				{"1.1.0"},
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("wrong result\n%s", diff)
	}
}

func TestGetDownloadPath(t *testing.T) {
	ms := modules.Modules{
		Modules: []module.Module{
			{"hashicorp", "consul", "gcp", "1.0.0"},
		},
	}

	r := gin.Default()
	r.GET("/modules/v1/:namespace/:name/:system/:version/download", getDownloadPath(&ms))

	req, _ := http.NewRequest("GET", "/modules/v1/hashicorp/consul/gcp/1.0.0/download", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	if w.Code != 204 {
		t.Errorf("w.Code = %d; want 200", w.Code)
	}

	want := "/modules/hashicorp/consul/gcp/1.0.0"
	if got := w.Header().Get("X-Terraform-Get"); got != want {
		t.Errorf("got != %s; got %s", want, got)
	}

	req, _ = http.NewRequest("GET", "/modules/v1/hashicorp/consul/aws/1.0.0/download", nil)
	w = httptest.NewRecorder()

	// Ensure if module does not exist 404 response is returned
	r.ServeHTTP(w, req)
	if w.Code != 404 {
		t.Errorf("w.Code = %d; want 404", w.Code)
	}
}

func TestDownloadModule(t *testing.T) {

}

func TestUploadModule(t *testing.T) {

}
