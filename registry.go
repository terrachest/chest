package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"privateterraformregistry/internal/datamanager"
	"privateterraformregistry/internal/module"
	"privateterraformregistry/internal/moduleprotocol"
	"privateterraformregistry/internal/modules"
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
	r.MaxMultipartMemory = maxUploadSize

	if err = r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
	}

	r.GET("/.well-known/terraform.json", getServiceDiscovery)
	r.GET("/modules/v1/:namespace/:name/:system/versions", listAvailableVersions(&ms))
	r.GET("/modules/v1/:namespace/:name/:system/:version/download", getDownloadPath)
	r.GET("/modules/:namespace/:name/:system/:version", downloadModule)
	r.POST("/modules/:namespace/:name/:system/:version", uploadModule(&ms, dm))

	log.Fatal(r.Run())
}

func getServiceDiscovery(c *gin.Context) {
	mp := moduleprotocol.New([]module.Module{})
	c.JSON(http.StatusOK, mp.ServiceDiscovery())
}

func listAvailableVersions(ms *modules.Modules) func(c *gin.Context) {
	return func(c *gin.Context) {
		namespace := c.Param("namespace")
		system := c.Param("system")
		name := c.Param("name")

		mp := moduleprotocol.New(ms.GetModules())
		c.JSON(http.StatusOK, mp.AvailableVersions(
			namespace,
			system,
			name,
		))
	}
}

func getDownloadPath(c *gin.Context) {
	m := module.Module{
		Namespace: c.Param("namespace"),
		Name:      c.Param("name"),
		System:    c.Param("system"),
		Version:   c.Param("version"),
	}
	m.Validate()

	var downloadPath = fmt.Sprintf("/modules/%s/%s/%s/%s", c.Param("namespace"), c.Param("name"), c.Param("system"), c.Param("version"))
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
		m.Validate()

		file, err := c.FormFile("File")
		if err != nil {
			log.Fatal(err)
		}

		if err = c.SaveUploadedFile(file, m.GetFileName()); err != nil {
			log.Fatal(err)
		}

		ms.Add(m)
		if err = dm.Persist(*ms); err != nil {
			log.Fatal(err)
		}
	}
}

func downloadModule(c *gin.Context) {
	m := module.Module{
		Namespace: c.Param("namespace"),
		Name:      c.Param("name"),
		System:    c.Param("system"),
		Version:   c.Param("version"),
	}
	m.Validate()

	c.FileAttachment(dataDir, m.GetFileName())
}
