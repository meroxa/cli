package ir

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/heimdalr/dag"
)

type (
	DirectionType string
	Lang          string
)

const (
	GoLang     Lang = "golang"
	JavaScript Lang = "javascript"
	Python     Lang = "python"
	Ruby       Lang = "ruby"

	PluginSource      DirectionType = "source"
	PluginDestination DirectionType = "destination"

	SpecVersion_v3    = "v3"
	LatestSpecVersion = SpecVersion_v3
)

var specVersions = []string{
	SpecVersion_v3,
}

type DeploymentSpec struct {
	mu          sync.Mutex
	turbineDag  dag.DAG
	dagInitOnce sync.Once
	Connectors  []ConnectorSpec `json:"connectors"`
	Functions   []FunctionSpec  `json:"functions,omitempty"`
	Streams     []StreamSpec    `json:"streams,omitempty"`
	Definition  DefinitionSpec  `json:"definition"`
}

type StreamSpec struct {
	UUID     string `json:"uuid"`
	Name     string `json:"name"`
	FromUUID string `json:"from_uuid"`
	ToUUID   string `json:"to_uuid"`
}

type ConnectorSpec struct {
	UUID         string            `json:"uuid"`
	Name         string            `json:"name"`
	PluginType   DirectionType     `json:"plugin_type"`
	PluginName   string            `json:"plugin_name"`
	PluginConfig map[string]string `json:"plugin_config,omitempty"`
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

func ValidateSpecVersion(ver string) error {
	for _, v := range specVersions {
		if v == ver {
			return nil
		}
	}

	return fmt.Errorf(
		"spec version %q is invalid, supported versions: %s",
		ver,
		strings.Join(specVersions, ", "),
	)
}

func (d *DeploymentSpec) SetImageForFunctions(image string) error {
	switch {
	case image == "" && len(d.Functions) > 0:
		return fmt.Errorf("empty image for functions")
	case image != "" && len(d.Functions) == 0:
		return fmt.Errorf("cannot set image without defined functions")
	}

	for i := range d.Functions {
		d.Functions[i].Image = image
	}
	return nil
}

func (d *DeploymentSpec) Marshal() ([]byte, error) {
	if _, err := d.BuildDAG(); err != nil {
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

	if c.PluginType != PluginSource {
		return fmt.Errorf("not a source connector")
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

	if c.PluginType != PluginDestination {
		return fmt.Errorf("not a destination connector")
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

func (d *DeploymentSpec) ValidateDAG(turbineDag *dag.DAG) error {
	if turbineDag == nil {
		return fmt.Errorf("invalid DAG, no resources found")
	}

	if len(turbineDag.GetRoots()) > 1 {
		return fmt.Errorf("invalid DAG, too many sources")
	}

	if len(turbineDag.GetRoots()) == 0 {
		return fmt.Errorf("invalid DAG, no sources found")
	}

	// No edges
	if turbineDag.GetSize() == 0 {
		return fmt.Errorf("invalid DAG, there has to be at least one source, at most one function, and zero or more destinations")
	}

	return nil
}

func (d *DeploymentSpec) getSpecVersion() string {
	return d.Definition.Metadata.SpecVersion
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

	return turbineDag, ValidateSpecVersion(d.Definition.Metadata.SpecVersion)
}

func (d *DeploymentSpec) init() {
	d.dagInitOnce.Do(func() {
		d.turbineDag = *dag.NewDAG()
	})
}
