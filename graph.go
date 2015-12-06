package graph

import "errors"

// Id is the type of a node's id
type Id uint64

// ConnectorType is the type of connection a particular connector holds in
// relation to its linker
type ConnectorType int
type ConnectorName string

const (
	// InputType represents the input type of a connector
	InputType ConnectorType = iota
	// OutputType represents the output type of a connector
	OutputType

	// InputName is the name of the default input connector
	InputName ConnectorName = "input"
	// OutputName is the name of the default output connector
	OutputName ConnectorName = "output"
)

var (
	ErrSameConnectorType = errors.New("Two connectors of the same type cannot be linked together")
	ErrInvalidConnector  = errors.New("The given connector is invalid")
)

// Node is a basic work unit within a graph
type Node interface {
	// Id returns the unique id of a node
	Id() Id
}

// Linker is a node container that defines its position within its siblings in
// a graph
type Linker interface {
	// Node returns the underlying node
	Node() Node
	// Connect connects to the target linker, via the current's source
	// connector and the target's sink connector. It returns an error if the
	// sink connector is of the same type as the source connector, or if any
	// one of them is nil
	Connect(target Linker, source, sink Connector) error
	// Disconnect breaks the connection to any linker that's currently
	// connected to the source connector
	Disconnect(source Connector)
	// Link is a helper method that connects the linker to the target via the
	// default output connector and the target's default input connector
	Link(target Linker)
	// Unlink disconnects any linker that's connected to default output
	// connector
	Unlink()
	// Connector returns the linker's connector of the given name and type. If
	// no type is provided, it returns the input connector for the given name
	Connector(name ConnectorName, kind ...ConnectorType) Connector
	// Connectors returns all connnectors of a given type. If no type is
	// provided, it returns the input connectors
	Connectors(kind ...ConnectorType) []Connector
	// Connection is a helper method that returns the given connector's target
	// linker and connector. If no connector is supplied, it uses the default
	// input connector
	Connection(source ...Connector) (Linker, Connector)
}

// Connector represents an input or output point via which a linker can connect
// to its siblings within a graph
type Connector interface {
	// Type returns the connector's type
	Type() ConnectorType
	// Name returns the connector's name
	Name() ConnectorName
	// Target returns the linker and connector that are connected to this one
	Target() (Linker, Connector)
	// Connect connects the target linker and connector to this one. It only
	// setups the link on its own end, the linker itself setups the reciprocal
	// connection. It will return an error if the target connector is of the
	// same type, or if the target connector is nil
	Connect(target Linker, connector Connector) error
	// Disconnect breaks a connection with any linker that's currently
	// connected. Similarly to the connect method, it only removes the link on
	// its own end
	Disconnect()
}
