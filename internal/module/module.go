package module

type Module struct {
	Namespace string
	Name      string
	System    string
	Version   string
}

func (*Module) Validate() error {
	return nil
}
