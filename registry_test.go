package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServiceDiscovery(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/.well-known/terraform.json", nil)
	w := httptest.NewRecorder()

	serviceDiscoveryHandler(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	if string(data) != `{"modules.v1":"/terraform/modules/v1/"}
` {
		t.Errorf(`expected {"modules.v1":"/terraform/modules/v1/"} got %v`, string(data))
	}
	header := res.Header.Get(`Content-Type`)
	if string(header) != "application/json" {
		t.Errorf(`expected Content-Type: application/json got Content-Type: %v`, header)
	}
}

func TestListAvailableVersions(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/modules/hashicorp/consul/aws/versions", nil)
	w := httptest.NewRecorder()

	listAvailableVersions(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	if string(data) != `{"modules":[{"versions":[{"version":"1.0.0"},{"version":"1.1.0"},{"version":"2.0.0"}]}]}
` {
		t.Errorf(`expected {"modules":[{"versions":[{"version":"1.0.0"},{"version":"1.1.0"},{"version":"2.0.0"}]}]} got %v`, string(data))
	}
	header := res.Header.Get(`Content-Type`)
	if string(header) != "application/json" {
		t.Errorf(`expected Content-Type: application/json got Content-Type: %v`, header)
	}
}

func TestGetDownloadPath(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/modules/hashicorp/consul/aws/0.0.1/download", nil)
	w := httptest.NewRecorder()

	getDownloadPath(w, req)
	res := w.Result()
	defer res.Body.Close()
	header := res.Header.Get(`https://api.github.com/repos/hashicorp/terraform-aws-consul/tarball/v0.0.1//*?archive=tar.gz`)
	if string(header) != "" {
		t.Errorf(`expected https://api.github.com/repos/hashicorp/terraform-aws-consul/tarball/v0.0.1//*?archive=tar.gz got %v`, header)
	}
	if res.StatusCode != 204 {
		t.Errorf(`expected status code 204 got %v`, res.StatusCode)
	}
}
