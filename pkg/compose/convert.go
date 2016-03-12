package compose

import (
	"encoding/json"
	"fmt"
	"github.com/conductant/gohm/pkg/resource"
	"github.com/conductant/gohm/pkg/template"
	"github.com/conductant/kat-compose/pkg/aurora"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"io"
)

type Convert struct {
	Load
	Template string `flag:"t, Template file"`

	aurora.Project
}

func (this *Convert) Help(w io.Writer) {
	fmt.Fprintln(w, "Convert Compose files to some format")
}

func (this *Convert) Run(args []string, w io.Writer) error {
	err := this.Load.Run(args, w)
	if err != nil {
		return err
	}

	glog.Infoln("Template file=", this.Template)
	glog.Infoln("Parsed=", this.Parsed)

	tpl, err := this.loadTemplate()
	if err != nil {
		return err
	}

	switch this.Parsed["version"] {
	case "2", 2:
		return this.processV2(args, w, tpl)
	default: // v1
		return this.processV1(args, w, tpl)
	}
	return nil
}

func (this *Convert) Close() error {
	return nil
}

func (this *Convert) loadTemplate() ([]byte, error) {
	if this.Template == "" {
		return []byte(aurora.DefaultJobJSON), nil
	}
	ctx := context.Background()
	return resource.Fetch(ctx, this.Template)
}

type KeyValue struct {
	Key   string
	Value interface{}
}

func addList(p []KeyValue, flag string, composeValue interface{}) []KeyValue {
	out := p
	switch l := composeValue.(type) {
	case []string:
		for _, e := range l {
			out = append(out, KeyValue{flag, e})
		}
	case []interface{}:
		for _, e := range l {
			out = append(out, KeyValue{flag, fmt.Sprintf("%v", e)})
		}
	}
	return out
}

// Returns a map of all k-v pairs that can be set on the docker run command line:
func getDockerParams(config interface{}) []KeyValue {
	m, is := config.(map[string]interface{})
	if !is {
		return []KeyValue{}
	}
	p := []KeyValue{}
	for k, v := range m {
		switch k {
		case "cpu_share":
			p = append(p, KeyValue{"cpu-share", v})
		case "cpu_quota":
			p = append(p, KeyValue{"cpu-quota", v})
		case "environment":
			p = addList(p, "e", v)
		case "env_file":
			// Generate list of env
			p = addList(p, "e", v)
		case "expose":
			addList(p, "expose", v)
		case "image": // Specified somewhere else and not part of Docker Params
		case "links":
		case "networks", "net":
			p = addList(p, "net", v)
		case "ports":
			p = addList(p, "p", v)
		case "volumes":
			p = addList(p, "v", v)
		default:
			p = append(p, KeyValue{k, v})
		}
	}
	return p
}

func (this *Convert) processServices(services map[string]interface{}, tpl []byte) ([]json.RawMessage, error) {
	jobs := []json.RawMessage{}
	for service, config := range services {
		data := map[string]interface{}{
			"cluster":       this.Project.Cluster,
			"contact":       this.Project.Contact,
			"role":          this.Project.Role,
			"environment":   this.Project.Environment,
			"is_production": this.Project.IsProduction,
			"is_service":    true,
			"name":          service,
			"compose":       config,
			"docker_params": getDockerParams(config), // map some fields to docker flags
		}
		raw, err := template.Apply(tpl, data)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, raw)
	}
	return jobs, nil
}

func (this *Convert) processV1(args []string, w io.Writer, tpl []byte) error {
	glog.Infoln("Processing v1")
	jobs, err := this.processServices(this.Parsed, tpl)
	if err != nil {
		return err
	}

	buff, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(buff))
	return nil
}

func (this *Convert) processV2(args []string, w io.Writer, tpl []byte) error {
	glog.Infoln("Processing v2")
	jobs, err := this.processServices(this.Parsed["services"].(map[string]interface{}), tpl)
	if err != nil {
		return err
	}

	buff, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}
	fmt.Print(string(buff))
	return nil
}
