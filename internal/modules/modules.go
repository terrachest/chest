package modules

import "privateterraformregistry/internal/module"

type Modules struct {
	modules []module.Module
}

func (modules *Modules) GetModules() []module.Module {
	return modules.modules
}

func (modules *Modules) Exists(m module.Module) bool {
	for _, x := range modules.modules {
		if x.Namespace == m.Namespace && x.System == m.System && x.Name == m.Name && x.Version == m.Version {
			return true
		}
	}
	return false
}

func (modules *Modules) Add(m module.Module) {
	if !modules.Exists(m) {
		modules.modules = append(modules.modules, m)
	}
}

func (modules *Modules) Validate() error {
	for _, m := range modules.modules {
		err := m.Validate()
		if err != nil {
			return err
		}
	}
	return nil
}
