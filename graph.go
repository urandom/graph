package graph

type Id int
type ConnectorType int

const (
	Input ConnectorType = iota
	Output

	InputName  = "input"
	OutputName = "output"
)

type Node interface {
	Id() Id
}

type Graph interface {
	Root() Node
}

type Connector interface {
	Type() ConnectorType
	Name() string
	Target() Node
}

type Linker interface {
	Connect(target Node, source, sink Connector)
	Link(target Node)
	Connection(source ...Connector) Node
}

type Processor interface {
	Process()
}
