package role

type Operator struct{}

func (Operator) Name() string {
	return "operator"
}

func (Operator) Chat(name string, msg string) string {
	return name + ": " + msg
}

func (Operator) Color(name string) string {
	return name
}

func (Operator) Inherits() Role {
	return nil
}
