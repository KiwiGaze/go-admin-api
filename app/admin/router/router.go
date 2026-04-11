package router

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

var (
	routerNoCheckRole = make([]func(*gin.RouterGroup), 0)
	routerCheckRole   = make([]func(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware), 0)
)

func InitExamplesRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.Engine {

	// Routes that do not require authentication
	examplesNoCheckRoleRouter(r)
	// Routes that require authentication
	examplesCheckRoleRouter(r, authMiddleware)

	return r
}

// Example of routes that do not require authentication
func examplesNoCheckRoleRouter(r *gin.Engine) {
	// API version can be set according to business requirements
	v1 := r.Group("/api/v1")
	for _, f := range routerNoCheckRole {
		f(v1)
	}
}

// Example of routes that require authentication
func examplesCheckRoleRouter(r *gin.Engine, authMiddleware *jwtauth.GinJWTMiddleware) {
	// API version can be set according to business requirements
	v1 := r.Group("/api/v1")
	for _, f := range routerCheckRole {
		f(v1, authMiddleware)
	}
}
