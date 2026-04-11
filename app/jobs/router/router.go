package router

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
	routerCheckRole   = make([]func(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware), 0)
)

// initRouter registers routes.
func initRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {

	// Routes that do not require authentication.
	noCheckRoleRouter(r)
	// Routes that require authentication.
	checkRoleRouter(r, authMiddleware)

	return r
}

// noCheckRoleRouter registers routes that do not require authentication.
func noCheckRoleRouter(r *gin.Engine) {
	// Configure the API version according to business requirements.
	v1 := r.Group("/api/v1")

	for _, f := range routerNoCheckRole {
		f(v1)
	}
}

// checkRoleRouter registers routes that require authentication.
func checkRoleRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) {
	// Configure the API version according to business requirements.
	v1 := r.Group("/api/v1")

	for _, f := range routerCheckRole {
		f(v1, authMiddleware)
	}
}
