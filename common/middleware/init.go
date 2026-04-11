package middleware

import (
	"go-admin-api/common/actions"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

const (
	JwtTokenCheck   string = "JwtToken"
	RoleCheck       string = "AuthCheckRole"
	PermissionCheck string = "PermissionAction"
)

func InitMiddleware(r *gin.Engine) {
	r.Use(RequestId(pkg.TrafficKey))
	r.Use(DemoEvn())
	// Database connection
	r.Use(WithContextDb)
	// Log handling
	r.Use(LoggerToFile())
	// Custom error handling
	r.Use(CustomError)
	// NoCache is a middleware function that appends headers
	r.Use(NoCache)
	// CORS handling
	r.Use(Options)
	// Secure is a middleware function that appends security
	r.Use(Secure)
	// Distributed tracing
	//r.Use(middleware.Trace())
	sdk.Runtime.SetMiddleware(JwtTokenCheck, (*jwt.GinJWTMiddleware).MiddlewareFunc)
	sdk.Runtime.SetMiddleware(RoleCheck, AuthCheckRole())
	sdk.Runtime.SetMiddleware(PermissionCheck, actions.PermissionAction())
}
