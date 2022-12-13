package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"privateterraformregistry/internal/datamanager"
	"privateterraformregistry/internal/downloader"
	"privateterraformregistry/internal/modules"
	"privateterraformregistry/internal/uploader"

	"github.com/gorilla/mux"
)

var dataDir = os.Getenv("DATA_DIR")

const (
	maxUploadSize = 32 << 20
)

func getServiceDiscovery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	type serviceDiscovery struct {
		ModulePath string `json:"modules.v1"`
	}

	json.NewEncoder(w).Encode(serviceDiscovery{
		ModulePath: "/terraform/modules/v1/",
	})
}

func listAvailableVersions(ms modules.Modules) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		namespace := vars["namespace"]
		system := vars["system"]
		name := vars["name"]

		type Version struct {
			Version string `json:"version"`
		}
		var availableVersions []Version

		for _, m := range ms.Modules {
			if m.Namespace == namespace && m.System == system && m.Name == name {
				availableVersions = append(availableVersions, Version{Version: m.Version})
			}
		}

		p := struct {
			Modules []struct {
				Versions []Version `json:"versions"`
			} `json:"modules"`
		}{
			Modules: []struct {
				Versions []Version `json:"versions"`
			}{
				{
					Versions: availableVersions,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(p)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getDownloadPath(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var downloadPath = fmt.Sprintf("/modules/%s.%s.%s.%s.tar.gz", vars["namespace"], vars["name"], vars["system"], vars["version"])
	w.Header().Set("X-Terraform-Get", downloadPath)
	w.WriteHeader(204)
}

func uploadModule(ms *modules.Modules, dm datamanager.DataManager) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(405)
			_, err := w.Write([]byte("Method Not Allowed"))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		vars := mux.Vars(r)

		m := modules.Module{
			Namespace: vars["namespace"],
			Name:      vars["name"],
			System:    vars["system"],
			Version:   vars["version"],
		}
		var err = m.Validate()

		if err != nil {
			log.Fatal(err)
			return
		}

		u := uploader.New(&uploader.Config{
			MaxUploadSize: maxUploadSize,
			DataDir:       dataDir,
		})
		err = u.Upload(r, m)

		if err != nil {
			log.Fatal(err)
			return
		}

		ms.Add(m)
		err = dm.Save()

		if err != nil {
			log.Fatal(err)
		}
	}
}

func downloadModule(ms *modules.Modules) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		d := downloader.New(&downloader.Config{
			DataDir: dataDir,
		})

		m := modules.Module{
			Namespace: vars["Namespace"],
			Name:      vars["Name"],
			System:    vars["System"],
			Version:   vars["Version"],
		}

		if !ms.Exists(m) {
			w.WriteHeader(404)
			_, err := w.Write([]byte("File not found."))

			if err != nil {
				log.Fatal(err)
			}
			return
		}

		err := d.Download(w, r, m)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	if dataDir == "" {
		dataDir = "/.privateterraformregistry/data"
	}

	ms := modules.Modules{}
	var dm = datamanager.New(&datamanager.Config{
		DataDir: dataDir,
	}, &ms)

	var err = dm.Load()

	if err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)

	// Request Logging
	router.Use(func(nxt http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.RequestURI)
			nxt.ServeHTTP(w, r)
		})
	})

	router.HandleFunc("/.well-known/terraform.json", getServiceDiscovery)                                      // terraform protocol
	router.HandleFunc("/terraform/modules/v1/{namespace}/{name}/{system}/versions", listAvailableVersions(ms)) // terraform module protocol
	router.HandleFunc("/terraform/modules/v1/{namespace}/{name}/{system}/{version}/download", getDownloadPath) // terraform module protocol
	router.HandleFunc("/modules/{namespace}/{name}/{system}/{version}", uploadModule(&ms, dm))
	router.HandleFunc("/modules/{filename}", downloadModule)

	log.Print("Server Ready")
	log.Fatal(http.ListenAndServe(":8080", router))
}
