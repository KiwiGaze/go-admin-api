package main

import (
	"go-admin-api/cmd"
)

//go:generate swag init --parseDependency --parseDepth=6 --instanceName admin -o ./docs/admin

// @title go-admin-api
// @version 1.0.0
// @description API documentation for a role-based access control system built with Gin + Vue + Element UI (front-end and back-end separated)
// @license.name MIT
// @license.url https://github.com/kiwi-gaze/go-admin-api/blob/master/LICENSE.md

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	cmd.Execute()
}
