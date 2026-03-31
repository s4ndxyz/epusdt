package command

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/assimon/luuu/bootstrap"
	"github.com/assimon/luuu/config"
	"github.com/assimon/luuu/middleware"
	"github.com/assimon/luuu/route"
	"github.com/assimon/luuu/util/constant"
	luluHttp "github.com/assimon/luuu/util/http"
	"github.com/assimon/luuu/util/log"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "http service",
	Long:  "http service commands",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	httpCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start",
	Long:  "start http service",
	Run: func(cmd *cobra.Command, args []string) {
		bootstrap.InitApp()
		printBanner()
		HttpServerStart()
	},
}

func HttpServerStart() {
	var err error
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = customHTTPErrorHandler

	MiddlewareRegister(e)
	route.RegisterRoute(e)
	e.Static(config.StaticPath, config.StaticFilePath)

	httpListen := viper.GetString("http_listen")
	go func() {
		if err = e.Start(httpListen); err != http.ErrServerClosed {
			log.Sugar.Error(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func MiddlewareRegister(e *echo.Echo) {
	if config.HTTPAccessLog {
		e.Use(echoMiddleware.Logger())
	}
	e.Use(middleware.RequestUUID())
}

func customHTTPErrorHandler(err error, e echo.Context) {
	code := http.StatusInternalServerError
	msg := "server error"
	resp := &luluHttp.Response{
		StatusCode: code,
		Message:    msg,
		RequestID:  e.Request().Header.Get(echo.HeaderXRequestID),
	}
	if he, ok := err.(*echo.HTTPError); ok {
		e.String(http.StatusOK, he.Message.(string))
		return
	}
	if he, ok := err.(*constant.RspError); ok {
		resp.StatusCode = he.Code
		resp.Message = he.Msg
	}
	_ = e.JSON(http.StatusOK, resp)
}
