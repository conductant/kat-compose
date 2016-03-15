package aurora

import (
	"encoding/json"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/conductant/kat-compose/pkg/aurora/api"
	"github.com/conductant/kat-compose/pkg/encoding"
	. "gopkg.in/check.v1"
	"os"
	"testing"
)

func TestAurora(t *testing.T) { TestingT(t) }

type TestSuiteAurora struct {
}

var _ = Suite(&TestSuiteAurora{})

func (suite *TestSuiteAurora) SetUpSuite(c *C) {
}

func (suite *TestSuiteAurora) TearDownSuite(c *C) {
}

func schedulerHostPort() string {
	return os.Getenv("SCHEDULER_HOSTPORT")
}

func (suite *TestSuiteAurora) TestConnectAurora(c *C) {
	c.Assert(schedulerHostPort(), Not(Equals), "")

	trans, err := thrift.NewTHttpPostClient("http://" + schedulerHostPort() + "/api")
	c.Assert(err, IsNil)
	err = trans.Open()
	c.Assert(err, IsNil)

	protocolFactory := thrift.NewTJSONProtocolFactory()
	client := api.NewReadOnlySchedulerClientFactory(trans, protocolFactory)
	defer client.Transport.Close()

	c.Log(client)

	resp, err := client.GetRoleSummary()
	c.Assert(err, IsNil)
	for summary, _ := range resp.GetResult_().GetRoleSummaryResult_().Summaries {
		c.Log("Role=", summary.Role, ",Count=", summary.JobCount, ",CronJobCount=", summary.CronJobCount)

		resp, err := client.GetJobSummary(summary.Role)
		c.Assert(err, IsNil)
		for summary, _ := range resp.GetResult_().GetJobSummaryResult_().Summaries {

			// Fake some data
			constraint := &api.Constraint{
				Name: "const1",
				Constraint: &api.TaskConstraint{
					Limit: &api.LimitConstraint{
						Limit: 2,
					},
				},
			}
			summary.Job.TaskConfig.Constraints[constraint] = true

			m := encoding.MarshalMap(summary, map[string]encoding.OverrideFunc{
				".job.taskConfig.executorConfig.data": func(v interface{}) interface{} {
					s, ok := v.(string)
					if !ok {
						panic("not ok")
					}
					m := map[string]interface{}{}
					err := json.Unmarshal([]byte(s), &m)
					if err != nil {
						panic(err)
					}
					return m
				},
			})
			j, err := json.MarshalIndent(m, "  ", "  ")
			c.Assert(err, IsNil)
			c.Log(string(j))

			// going backwards
			summary2 := new(api.JobSummary)
			encoding.UnmarshalMap(m, summary2, map[string]encoding.OverrideFunc{
				".job.taskConfig.executorConfig.data": func(v interface{}) interface{} {
					m, ok := v.(map[string]interface{})
					if !ok {
						panic("not ok")
					}
					buff, err := json.Marshal(m)
					if err != nil {
						panic(err)
					}
					return string(buff)
				},
			})

			m2 := encoding.MarshalMap(summary, map[string]encoding.OverrideFunc{
				".job.taskConfig.executorConfig.data": func(v interface{}) interface{} {
					s, ok := v.(string)
					if !ok {
						panic("not ok")
					}
					m := map[string]interface{}{}
					err := json.Unmarshal([]byte(s), &m)
					if err != nil {
						panic(err)
					}
					return m
				},
			})
			j2, err := json.MarshalIndent(m2, "  ", "  ")
			c.Assert(err, IsNil)
			c.Log(string(j2))

			c.Assert(string(j), Equals, string(j2))
		}
	}
}
