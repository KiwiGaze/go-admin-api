package global

// In-memory queue topics for asynchronous log processing and API validation
const (
	// LoginLog Login log queue, records user login success/failure events
	LoginLog = "login_log_queue"
	// OperateLog Operation log queue, records API request operation logs
	OperateLog = "operate_log_queue"
	// ApiCheck API validation queue, validates consistency between registered routes and sys_api table
	ApiCheck = "api_check_queue"
)