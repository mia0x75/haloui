package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"

	"github.com/mia0x75/venus/g"
	"github.com/mia0x75/venus/routes"
)

// var site struct {
// 	DefaultTitle string
// }

// var author struct {
// 	Name  string
// 	Email string
// }

// var page struct {
// 	Title string
// }

func main() {
	cfg := flag.String("c", "", "configuration file")
	version := flag.Bool("v", false, "show version")

	flag.Parse()

	fmt.Printf("%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n%-11s: %s\n\n",
		"Version", g.Version,
		"Git commit", g.Git,
		"Compile", g.Compile,
		"Distro", g.Distro,
		"Kernel", g.Kernel,
		"Branch", g.Branch,
	)
	if *version {
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitLog()

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = strings.EqualFold(g.Config().Log.Level, "debug")
	renderer := g.NewRenderer()
	e.Renderer = renderer
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		var code = http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		if !c.Response().Committed {
			page := fmt.Sprintf("%d.html", code)
			if _, err := renderer.GetTemplate(page); err != nil {
				c.NoContent(code)
			} else {
				c.Render(code, page, nil)
			}
		}
	}

	e.Static("/assets", "public/assets")
	e.Static("/demo", "public/demo")
	e.File("/favicon.ico", "public/assets/images/favicon.ico")

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "Echo/4.0")
			return next(c)
		}
	})
	e.Use(middleware.SecureWithConfig(middleware.DefaultSecureConfig))
	e.Use(middleware.MethodOverride())
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `time:${time_unix} remote_ip:${remote_ip} host:${host} ` +
			`method:${method} uri:${uri} status:${status} bytes_in:${bytes_in} ` +
			`bytes_out:${bytes_out}` + "\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret-key"))))

	routes.Routes(e)

	// Startup http service
	go func() {
		addr := g.Config().Listen
		log.Infof("[I] http listening %s", addr)
		e.Logger.Fatal(e.StartTLS(addr, g.Config().Cert, g.Config().Key))
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	sc := make(chan os.Signal)
	signal.Notify(sc,
		os.Kill,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sc
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer func() {
	}()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
