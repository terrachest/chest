package listavailableversions

type AvailableVersions struct {
	Modules Modules `json:"modules"`
}

type Modules struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	Version string `json:"version"`
}
