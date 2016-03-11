package compose

import (
	"encoding/json"
	"fmt"
	"github.com/conductant/gohm/pkg/encoding"
	"github.com/golang/glog"
	"io"
	"os"
)

type Load struct {
	ComposeFiles []string `flag:"f,          The Compose file"`
	ProjectName  string   `flag:"p,          The project name"`
	Dump         bool     `flag:"dump,       Dump the project structure as JSON"`

	Parsed map[string]interface{}
}

func (this *Load) Help(w io.Writer) {
	fmt.Fprintln(w, "Load and parse Compose files")
}

func (this *Load) Run(args []string, w io.Writer) error {
	err := this.loadByMap(args, w)
	if err != nil {
		return err
	}

	if this.Dump {
		glog.Infoln("Parsed:", this.Parsed)
		buff, err := json.MarshalIndent(this.Parsed, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(buff))
	}
	return nil
}

// converts to the map type that JSON encoder can work with
func convert(in map[interface{}]interface{}) map[string]interface{} {
	out := map[string]interface{}{}
	for k, v := range in {
		key := fmt.Sprintf("%v", k)
		switch v := v.(type) {
		case map[interface{}]interface{}:
			out[key] = convert(v)
		default:
			out[key] = v
		}
	}
	return out
}

func (this *Load) loadByMap(args []string, w io.Writer) error {
	view := map[interface{}]interface{}{}
	for _, f := range this.ComposeFiles {
		fd, err := os.Open(f)
		if err != nil {
			return err
		}
		err = encoding.UnmarshalYAML(fd, view)
		if err != nil {
			return err
		}
	}
	this.Parsed = convert(view)
	return nil
}

func (this *Load) Close() error {
	return nil
}
