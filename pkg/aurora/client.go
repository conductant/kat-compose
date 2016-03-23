package aurora

import (
	"errors"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/conductant/kat-compose/pkg/aurora/api"
	"github.com/golang/glog"
)

var (
	ErrNotRun = errors.New("err-not-run")
)

type client struct {
	transport        thrift.TTransport
	Readonly         *api.ReadOnlySchedulerClient
	SchedulerManager *api.AuroraSchedulerManagerClient
	Admin            *api.AuroraAdminClient
}

func (this *client) Close() error {
	return this.transport.Close()
}

// Connects to the server and returns a client reference.  hostport is in the form of <host>:<port>.
func Connect(hostport string) (*client, error) {
	trans, err := thrift.NewTHttpPostClient("http://" + schedulerHostPort() + "/api")
	if err != nil {
		return nil, err
	}
	err = trans.Open()
	if err != nil {
		return nil, err
	}

	protocolFactory := thrift.NewTJSONProtocolFactory()

	return &client{
		transport:        trans,
		Readonly:         api.NewReadOnlySchedulerClientFactory(trans, protocolFactory),
		SchedulerManager: api.NewAuroraSchedulerManagerClientFactory(trans, protocolFactory),
		Admin:            api.NewAuroraAdminClientFactory(trans, protocolFactory),
	}, nil
}

func (this *client) AcquireLock(role, environment, name string) (*api.Lock, error) {
	jobKey := &api.JobKey{Role: role, Environment: environment, Name: name}
	lockKey := &api.LockKey{jobKey}
	resp, err := this.SchedulerManager.AcquireLock(lockKey)
	glog.Infoln("AcquireLock: resp=", resp, "err=", err)
	if err != nil {
		return nil, err
	}
	if resp.ResponseCode == api.ResponseCode_OK {
		return resp.GetResult_().GetAcquireLockResult_().Lock, nil
	} else {
		return nil, fmt.Errorf("no-lock:role=%s,environment=%s,name=%s", role, environment, name)
	}
}

func (this *client) ReleaseLock(lock *api.Lock, force bool) (bool, error) {
	validation := api.LockValidation_CHECKED
	if force {
		validation = api.LockValidation_UNCHECKED
	}
	resp, err := this.SchedulerManager.ReleaseLock(lock, validation)
	if err != nil {
		return false, err
	}
	return resp.ResponseCode == api.ResponseCode_OK, nil
}

type lockContext struct {
	Locking bool
	lock    *api.Lock
	client  *client
	Err     error
}

func (this *lockContext) Run(f func(lock *api.Lock) error) error {
	if !this.Locking || (this.lock != nil && this.Err == nil) {
		defer func() {
			recover() // in case user code panics
			if this.Locking {
				this.client.ReleaseLock(this.lock, true)
			}
		}()
		this.Err = f(this.lock)
		return this.Err
	} else {
		this.Err = ErrNotRun
		return this.Err
	}
}

func (this *client) WithoutLock(role, environment, name string) *lockContext {
	return this.withLock(role, environment, name, false)
}

func (this *client) WithLock(role, environment, name string) *lockContext {
	return this.withLock(role, environment, name, true)
}

func (this *client) withLock(role, environment, name string, acquireLock bool) *lockContext {
	var lock *api.Lock
	var err error
	if acquireLock {
		lock, err = this.AcquireLock(role, environment, name)
		glog.Infoln("role=", role, "environment=", environment, "name=", name, "lock=", lock, "err=", err)
	}
	return &lockContext{
		Locking: acquireLock,
		lock:    lock,
		Err:     err,
		client:  this,
	}
}

func (this *client) GetJobs(role string) ([]*api.JobConfiguration, error) {
	resp, err := this.Readonly.GetJobs(role)
	if err != nil {
		return nil, err
	}
	result := []*api.JobConfiguration{}
	for config, _ := range resp.GetResult_().GetGetJobsResult_().Configs {
		result = append(result, config)
	}
	return result, nil
}

func (this *client) GetJobSummary(role string) ([]*api.JobSummary, error) {
	resp, err := this.Readonly.GetJobSummary(role)
	if err != nil {
		return nil, err
	}
	result := []*api.JobSummary{}
	for summary, _ := range resp.GetResult_().GetJobSummaryResult_().Summaries {
		result = append(result, summary)
	}
	return result, nil
}

func (this *client) GetQuota(role string) (*api.GetQuotaResult_, error) {
	resp, err := this.Readonly.GetQuota(role)
	if err != nil {
		return nil, err
	}
	return resp.GetResult_().GetGetQuotaResult_(), nil
}

// There is a bug currently with Aurora that it doesn't handle empty lists properly and will therefore
// not be able to match results.
func (this *client) GetJobTasks(role, environment, name string) ([]*api.ScheduledTask, error) {
	resp, err := this.Readonly.GetTasksStatus(&api.TaskQuery{
		JobKeys: map[*api.JobKey]bool{
			&api.JobKey{
				Role:        role,
				Environment: environment,
				Name:        name,
			}: true,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.GetResult_().GetScheduleStatusResult_().Tasks, nil
}

// TODO - This fails on HTTP Response code 401 - need to figure out auth.
func (this *client) CreateJob(job *api.JobConfiguration) (bool, error) {
	glog.Infoln("Creating Job:", job)
	resp, err := this.SchedulerManager.CreateJob(job, &api.Lock{
		Key: &api.LockKey{
			Job: job.Key,
		},
	})
	glog.Infoln("resp=", resp, "err=", err)
	if err != nil {
		return false, err
	}
	return resp.ResponseCode == api.ResponseCode_OK, nil
}

func (this *client) CreateJob2(job *api.JobConfiguration) (bool, error) {
	success := new(bool)
	err := this.WithLock(job.Key.Role, job.Key.Environment, job.Key.Name).Run(
		func(lock *api.Lock) error {
			glog.Infoln("Creating Job:", job)
			resp, err := this.SchedulerManager.CreateJob(job, nil)
			glog.Infoln("Lock=", lock, "resp=", resp, "err=", err)
			if err != nil {
				return err
			}
			*success = resp.ResponseCode == api.ResponseCode_OK
			return nil
		})
	return *success, err
}
