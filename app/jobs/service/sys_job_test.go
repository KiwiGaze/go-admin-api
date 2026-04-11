package service

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/cronjob"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"go-admin-api/app/jobs"
	"go-admin-api/app/jobs/models"
	"go-admin-api/common/actions"
	"go-admin-api/common/dto"
	commonmodels "go-admin-api/common/models"
	"go-admin-api/internal/testutil"
)

func TestSysJobStartJob(t *testing.T) {
	t.Run("disabled job returns existing disabled error", func(t *testing.T) {
		db := newJobServiceTestDB(t)
		seedJob(t, db, models.SysJob{JobId: 1, JobName: "disabled", Status: 1, JobType: 1, CronExpression: "0 0 0 1 1 ?", InvokeTarget: "http://127.0.0.1"})
		service := newJobService(db, cronjob.NewWithSeconds())

		err := service.StartJob(&dto.GeneralGetDto{Id: 1}, &actions.DataPermission{})
		if err == nil || !strings.Contains(err.Error(), "disabled") {
			t.Fatalf("StartJob() error = %v, want disabled error", err)
		}
	})

	t.Run("already started job is rejected", func(t *testing.T) {
		db := newJobServiceTestDB(t)
		seedJob(t, db, models.SysJob{JobId: 1, JobName: "running", Status: 2, JobType: 1, EntryId: 7, CronExpression: "0 0 0 1 1 ?", InvokeTarget: "http://127.0.0.1"})
		crontab := cronjob.NewWithSeconds()
		service := newJobService(db, crontab)

		err := service.StartJob(&dto.GeneralGetDto{Id: 1}, &actions.DataPermission{})
		if !errors.Is(err, ErrJobAlreadyStarted) {
			t.Fatalf("StartJob() error = %v, want %v", err, ErrJobAlreadyStarted)
		}
		if entries := crontab.Entries(); len(entries) != 0 {
			t.Fatalf("cron entries = %d, want 0", len(entries))
		}
	})

	t.Run("nil scheduler returns initialization error", func(t *testing.T) {
		db := newJobServiceTestDB(t)
		seedJob(t, db, models.SysJob{JobId: 1, JobName: "ready", Status: 2, JobType: 1, CronExpression: "0 0 0 1 1 ?", InvokeTarget: "http://127.0.0.1"})
		service := newJobService(db, nil)

		err := service.StartJob(&dto.GeneralGetDto{Id: 1}, &actions.DataPermission{})
		if !errors.Is(err, ErrSchedulerNotInitialized) {
			t.Fatalf("StartJob() error = %v, want %v", err, ErrSchedulerNotInitialized)
		}
	})

	t.Run("valid HTTP job stores new entry id", func(t *testing.T) {
		db := newJobServiceTestDB(t)
		seedJob(t, db, models.SysJob{JobId: 1, JobName: "http", Status: 2, JobType: 1, CronExpression: "0 0 0 1 1 ?", InvokeTarget: "http://127.0.0.1"})
		service := newJobService(db, cronjob.NewWithSeconds())

		if err := service.StartJob(&dto.GeneralGetDto{Id: 1}, &actions.DataPermission{}); err != nil {
			t.Fatalf("StartJob() error = %v", err)
		}
		assertEntryID(t, db, 1, true)
	})

	t.Run("valid exec job stores new entry id", func(t *testing.T) {
		db := newJobServiceTestDB(t)
		jobs.InitJob()
		seedJob(t, db, models.SysJob{JobId: 1, JobName: "exec", Status: 2, JobType: 2, CronExpression: "0 0 0 1 1 ?", InvokeTarget: "ExamplesOne", Args: "arg"})
		service := newJobService(db, cronjob.NewWithSeconds())

		if err := service.StartJob(&dto.GeneralGetDto{Id: 1}, &actions.DataPermission{}); err != nil {
			t.Fatalf("StartJob() error = %v", err)
		}
		assertEntryID(t, db, 1, true)
	})

	t.Run("self scoped user can start own job", func(t *testing.T) {
		enableDataPermissionForTest(t)

		db := newJobServiceTestDB(t)
		seedJob(t, db, models.SysJob{
			JobId:          1,
			JobName:        "owned",
			Status:         2,
			JobType:        1,
			CronExpression: "0 0 0 1 1 ?",
			InvokeTarget:   "http://127.0.0.1",
			ControlBy:      commonmodels.ControlBy{CreateBy: 1},
		})
		service := newJobService(db, cronjob.NewWithSeconds())

		if err := service.StartJob(&dto.GeneralGetDto{Id: 1}, selfPermission(1)); err != nil {
			t.Fatalf("StartJob() error = %v", err)
		}
		assertEntryID(t, db, 1, true)
	})

	t.Run("self scoped user cannot start another user's job", func(t *testing.T) {
		enableDataPermissionForTest(t)

		db := newJobServiceTestDB(t)
		crontab := cronjob.NewWithSeconds()
		seedJob(t, db, models.SysJob{
			JobId:          1,
			JobName:        "hidden",
			Status:         2,
			JobType:        1,
			CronExpression: "0 0 0 1 1 ?",
			InvokeTarget:   "http://127.0.0.1",
			ControlBy:      commonmodels.ControlBy{CreateBy: 2},
		})
		service := newJobService(db, crontab)

		err := service.StartJob(&dto.GeneralGetDto{Id: 1}, selfPermission(1))
		if !errors.Is(err, ErrJobNotFoundOrNotVisible) {
			t.Fatalf("StartJob() error = %v, want %v", err, ErrJobNotFoundOrNotVisible)
		}
		assertEntryID(t, db, 1, false)
		if entries := crontab.Entries(); len(entries) != 0 {
			t.Fatalf("cron entries = %d, want 0", len(entries))
		}
	})
}

