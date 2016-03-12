package aurora

import (
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
}
