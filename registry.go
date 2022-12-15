package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"privateterraformregistry/internal/datamanager"
	"privateterraformregistry/internal/downloader"
	"privateterraformregistry/internal/module"
	"privateterraformregistry/internal/moduleprotocol"
	"privateterraformregistry/internal/modules"
	"privateterraformregistry/internal/uploader"
	"privateterraformregistry/internal/utils/env"
)

var dataDir = env.Get("DATA_DIR", "/.privateterraformregistry/data")

const (
	maxUploadSize = 32 << 20
)

func main() {
	ms := modules.Modules{}
	dm := datamanager.New(&datamanager.Config{DataDir: dataDir})
	err := dm.Load(&ms)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.GET("/.well-known/terraform.json", getServiceDiscovery)
	r.GET("/terraform/modules/v1/:namespace/:name/:system/versions", listAvailableVersions(ms))
	r.GET("/terraform/modules/v1/:namespace/:name/:system/:version/download", getDownloadPath)
	r.GET("/download/:filename", downloadModule)
	r.GET("/modules/:namespace/:name/:system/:version", uploadModule(&ms, dm))

	log.Fatal(r.Run())
}

func getServiceDiscovery(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	type serviceDiscovery struct {
		ModulePath string `json:"modules.v1"`
	}

	err := json.NewEncoder(c.Writer).Encode(serviceDiscovery{
		ModulePath: "/terraform/modules/v1/",
	})

	if err != nil {
		log.Fatal(err)
	}
}

func listAvailableVersions(ms modules.Modules) func(c *gin.Context) {
	return func(c *gin.Context) {
		moduleProtocol := moduleprotocol.New(ms.GetModules())
		namespace := c.Param("namespace")
		system := c.Param("system")
		name := c.Param("name")

		c.JSON(http.StatusOK, moduleProtocol.AvailableVersions(
			namespace,
			system,
			name,
		))
	}
}

func getDownloadPath(c *gin.Context) {
	var downloadPath = fmt.Sprintf("/download/%s.%s.%s.%s.tar.gz", c.Param("namespace"), c.Param("name"), c.Param("system"), c.Param("version"))
	c.Header("X-Terraform-Get", downloadPath)
	c.Status(204)
}

func uploadModule(ms *modules.Modules, dm datamanager.DataManager) func(c *gin.Context) {
	return func(c *gin.Context) {
		m := module.Module{
			Namespace: c.Param("namespace"),
			Name:      c.Param("name"),
			System:    c.Param("system"),
			Version:   c.Param("version"),
		}
		var err = m.Validate()

		if err != nil {
			log.Println(err)
			return
		}

		u := uploader.New(&uploader.Config{
			MaxUploadSize: maxUploadSize,
			DataDir:       dataDir,
		})

		err = u.Upload(c.Request, m)

		if err != nil {
			log.Fatal(err)
		}

		ms.Add(m)
		err = dm.Persist(*ms)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func downloadModule(c *gin.Context) {
	d := downloader.New(&downloader.Config{
		DataDir: dataDir,
	})

	fn := c.Param("filename")

	err := d.Download(c.Writer, c.Request, fn)

	if err != nil {
		log.Fatal(err)
	}
}