func TestSysJobRemoveJob(t *testing.T) {
	t.Run("nil scheduler returns initialization error", func(t *testing.T) {
		db := newJobServiceTestDB(t)
		seedJob(t, db, models.SysJob{JobId: 1, JobName: "running", Status: 2, EntryId: 3})
		service := newJobService(db, nil)

		err := service.RemoveJob(&dto.GeneralDelDto{Id: 1}, &actions.DataPermission{})
		if !errors.Is(err, ErrSchedulerNotInitialized) {
			t.Fatalf("RemoveJob() error = %v, want %v", err, ErrSchedulerNotInitialized)
		}
	})

	t.Run("already stopped job succeeds without removal", func(t *testing.T) {
		db := newJobServiceTestDB(t)
		seedJob(t, db, models.SysJob{JobId: 1, JobName: "stopped", Status: 2, EntryId: 0})
		service := newJobService(db, cronjob.NewWithSeconds())

		if err := service.RemoveJob(&dto.GeneralDelDto{Id: 1}, &actions.DataPermission{}); err != nil {
			t.Fatalf("RemoveJob() error = %v", err)
		}
		if service.Msg != "Job is not running." {
			t.Fatalf("service msg = %q, want stopped message", service.Msg)
		}
		assertEntryID(t, db, 1, false)
	})

	t.Run("valid removal resets entry id", func(t *testing.T) {
		db := newJobServiceTestDB(t)
		crontab := cronjob.NewWithSeconds()
		entryID, err := jobs.AddJob(crontab, &jobs.HttpJob{JobCore: jobs.JobCore{
			Name:           "http",
			CronExpression: "0 0 0 1 1 ?",
			InvokeTarget:   "http://127.0.0.1",
		}})
		if err != nil {
			t.Fatalf("AddJob() error = %v", err)
		}
		seedJob(t, db, models.SysJob{JobId: 1, JobName: "running", Status: 2, EntryId: entryID})
		service := newJobService(db, crontab)

		if err := service.RemoveJob(&dto.GeneralDelDto{Id: 1}, &actions.DataPermission{}); err != nil {
			t.Fatalf("RemoveJob() error = %v", err)
		}
		assertEntryID(t, db, 1, false)
	})

	t.Run("self scoped user can remove own job", func(t *testing.T) {
		enableDataPermissionForTest(t)

		db := newJobServiceTestDB(t)
		crontab := cronjob.NewWithSeconds()
		entryID, err := jobs.AddJob(crontab, &jobs.HttpJob{JobCore: jobs.JobCore{
			Name:           "owned",
			CronExpression: "0 0 0 1 1 ?",
			InvokeTarget:   "http://127.0.0.1",
		}})
		if err != nil {
			t.Fatalf("AddJob() error = %v", err)
		}
		seedJob(t, db, models.SysJob{
			JobId:     1,
			JobName:   "owned",
			Status:    2,
			EntryId:   entryID,
			ControlBy: commonmodels.ControlBy{CreateBy: 1},
		})
		service := newJobService(db, crontab)

		if err := service.RemoveJob(&dto.GeneralDelDto{Id: 1}, selfPermission(1)); err != nil {
			t.Fatalf("RemoveJob() error = %v", err)
		}
		assertEntryID(t, db, 1, false)
	})

	t.Run("self scoped user cannot remove another user's job", func(t *testing.T) {
		enableDataPermissionForTest(t)

		db := newJobServiceTestDB(t)
		crontab := cronjob.NewWithSeconds()
		entryID, err := jobs.AddJob(crontab, &jobs.HttpJob{JobCore: jobs.JobCore{
			Name:           "hidden",
			CronExpression: "0 0 0 1 1 ?",
			InvokeTarget:   "http://127.0.0.1",
		}})
		if err != nil {
			t.Fatalf("AddJob() error = %v", err)
		}
		seedJob(t, db, models.SysJob{
			JobId:     1,
			JobName:   "hidden",
			Status:    2,
			EntryId:   entryID,
			ControlBy: commonmodels.ControlBy{CreateBy: 2},
		})
		service := newJobService(db, crontab)

		err = service.RemoveJob(&dto.GeneralDelDto{Id: 1}, selfPermission(1))
		if !errors.Is(err, ErrJobNotFoundOrNotVisible) {
			t.Fatalf("RemoveJob() error = %v, want %v", err, ErrJobNotFoundOrNotVisible)
		}
		assertEntryID(t, db, 1, true)
	})
}

func newJobServiceTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	return testutil.NewTestDB(t, &models.SysJob{})
}

func newJobService(db *gorm.DB, crontab *cron.Cron) *SysJob {
	return &SysJob{
		Service: testutil.NewTestService(db),
		Cron:    crontab,
	}
}

func seedJob(t *testing.T, db *gorm.DB, job models.SysJob) {
	t.Helper()
	if err := db.Create(&job).Error; err != nil {
		t.Fatalf("seed job: %v", err)
	}
}

func assertEntryID(t *testing.T, db *gorm.DB, jobID int, wantSet bool) {
	t.Helper()
	var job models.SysJob
	if err := db.First(&job, "job_id = ?", jobID).Error; err != nil {
		t.Fatalf("load job: %v", err)
	}
	if wantSet && job.EntryId == 0 {
		t.Fatalf("entry id = 0, want set")
	}
	if !wantSet && job.EntryId != 0 {
		t.Fatalf("entry id = %d, want 0", job.EntryId)
	}
}

func selfPermission(userID int) *actions.DataPermission {
	return &actions.DataPermission{
		DataScope: "5",
		UserId:    userID,
	}
}

func enableDataPermissionForTest(t *testing.T) {
	t.Helper()
	previous := config.ApplicationConfig.EnableDP
	config.ApplicationConfig.EnableDP = true
	t.Cleanup(func() {
		config.ApplicationConfig.EnableDP = previous
	})
}
