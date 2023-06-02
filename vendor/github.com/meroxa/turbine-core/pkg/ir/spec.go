package ir

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/heimdalr/dag"
)

type (
	ConnectorDirection string
	Lang               string
)

const (
	GoLang     Lang = "golang"
	JavaScript Lang = "javascript"
	Python     Lang = "python"
	Ruby       Lang = "ruby"
	Java       Lang = "java"

	ConnectorSource      ConnectorDirection = "source"
	ConnectorDestination ConnectorDirection = "destination"

	SpecVersion_0_1_1 = "0.1.1"
	SpecVersion_0_2_0 = "0.2.0"

	LatestSpecVersion = SpecVersion_0_2_0
)

var specVersions = []string{
	SpecVersion_0_1_1,
	SpecVersion_0_2_0,
}

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
	PluginName string                 `json:"plugin_name"`
	Direction  ConnectorDirection     `json:"type"`
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

	if c.Direction != ConnectorSource {
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

	if c.Direction != ConnectorDestination {
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

	err := d.upgradeToLatestSpecVersion()
	if err != nil {
		return turbineDag, err
	}

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

// upgradeToLatestSpecVersion will ensure that simple topologies as defined in 0.1.1 are still compatible
// currently, only supported migration is from 0.1.1 to 0.2.0.
func (d *DeploymentSpec) upgradeToLatestSpecVersion() error {
	if d.getSpecVersion() == LatestSpecVersion {
		return nil
	}

	if d.getSpecVersion() == "" {
		return fmt.Errorf("cannot upgrade to the latest version. spec version is not specified")
	}

	if d.getSpecVersion() != SpecVersion_0_1_1 || LatestSpecVersion != SpecVersion_0_2_0 {
		return fmt.Errorf("unsupported upgrade from spec version %q to %q", d.getSpecVersion(), LatestSpecVersion)
	}

	if d.getSpecVersion() != SpecVersion_0_1_1 || LatestSpecVersion != SpecVersion_0_2_0 {
		return fmt.Errorf("unsupported upgrade from spec version %q to %q", d.getSpecVersion(), LatestSpecVersion)
	}

	var sources, destinations []ConnectorSpec

	// assign UUIDs to all connectors
	for i, c := range d.Connectors {
		if c.UUID == "" {
			c.UUID = uuid.New().String()
			d.Connectors[i].UUID = c.UUID
		}

		switch c.Direction {
		case ConnectorSource:
			sources = append(sources, c)
		case ConnectorDestination:
			destinations = append(destinations, c)
		}
	}

	// validate supported DAG in 0.1.1
	if d.getSpecVersion() == SpecVersion_0_1_1 {
		switch {
		case len(d.Functions) > 1:
			return fmt.Errorf("unsupported number of functions in spec version %q", SpecVersion_0_1_1)
		case len(sources) > 1:
			return fmt.Errorf("unsupported number of sources in spec version %q", SpecVersion_0_1_1)
		}
	}

	if len(sources) == 0 {
		return errors.New("not sources found in spec")
	}

	// assign UUIDs to all functions
	for i, fn := range d.Functions {
		if fn.UUID == "" {
			fn.UUID = uuid.New().String()
			d.Functions[i].UUID = fn.UUID
		}
	}

	// create streams
	// s -> n(d)
	if len(d.Functions) == 0 {
		for _, dest := range destinations {
			d.Streams = append(d.Streams, StreamSpec{
				UUID:     uuid.New().String(),
				FromUUID: sources[0].UUID,
				ToUUID:   dest.UUID,
			})
		}
	}

	// s -> f -> n(d)
	if len(d.Functions) == 1 {
		// s -> f
		d.Streams = append(d.Streams, StreamSpec{
			UUID:     uuid.New().String(),
			FromUUID: sources[0].UUID,
			ToUUID:   d.Functions[0].UUID,
		})

		// f -> n(d)
		for _, dest := range destinations {
			d.Streams = append(d.Streams, StreamSpec{
				UUID:     uuid.New().String(),
				FromUUID: d.Functions[0].UUID,
				ToUUID:   dest.UUID,
			})
		}
	}
	return nil
}

func (d *DeploymentSpec) init() {
	d.dagInitOnce.Do(func() {
		d.turbineDag = *dag.NewDAG()
	})
}
