package router

import (
	"os"

	common "go-admin-api/common/middleware"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
)

// InitRouter initializes routes. Don't doubt it — this is in use.
func InitRouter() {
	var r *gin.Engine
	h := sdk.Runtime.GetEngine()
	if h == nil {
		log.Fatal("not found engine...")
		os.Exit(-1)
	}
	switch engine := h.(type) {
	case *gin.Engine:
		r = engine
	default:
		log.Fatal("not support other engine")
		os.Exit(-1)
	}

	// the jwt middleware
	authMiddleware, err := common.AuthInit()
	if err != nil {
		log.Fatalf("JWT Init Error, %s", err.Error())
	}

	// Register system routes
	InitSysRouter(r, authMiddleware)

	// Register business routes
	// TODO: Place business routes here. Currently contains only demo code, no real routes.
	InitExamplesRouter(r, authMiddleware)
}
