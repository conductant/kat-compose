package aurora

import (
	"encoding/json"
	"fmt"
	"github.com/conductant/kat-compose/pkg/aurora/api"
	. "gopkg.in/check.v1"
	"os"
	"testing"
	"time"
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

func (suite *TestSuiteAurora) _TestAuroraReadonly(c *C) {
	c.Assert(schedulerHostPort(), Not(Equals), "")

	client, err := Connect(schedulerHostPort())
	c.Assert(err, IsNil)
	defer client.Close()

	resp, err := client.Readonly.GetRoleSummary()
	c.Assert(err, IsNil)
	for summary, _ := range resp.GetResult_().GetRoleSummaryResult_().Summaries {
		resp, err := client.Readonly.GetJobSummary(summary.Role)
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

			m := Codec.MarshalMap(summary)
			j, err := json.MarshalIndent(m, "  ", "  ")
			c.Assert(err, IsNil)

			// going backwards
			summary2 := new(api.JobSummary)
			Codec.UnmarshalMap(m, summary2)

			m2 := Codec.MarshalMap(summary)
			j2, err := json.MarshalIndent(m2, "  ", "  ")
			c.Assert(err, IsNil)
			c.Assert(string(j), Equals, string(j2))

			m3 := Codec.MarshalMap(summary.Job)
			j3, err := json.MarshalIndent(m3, "  ", "  ")
			c.Assert(err, IsNil)
			c.Log(string(j3))
		}
	}
}

func (suite *TestSuiteAurora) TestGetJobs(c *C) {
	client, err := Connect(schedulerHostPort())
	c.Assert(err, IsNil)
	defer client.Close()

	jobs, err := client.GetJobs("www-data")
	c.Assert(err, IsNil)
	c.Log(jobs)
	c.Assert(len(jobs) > 0, Equals, true)
}

func (suite *TestSuiteAurora) TestGetJobSummary(c *C) {
	client, err := Connect(schedulerHostPort())
	c.Assert(err, IsNil)
	defer client.Close()

	summaries, err := client.GetJobSummary("www-data")
	c.Assert(err, IsNil)
	c.Log(summaries)
	c.Assert(len(summaries) > 0, Equals, true)
}

func (suite *TestSuiteAurora) TestGetQuota(c *C) {
	client, err := Connect(schedulerHostPort())
	c.Assert(err, IsNil)
	defer client.Close()

	quota, err := client.GetQuota("www-data")
	c.Assert(err, IsNil)
	c.Log(quota)
}

// Disabled until Auth is figured out.
func (suite *TestSuiteAurora) DISALBE_TestCreateJob(c *C) {
	client, err := Connect(schedulerHostPort())
	c.Assert(err, IsNil)
	defer client.Close()

	jobs, err := client.GetJobs("www-data")
	c.Assert(err, IsNil)
	c.Assert(len(jobs) > 0, Equals, true)

	job := jobs[0]
	job.Key.Name = fmt.Sprintf("%s-%d", job.Key.Name, time.Now().Unix())
	c.Log("Creating: ", job.Key.Name)

	created, err := client.CreateJob(job)
	c.Assert(err, IsNil)
	c.Assert(created, Equals, true)
}

// Disabled until Aurora fixes a bug on normalizing treatment of null and empty lists.
func (suite *TestSuiteAurora) DISABLE_TestGetTasksStatus(c *C) {
	client, err := Connect(schedulerHostPort())
	c.Assert(err, IsNil)
	defer client.Close()

	tasks, err := client.GetJobTasks("www-data", "prod", "hello")
	c.Assert(err, IsNil)
	c.Log(tasks)
	c.Assert(len(tasks) > 0, Equals, true)
}
