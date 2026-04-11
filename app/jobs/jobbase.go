package jobs

import (
	"fmt"
	models2 "go-admin-api/app/jobs/models"
	"io"
	"net/http"
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	"gorm.io/gorm"

	"github.com/robfig/cron/v3"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/cronjob"
)

var timeFormat = "2006-01-02 15:04:05"
var retryCount = 3
var sleep = time.Sleep

var jobList map[string]JobExec

//var lock sync.Mutex

type JobCore struct {
	InvokeTarget   string
	Name           string
	JobId          int
	EntryId        int
	CronExpression string
	Args           string
}

// HttpJob is an HTTP job type.
type HttpJob struct {
	JobCore
}

type ExecJob struct {
	JobCore
}

func (e *ExecJob) Run() {
	startTime := time.Now()
	err := e.run()
	latencyTime := time.Since(startTime)
	if err != nil {
		log.Errorf("[Job] JobCore %s exec failed, spend :%v, error: %v", e.Name, latencyTime, err)
		return
	}
	log.Infof("[Job] JobCore %s exec success , spend :%v", e.Name, latencyTime)
}

// Run executes an HTTP job.
func (h *HttpJob) Run() {
	startTime := time.Now()
	err := h.run()
	latencyTime := time.Since(startTime)
	if err != nil {
		log.Errorf("[Job] JobCore %s exec failed after retries, spend :%v, error: %v", h.Name, latencyTime, err)
		return
	}
	log.Infof("[Job] JobCore %s exec success , spend :%v", h.Name, latencyTime)
}

func (e *ExecJob) run() error {
	obj, ok := jobList[e.InvokeTarget]
	if !ok || obj == nil {
		return fmt.Errorf("job %q is not registered", e.InvokeTarget)
	}
	return CallExec(obj, e.Args)
}

func (h *HttpJob) run() error {
	var lastErr error
	for attempt := 0; attempt < retryCount; attempt++ {
		responseBody, err := runHTTPGet(h.InvokeTarget)
		if err == nil {
			return nil
		}

		lastErr = err
		log.Warnf("[Job] mission failed! %v", err)
		log.Warnf("[Job] Retry after the task fails %d seconds! %s \n", (attempt+1)*5, responseBody)
		sleep(time.Duration(attempt+1) * 5 * time.Second)
	}
	return lastErr
}

func runHTTPGet(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return "", readErr
	}
	responseBody := string(body)
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return responseBody, fmt.Errorf("received HTTP %d", resp.StatusCode)
	}

	return responseBody, nil
}

// Setup initializes scheduled jobs.
func Setup(dbs map[string]*gorm.DB) {

	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Starting...")

	for k, db := range dbs {
		sdk.Runtime.SetCrontabByTenant(k, cronjob.NewWithSeconds())
		setup(k, db)
	}
}

func setup(key string, db *gorm.DB) {
	crontab := sdk.Runtime.GetCrontabByTenant(key)
	sysJob := models2.SysJob{}
	jobList := make([]models2.SysJob, 0)
	err := sysJob.GetList(db, &jobList)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore init error", err)
	}
	if len(jobList) == 0 {
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore total:0")
	}

	_, err = sysJob.RemoveAllEntryID(db)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore remove entry_id error", err)
	}

	for i := 0; i < len(jobList); i++ {
		if jobList[i].JobType == 1 {
			j := &HttpJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName

			sysJob.EntryId, err = AddJob(crontab, j)
		} else if jobList[i].JobType == 2 {
			j := &ExecJob{}
			j.InvokeTarget = jobList[i].InvokeTarget
			j.CronExpression = jobList[i].CronExpression
			j.JobId = jobList[i].JobId
			j.Name = jobList[i].JobName
			j.Args = jobList[i].Args
			sysJob.EntryId, err = AddJob(crontab, j)
		}
		err = sysJob.Update(db, jobList[i].JobId)
	}

	// Start tasks.
	crontab.Start()
	fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore start success.")
}

// AddJob adds a job. AddJob(invokeTarget string, jobId int, jobName string, cronExpression string)
func AddJob(c *cron.Cron, job Job) (int, error) {
	if job == nil {
		fmt.Println("unknown")
		return 0, nil
	}
	return job.addJob(c)
}

func (h *HttpJob) addJob(c *cron.Cron) (int, error) {
	id, err := c.AddJob(h.CronExpression, h)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

func (e *ExecJob) addJob(c *cron.Cron) (int, error) {
	id, err := c.AddJob(e.CronExpression, e)
	if err != nil {
		fmt.Println(time.Now().Format(timeFormat), " [ERROR] JobCore AddJob error", err)
		return 0, err
	}
	EntryId := int(id)
	return EntryId, nil
}

// Remove removes a task.
func Remove(c *cron.Cron, entryID int) chan bool {
	ch := make(chan bool)
	go func() {
		c.Remove(cron.EntryID(entryID))
		fmt.Println(time.Now().Format(timeFormat), " [INFO] JobCore Remove success ,info entryID :", entryID)
		ch <- true
	}()
	return ch
}

// Stop tasks.
//func Stop() chan bool {
//	ch := make(chan bool)
//	go func() {
//		global.GADMCron.Stop()
//		ch <- true
//	}()
//	return ch
//}
