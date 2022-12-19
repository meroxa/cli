package ir

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/heimdalr/dag"
)

type ConnectorType string
type Lang string

const (
	GoLang               Lang          = "golang"
	JavaScript           Lang          = "javascript"
	NodeJs               Lang          = "nodejs"
	Python               Lang          = "python"
	Python3              Lang          = "python3"
	Ruby                 Lang          = "ruby"
	ConnectorSource      ConnectorType = "source"
	ConnectorDestination ConnectorType = "destination"

	LatestSpecVersion = "0.2.0"
)

type DeploymentSpec struct {
	mu          sync.Mutex
	turbineDag  dag.DAG
	dagInitOnce sync.Once
	Secrets     map[string]string `json:"secrets,omitempty"`
	Connectors  []ConnectorSpec   `json:"connectors"`
	Functions   []FunctionSpec    `json:"functions,omitempty"`
	Streams     []StreamSpec      `json:"streams,omitempty"`
	Definition  DefinitionSpec    `json:"definition"`
}

type StreamSpec struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	FromUUID string `json:"from_uuid"`
	ToUUID   string `json:"to_uuid"`
}

type ConnectorSpec struct {
	UUID       string                 `json:"uuid"`
	Type       ConnectorType          `json:"type"`
	Resource   string                 `json:"resource"`
	Collection string                 `json:"collection"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

type FunctionSpec struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

type DefinitionSpec struct {
	GitSha   string       `json:"git_sha"`
	Metadata MetadataSpec `json:"metadata"`
}

type MetadataSpec struct {
	Turbine     TurbineSpec `json:"turbine"`
	SpecVersion string      `json:"spec_version"`
}

type TurbineSpec struct {
	Language Lang   `json:"language"`
	Version  string `json:"version"`
}

func ValidateSpecVersion(specVersion string) error {
	if specVersion != LatestSpecVersion {
		return fmt.Errorf("spec version %q is not a supported. use version %q instead", specVersion, LatestSpecVersion)
	}
	return nil
}

func (d *DeploymentSpec) SetImageForFunctions(image string) {
	for i := range d.Functions {
		d.Functions[i].Image = image
	}
}

func (d *DeploymentSpec) Marshal() ([]byte, error) {
	if err := d.ValidateDAG(); err != nil {
		return nil, err
	}
	return json.Marshal(d)
}

func Unmarshal(data []byte) (*DeploymentSpec, error) {
	spec := &DeploymentSpec{}
	if err := json.Unmarshal(data, spec); err != nil {
		return nil, err
	}
	return spec, nil
}

func (d *DeploymentSpec) AddSource(c *ConnectorSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.init()

	if len(d.turbineDag.GetRoots()) >= 1 {
		return fmt.Errorf("can only add one source connector per application")
	}
	if c.Type != ConnectorSource {
		return fmt.Errorf("not a source connector.")
	}
	d.Connectors = append(d.Connectors, *c)
	return d.turbineDag.AddVertexByID(c.UUID, &c)
}

func (d *DeploymentSpec) AddFunction(f *FunctionSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.init()

	d.Functions = append(d.Functions, *f)
	return d.turbineDag.AddVertexByID(f.UUID, &f)
}

func (d *DeploymentSpec) AddDestination(c *ConnectorSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.init()

	if c.Type != ConnectorDestination {
		return fmt.Errorf("not a destination connector.")
	}
	d.Connectors = append(d.Connectors, *c)
	return d.turbineDag.AddVertexByID(c.UUID, &c)
}

func (d *DeploymentSpec) AddStream(s *StreamSpec) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.init()

	if _, err := d.turbineDag.GetVertex(s.FromUUID); err != nil {
		return fmt.Errorf("source %s does not exist", s.FromUUID)
	}

	if _, err := d.turbineDag.GetVertex(s.ToUUID); err != nil {
		return fmt.Errorf("destination %s does not exist", s.ToUUID)
	}

	d.Streams = append(d.Streams, *s)
	return d.turbineDag.AddEdge(s.FromUUID, s.ToUUID)
}

func (d *DeploymentSpec) ValidateDAG() error {
	turbineDag, err := d.BuildDAG()
	if err != nil {
		return err
	}
	if turbineDag == nil {
		return fmt.Errorf("unable to build dag, no resources found")
	}
	if len(turbineDag.GetRoots()) > 1 {
		return fmt.Errorf("too many source connectors")
	}
	if len(turbineDag.GetRoots()) == 0 {
		return fmt.Errorf("no source connector found.")
	}
	return nil
}

func (d *DeploymentSpec) BuildDAG() (*dag.DAG, error) {
	turbineDag := dag.NewDAG()
	for i := range d.Connectors {
		con := &d.Connectors[i]
		if err := turbineDag.AddVertexByID(con.UUID, con); err != nil {
			return nil, err
		}
	}
	for i := range d.Functions {
		fun := &d.Functions[i]
		if err := turbineDag.AddVertexByID(fun.UUID, fun); err != nil {
			return nil, err
		}
	}
	for _, stream := range d.Streams {
		if err := turbineDag.AddEdge(stream.FromUUID, stream.ToUUID); err != nil {
			return nil, err
		}
	}
	return turbineDag, nil
}

func (d *DeploymentSpec) init() {
	d.dagInitOnce.Do(func() {
		d.turbineDag = *dag.NewDAG()
	})
}
