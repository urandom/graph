package graph

type Id int
type ConnectorType int

const (
	InputType ConnectorType = iota
	OutputType

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
	Target() (Node, Connector)
	Connect(target Node, connector Connector)
	Disconnect()
}

type Linker interface {
	Node() Node
	Connect(target Node, source, sink Connector)
	Disconnect(source Connector)
	Link(target Node)
	Unlink()
	Connector(name string, kind ...ConnectorType) Connector
	Connection(source ...Connector) (Node, Connector)
}

type Processor interface {
	Process()
}
