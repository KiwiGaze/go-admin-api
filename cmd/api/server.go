package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/config/source/file"
	log "github.com/go-admin-team/go-admin-core/logger"
	coresdk "github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"go-admin-api/app/admin/models"
	"go-admin-api/app/admin/router"

	"go-admin-api/app/jobs"
	"go-admin-api/common/database"
	"go-admin-api/common/global"
	common "go-admin-api/common/middleware"
	appstorage "go-admin-api/common/storage"
	ext "go-admin-api/config"
	localsdk "go-admin-api/sdk"
	localstorage "go-admin-api/storage"
)

var (
	configYml string
	apiCheck  bool
	StartCmd  = &cobra.Command{
		Use:          "server",
		Short:        "Start API server",
		Example:      "go-admin server -c config/settings.yml",
		SilenceUsage: true,
		PreRun: func(cmd *cobra.Command, args []string) {
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

var AppRouters = make([]func(), 0)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().BoolVarP(&apiCheck, "api", "a", false, "Start server with check api data")

	// Register routes. fixme: for other application routes, create a new file in this directory and add them in init.
	AppRouters = append(AppRouters, router.InitRouter)
}

func setup() {
	// Inject extended configuration options.
	config.ExtendConfig = &ext.ExtConfig
	// 1. Load configuration.
	config.Setup(
		file.NewSource(file.WithPath(configYml)),
		database.Setup,
		appstorage.Setup,
	)
	// Register listeners.
	queue := localsdk.Runtime.GetMemoryQueue("")
	queue.Register(global.LoginLog, func(message localstorage.Messager) error {
		return models.SaveLoginLog(message)
	})
	queue.Register(global.OperateLog, func(message localstorage.Messager) error {
		return models.SaveOperaLog(message)
	})
	queue.Register(global.ApiCheck, func(message localstorage.Messager) error {
		return models.SaveSysApi(message)
	})
	go queue.Run()

	usageStr := `starting api server...`
	log.Info(usageStr)
}

func run() error {
	if config.ApplicationConfig.Mode == pkg.ModeProd.String() {
		gin.SetMode(gin.ReleaseMode)
	}
	initRouter()

	for _, f := range AppRouters {
		f()
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.ApplicationConfig.Host, config.ApplicationConfig.Port),
		Handler:      coresdk.Runtime.GetEngine(),
		ReadTimeout:  time.Duration(config.ApplicationConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(config.ApplicationConfig.WriterTimeout) * time.Second,
	}

	go func() {
		jobs.InitJob()
		jobs.Setup(coresdk.Runtime.GetAllDb())

	}()

	if apiCheck {
		var routers = coresdk.Runtime.GetRouter()
		q := localsdk.Runtime.GetMemoryQueue("")
		mp := make(map[string]interface{})
		mp["List"] = routers
		message, err := localsdk.Runtime.GetStreamMessage("", global.ApiCheck, mp)
		if err != nil {
			log.Infof("GetStreamMessage error, %s \n", err.Error())
			// Log the error, but do not interrupt the request.
		} else {
			err = q.Append(message)
			if err != nil {
				log.Infof("Append message error, %s \n", err.Error())
			}
		}
	}

	go func() {
		// Service connection.
		if config.SslConfig.Enable {
			if err := srv.ListenAndServeTLS(config.SslConfig.Pem, config.SslConfig.KeyStr); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("listen: ", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatal("listen: ", err)
			}
		}
	}()
	fmt.Println(pkg.Red(string(global.LogoContent)))
	tip()
	fmt.Println(pkg.Green("Server run at:"))
	fmt.Printf("-  Local:   %s://localhost:%d/ \r\n", "http", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: %s://%s:%d/ \r\n", "http", pkg.GetLocalHost(), config.ApplicationConfig.Port)
	fmt.Println(pkg.Green("Swagger run at:"))
	fmt.Printf("-  Local:   http://localhost:%d/swagger/admin/index.html \r\n", config.ApplicationConfig.Port)
	fmt.Printf("-  Network: %s://%s:%d/swagger/admin/index.html \r\n", "http", pkg.GetLocalHost(), config.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", pkg.GetCurrentTimeStr())
	// Wait for an interrupt signal to gracefully shut down the server, with a 5-second timeout.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Info("Shutdown Server ... ")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Info("Server exiting")

	return nil
}

//var Router runtime.Router

func tip() {
	usageStr := `Welcome to ` + pkg.Green(`go-admin `+global.Version) + `, use ` + pkg.Red(`-h`) + ` to view commands`
	fmt.Printf("%s \n\n", usageStr)
}

func initRouter() {
	var r *gin.Engine
	h := coresdk.Runtime.GetEngine()
	if h == nil {
		h = gin.New()
		coresdk.Runtime.SetEngine(h)
	}
	switch engine := h.(type) {
	case *gin.Engine:
		r = engine
	default:
		log.Fatal("not support other engine")
		//os.Exit(-1)
	}
	if config.SslConfig.Enable {
		// ToDo: Implement SSL handler
		// r.Use(handler.TlsHandler())
	}
	//r.Use(middleware.Metrics())
	r.Use(common.Sentinel()).
		Use(common.RequestId(pkg.TrafficKey)).
		Use(api.SetRequestLogger)

	common.InitMiddleware(r)

}
