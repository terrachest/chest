package module

import "fmt"

type Module struct {
	Namespace string
	Name      string
	System    string
	Version   string
}

func (*Module) Validate() {
	// log fatal
}

func (module *Module) GetFileName() string {
	return fmt.Sprintf(
		"%s.%s.%s.%s.tar.gz",
		module.Namespace,
		module.Name,
		module.System,
		module.Version,
	)
}
