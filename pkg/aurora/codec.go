package aurora

import (
	"encoding/json"
	"github.com/conductant/kat-compose/pkg/aurora/api"
	"github.com/conductant/kat-compose/pkg/encoding"
)

type JSONObject map[string]interface{}

var (
	Codec = encoding.NewCodec().OverrideMarshal(api.ExecutorConfig{}, "Data",
		func(v interface{}) interface{} {
			s, ok := v.(string)
			if !ok {
				panic("Data not a string")
			}
			m := map[string]interface{}{}
			err := json.Unmarshal([]byte(s), &m)
			if err != nil {
				panic(err)
			}
			return m
		}).OverrideUnmarshal(api.ExecutorConfig{}, "Data",
		func(v interface{}) interface{} {
			m, ok := v.(map[string]interface{})
			if !ok {
				panic("Data not map/JSONObject")
			}
			buff, err := json.Marshal(m)
			if err != nil {
				panic(err)
			}
			return string(buff)
		})
)
