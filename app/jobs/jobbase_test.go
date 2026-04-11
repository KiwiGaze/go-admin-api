package jobs

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-admin-team/go-admin-core/sdk"
	"gorm.io/gorm"

	"go-admin-api/app/jobs/models"
	"go-admin-api/internal/testutil"
)

func TestSetupInitializesEveryTenantAndReturns(t *testing.T) {
	tenantA := "tenant-a-" + t.Name()
	tenantB := "tenant-b-" + t.Name()
	dbA := testutil.NewTestDB(t, &models.SysJob{})
	dbB := testutil.NewTestDB(t, &models.SysJob{})
	seedSetupJob(t, dbA, 1)
	seedSetupJob(t, dbB, 2)

	done := make(chan struct{})
	go func() {
		Setup(map[string]*gorm.DB{
			tenantA: dbA,
			tenantB: dbB,
		})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Setup() did not return")
	}

	crontabA := sdk.Runtime.GetCrontabByTenant(tenantA)
	crontabB := sdk.Runtime.GetCrontabByTenant(tenantB)
	if crontabA == nil || crontabB == nil {
		t.Fatalf("crontabs = (%v, %v), want both initialized", crontabA, crontabB)
	}
	t.Cleanup(func() {
		crontabA.Stop()
		crontabB.Stop()
	})

	assertSetupEntryID(t, dbA, 1)
	assertSetupEntryID(t, dbB, 2)
}

func TestExecJobRunReturnsUnderlyingFailure(t *testing.T) {
	previousJobList := jobList
	jobList = map[string]JobExec{
		"failing": failingJobExec{},
	}
	t.Cleanup(func() {
		jobList = previousJobList
	})

	job := ExecJob{
		JobCore: JobCore{
			Name:         "failing",
			InvokeTarget: "failing",
			Args:         "arg",
		},
	}

	err := job.run()
	if err == nil || err.Error() != "boom" {
		t.Fatalf("run() error = %v, want boom", err)
	}
}

func TestHttpJobRunTreatsHTTPFailuresAsErrors(t *testing.T) {
	previousSleep := sleep
	previousRetryCount := retryCount
	sleep = func(time.Duration) {}
	retryCount = 2
	t.Cleanup(func() {
		sleep = previousSleep
		retryCount = previousRetryCount
	})

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("failed"))
	}))
	defer server.Close()

	job := HttpJob{
		JobCore: JobCore{
			Name:         "http-failing",
			InvokeTarget: server.URL,
		},
	}

	err := job.run()
	if err == nil || err.Error() != "received HTTP 500" {
		t.Fatalf("run() error = %v, want received HTTP 500", err)
	}
	if requestCount != retryCount {
		t.Fatalf("request count = %d, want %d", requestCount, retryCount)
	}
}

func TestHttpJobRunSucceedsOnOKResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	job := HttpJob{
		JobCore: JobCore{
			Name:         "http-ok",
			InvokeTarget: server.URL,
		},
	}

	if err := job.run(); err != nil {
		t.Fatalf("run() error = %v", err)
	}
}

func seedSetupJob(t *testing.T, db *gorm.DB, jobID int) {
	t.Helper()
	if err := db.Create(&models.SysJob{
		JobId:          jobID,
		JobName:        "setup",
		JobType:        1,
		Status:         2,
		CronExpression: "0 0 0 1 1 ?",
		InvokeTarget:   "http://127.0.0.1",
	}).Error; err != nil {
		t.Fatalf("seed setup job: %v", err)
	}
}

func assertSetupEntryID(t *testing.T, db *gorm.DB, jobID int) {
	t.Helper()
	var job models.SysJob
	if err := db.First(&job, "job_id = ?", jobID).Error; err != nil {
		t.Fatalf("load setup job: %v", err)
	}
	if job.EntryId == 0 {
		t.Fatalf("entry id for job %d = 0, want set", jobID)
	}
}

type failingJobExec struct{}

func (failingJobExec) Exec(arg interface{}) error {
	return errors.New("boom")
}
