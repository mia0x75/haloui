package routes

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// Routes TODO
func Routes(e *echo.Echo) {
	p := e.Group("/") // 公开组
	p.GET("about.html", about)
	p.Any("login.html", login)
	p.GET("register.html", register)
	p.GET("forgot-password.html", forgot)
	p.GET("400.html", error400)
	p.GET("401.html", error401)
	p.GET("402.html", error402)
	p.GET("403.html", error403)
	p.GET("404.html", error404)
	p.GET("500.html", error500)
	p.GET("503.html", error503)

	r := e.Group("/") // 限制组
	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)
			if _, ok := sess.Values["token"]; ok {
				sess.Save(c.Request(), c.Response())
				return next(c)
			}

			return c.Redirect(http.StatusFound, "/login.html")
		}
	})

	r.GET("index.html", dashboard)
	r.GET("users-list.html", users)
	r.GET("crypto-currencies.html", currencies)
	r.GET("pagination.html", pagination)
	r.GET("lookup.html", lookup)
	r.GET("invoice.html", invoice)
	r.GET("sample-cards.html", sample)
	r.GET("clusters-list.html", clusters)
	r.GET("crons-list.html", crons)
	r.GET("options-list.html", options)
	r.GET("queries-list.html", queries)
	r.GET("rules-list.html", rules)
	r.GET("tasks-list.html", tasks)
	r.GET("tickets-list.html", tickets)
	r.GET("users-list.html", users)

	// e.GET("/captcha", captcha.Server(captcha.StdWidth, captcha.StdHeight))
}
