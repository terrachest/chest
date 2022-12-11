package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func serviceDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	type serviceDiscovery struct {
		ModulePath string `json:"modules.v1"`
	}

	json.NewEncoder(w).Encode(serviceDiscovery{
		ModulePath: "/terraform/modules/v1/",
	})
}

func listAvailableVersions(w http.ResponseWriter, r *http.Request) {
	type version struct {
		Version string `json:"version"`
	}

	type module struct {
		Versions []version `json:"versions"`
	}

	type availableVersions struct {
		Modules []module `json:"modules"`
	}

	json.NewEncoder(w).Encode(availableVersions{Modules: []module{
		{Versions: []version{
			{Version: "1.0.0"},
			{Version: "1.1.0"},
			{Version: "2.0.0"},
		}},
	}})
}

func getDownloadPath(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-Terraform-Get", "https://api.github.com/repos/hashicorp/terraform-aws-consul/tarball/v0.0.1//*?archive=tar.gz")
	w.WriteHeader(204) // No Content
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/.well-known/terraform.json", serviceDiscoveryHandler)
	router.HandleFunc("/v1/{namespace}/{name}/{system}/versions", listAvailableVersions)
	router.HandleFunc("/v1/{namespace}/{name}/{system}/{version}/download", getDownloadPath)
	log.Fatal(http.ListenAndServe(":8080", router))
}
