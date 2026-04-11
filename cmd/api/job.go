package api

import "go-admin-api/app/jobs/router"

func init() {
	// Register routes for the jobs application.
	AppRouters = append(AppRouters, router.InitRouter)
}
