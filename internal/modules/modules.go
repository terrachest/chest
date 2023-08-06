package modules

import "github.com/terrachest/chest/internal/module"

type Modules struct {
	Modules []module.Module
}

func (modules *Modules) GetModules() []module.Module {
	return modules.Modules
}

func (modules *Modules) Add(m module.Module) {
	if !modules.Exists(m) {
		modules.Modules = append(modules.Modules, m)
	}
}

func (modules *Modules) Exists(m module.Module) bool {
	for _, x := range modules.Modules {
		if x.Namespace == m.Namespace && x.System == m.System && x.Name == m.Name && x.Version == m.Version {
			return true
		}
	}
	return false
}

func (modules *Modules) Validate() {
	for _, m := range modules.Modules {
		m.Validate()
	}
}
