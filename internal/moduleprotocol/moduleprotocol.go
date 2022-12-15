package moduleprotocol

import (
	"privateterraformregistry/internal/module"
	"privateterraformregistry/internal/moduleprotocol/listavailableversions"
)

type ModuleProtocol interface {
	AvailableVersions(namespace string, system string, name string) listavailableversions.AvailableVersions
}

func New(modules []module.Module) ModuleProtocol {
	return &moduleProtocol{
		modules: modules,
	}
}

type moduleProtocol struct {
	modules []module.Module
}

func (protocol *moduleProtocol) AvailableVersions(
	namespace string,
	system string,
	name string,
) listavailableversions.AvailableVersions {
	var availableVersions []listavailableversions.Version

	for _, m := range protocol.modules {
		if m.Namespace == namespace && m.System == system && m.Name == name {
			availableVersions = append(availableVersions, listavailableversions.Version{Version: m.Version})
		}
	}

	return listavailableversions.AvailableVersions{
		Modules: listavailableversions.Modules{
			Versions: availableVersions,
		},
	}
}
