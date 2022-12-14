package modules

type Module struct {
	Namespace string
	Name      string
	System    string
	Version   string
}

type Modules struct {
	Modules []Module `json:"modules"`
}

func (ms *Modules) Exists(rhs Module) bool {
	for _, lhs := range ms.Modules {
		if lhs.Namespace == rhs.Namespace && lhs.System == rhs.System && lhs.Name == rhs.Name && lhs.Version == rhs.Version {
			return true
		}
	}
	return false
}

func (ms *Modules) Add(m Module) {
	if !ms.Exists(m) {
		ms.Modules = append(ms.Modules, m)
	}
}

func (m *Module) Validate() error {
	return nil
}
