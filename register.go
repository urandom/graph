package graph

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"sync"
	"text/template"
)

type LinkerJSONConstructor func(opts json.RawMessage) (Linker, error)
type JSONTemplateData struct {
	Args []string
}

var (
	operationsMu sync.Mutex
	operations   = make(map[string]LinkerJSONConstructor)
)

func RegisterLinker(name string, constructor LinkerJSONConstructor) {
	operationsMu.Lock()
	defer operationsMu.Unlock()

	if constructor == nil {
		panic("drawgl: Register operation constructor is nil")
	}

	if _, dup := operations[name]; dup {
		panic("drawgl: Register called twice for constructor " + name)
	}

	operations[name] = constructor
}

type jsonLinker struct {
	// The registered name of the linker
	Name string `json:"name,omitempty"`
	// A linker with this reference id to connect to
	ReferenceTo uint16 `json:"referenceTo,omitempty"`
	// A reference id, used to connect a separate branch to a linker
	ReferenceId uint16 `json:"referenceId,omitempty"`
	// The constructor options for this linker
	Options json.RawMessage `json:"options,omitempty"`
	// The input connector name. If empty, the default name is used
	Input ConnectorName `json:"input,omitempty"`
	// A map of all child linkers that are connected to the corresponding
	// output connector names
	Outputs map[ConnectorName]jsonLinker `json:"outputs,omitempty"`
}

type convertError struct {
	linker jsonLinker
	err    error
}

type deferredLinker struct {
	linker     Linker
	outputName ConnectorName
	inputName  ConnectorName
}

/*
A list of json objects. Each line represents a root

{
	"Name": "Load",
	"Options": {
		"Path": "{{ index .Args 0 }}"
	},
	"Outputs": {
		"Output": {
			"Name": "Convolution",
			"Options": {
				"Kernel": [-1, -1, -1, -1, 8, -1, -1, -1, -1],
				"Noralize": true
			},
			"Outputs": {
				"Output": {
					"Name": "Save",
					"Options": {
						"Path": {{ if gt (len .Args) 1 }} "{{ index .Args 1 }}" {{ else }} "/tmp/out.png" {{ end }}
					}
				}
			}
		}
	}
}
*/

func ProcessJSON(input interface{}, templateData *JSONTemplateData) (roots []Linker, err error) {
	defer func() {
		if r := recover(); r != nil {
			if ce, ok := r.(convertError); ok {
				roots = []Linker{}
				err = fmt.Errorf("processing json linker data for node '%#v': %v", ce.linker, ce.err)
			} else {
				panic(r)
			}
		}
	}()

	var jsonInput string
	var dec *json.Decoder

	switch t := input.(type) {
	case string:
		jsonInput = t
	case []byte:
		jsonInput = string(t)
	case io.Reader:
		var b []byte
		if b, err = ioutil.ReadAll(t); err != nil {
			err = fmt.Errorf("reading json from reader: %v", err)
			return
		}
		jsonInput = string(b)
	case *json.Decoder:
		dec = t
	default:
		panic(fmt.Sprintf("Unkown type: %T", input))
	}

	if dec == nil {
		var r io.Reader

		if templateData == nil {
			r = strings.NewReader(jsonInput)
		} else {
			var t *template.Template

			if t, err = template.New("json").Parse(jsonInput); err != nil {
				err = fmt.Errorf("parsing template: %v", err)
				return
			}

			var b bytes.Buffer
			if err = t.Execute(&b, templateData); err != nil {
				err = fmt.Errorf("executing template: %v", err)
				return
			}

			r = &b
		}
		dec = json.NewDecoder(r)
	}

	var references = make(map[uint16]Linker)
	var deferred = make(map[uint16][]deferredLinker)

	for {
		var root jsonLinker
		if err = dec.Decode(&root); err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				panic(convertError{linker: jsonLinker{}, err: fmt.Errorf("decoding root: %v", err)})
			}
		}

		r, rId := jsonToLinker(root, references)
		if rId > 0 {
			panic(convertError{linker: jsonLinker{}, err: errors.New("roots cannot be references")})
		}

		processLinkerTree(r, root, references, deferred)
		roots = append(roots, r)
	}

	return
}

func processLinkerTree(p Linker, rj jsonLinker, references map[uint16]Linker, deferred map[uint16][]deferredLinker) {
	for name, cj := range rj.Outputs {
		c, ref := jsonToLinker(cj, references)

		inputName := InputName
		if cj.Input != "" {
			inputName = cj.Input
		}

		if c != nil {
			if ops, ok := deferred[cj.ReferenceId]; ok {
				for _, op := range ops {
					opInputName := InputName
					if op.inputName != "" {
						opInputName = op.inputName
					}

					op.linker.Connect(c, op.linker.Connector(op.outputName, OutputType), c.Connector(opInputName, InputType))
				}

				delete(deferred, cj.ReferenceId)
			}
			p.Connect(c, p.Connector(name, OutputType), c.Connector(inputName, InputType))

			processLinkerTree(c, cj, references, deferred)
		} else if ref > 0 {
			deferred[ref] = append(deferred[ref], deferredLinker{linker: p, outputName: name, inputName: inputName})
		} else {
			panic(convertError{linker: cj, err: errors.New("no child linker or reference id")})
		}
	}
}

func jsonToLinker(j jsonLinker, references map[uint16]Linker) (Linker, uint16) {
	if j.Name != "" {
		c := operations[j.Name]
		if c == nil {
			panic(convertError{linker: j, err: fmt.Errorf("unknown name %s", j.Name)})
		}

		l, err := c(j.Options)
		if err != nil {
			panic(convertError{linker: j, err: fmt.Errorf("constructor failed for %s: %v", j.Name, err)})
		}

		if j.ReferenceId != 0 {
			references[j.ReferenceId] = l
		}

		return l, 0
	} else if j.ReferenceTo != 0 {
		if l, ok := references[j.ReferenceTo]; ok {
			return l, 0
		} else {
			return nil, j.ReferenceTo
		}
	} else {
		panic(convertError{linker: j, err: fmt.Errorf("json linker contains no name or link target")})
	}
}
