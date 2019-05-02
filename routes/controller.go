package routes

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/machinebox/graphql"
)

var (
	startTime = time.Now()
)

func init() {
}

func humanateBytes(s uint64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(math.Log(float64(s)) / math.Log(base))
	suffix := sizes[int(e)]
	val := float64(s) / math.Pow(base, math.Floor(e))
	f := "%.0f"
	if val < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", val, suffix)
}

// FileSize calculates the file size and generate user-friendly string.
func FileSize(s int64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(uint64(s), 1024, sizes)
}

func dashboard(c echo.Context) error {
	req := graphql.NewRequest(`query index {
  statistics (Groups: ["overall" "tickets-daily" "tickets-users"]) {
    Group
    Key
    Value
  }
  crons(first: 10) {
    edges {
      node {
        UUID
        Status
        Name
        Cmd
        Params
        Interval
        Duration
        LastRun
        NextRun
        Recurrent
        UpdateAt
        CreateAt
      }
    }
  }
  me {
    ...UserInfo
  }
  environments {
    CPUStats {
      Kernel
      Idle
      User
      Nice
      IOWait
      LoadMin1
      LoadMin5
      LoadMin15
    }
    HostInfos {
      NCPUs
      OSName
      HostName
      BitWidth
      OSRelease
      OSVersion
      Platform
    }
    MemStats {
      Used
      Free
      Cache
      Total
      SwapUsed
      SwapFree
      SwapTotal
    }
    ProcessStats {
      Stopped
      Running
      Zombie
      Total
    }
  }
  tickets(first: 10){
    edges {
      node {
        UUID
        Subject
        Database
        Status
        CreateAt
        User {
          ...UserInfo
        }
        Reviewer {
          ...UserInfo
        }
        Cluster {
          ...ClusterInfo
        }
      }
    }
  }
}
fragment UserInfo on User {
  Name
  UUID
  Avatar {
    URL
  }
}
fragment ClusterInfo on Cluster {
  UUID
  Alias
  Host
  IP
  Port
}
`)

	// set any variables
	// req.Var("email", "ces365@163.com")
	// req.Var("password", "root")

	sess, _ := session.Get("session", c)
	token, _ := sess.Values["token"].(string)

	// set header fields
	req.Header.Set("Authentication", token)
	req.Header.Set("Cache-Control", "no-cache")

	var resp struct {
		Me struct {
			UUID string
			Name string
		}
		Statistics []struct {
			Group string
			Key   string
			Value float64
		}
		Environments struct {
			CPUStats struct {
				User      float64
				Kernel    float64
				Idle      float64
				IOWait    float64
				Swap      float64
				Nice      float64
				LoadMin1  float64
				LoadMin5  float64
				LoadMin15 float64
			}
			HostInfos struct {
				NCPUs     uint
				OSName    string
				HostName  string
				BitWidth  uint
				OSRelease string
				OSVersion string
				Platform  string
			}
			MemStats struct {
				Used      float64
				Free      float64
				Cache     float64
				Total     float64
				SwapUsed  float64
				SwapFree  float64
				SwapTotal float64
			}
			ProcessStats struct {
				Stopped uint
				Running uint
				Zombie  uint
				Total   uint
			}
		}
		Crons struct {
			Edges []struct {
				Node struct {
					UUID      string
					Status    string
					Name      string
					Cmd       string
					Params    string
					Interval  string
					Duration  string
					LastRun   string
					NextRun   string
					Recurrent uint8
					UpdateAt  uint
					CreateAt  uint
				}
			}
		}
		Tickets struct {
			Edges []struct {
				Node struct {
					UUID     string
					Subject  string
					Database string
					Status   uint8
					CreateAt uint
					User     struct {
						Name   string
						UUID   string
						Avatar struct {
							URL string
						}
					}
					Reviewer struct {
						Name   string
						UUID   string
						Avatar struct {
							URL string
						}
					}
					Cluster struct {
						UUID  string
						Alias string
						Host  string
						IP    string
						Port  uint16
					}
				}
			}
		}
	}

	if err := request(req, &resp); err != nil {
		log.Println(err)
	}

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"data": resp,
	})
}

func users(c echo.Context) error {
	req := graphql.NewRequest(`query index {
  me {
    ...UserInfo
  }
  users (first: 15){
    edges {
      node {
        UUID
        Name
		Email
		Status
        Avatar {
          URL
        }
        Phone
        CreateAt
      }
    }
  }
}
fragment UserInfo on User {
  Name
  UUID
  Avatar {
    URL
  }
}
`)

	// set any variables
	// req.Var("email", "ces365@163.com")
	// req.Var("password", "root")

	sess, _ := session.Get("session", c)
	token, _ := sess.Values["token"].(string)

	// set header fields
	req.Header.Set("Authentication", token)
	req.Header.Set("Cache-Control", "no-cache")

	var resp struct {
		Me struct {
			UUID   string
			Name   string
			Avatar struct {
				URL string
			}
		}
		Users struct {
			Edges []struct {
				Node struct {
					UUID     string
					Name     string
					Email    string
					Phone    uint64
					CreateAt uint
					Status   uint8
					Avatar   struct {
						URL string
					}
				}
			}
		}
	}

	if err := request(req, &resp); err != nil {
		log.Println(err)
	}

	return c.Render(http.StatusOK, "users-list.html", map[string]interface{}{
		"data": resp,
	})
}

func tickets(c echo.Context) error {
	return c.Render(http.StatusOK, "tickets-list.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func rules(c echo.Context) error {
	req := graphql.NewRequest(`query index {
  me {
    ...UserInfo
  }
  rules {
    UUID
    Name
    Group
    Description
    VldrGroup
    Values
    Bitwise
    Func
    Element
    CreateAt
    UpdateAt
  }
}
fragment UserInfo on User {
  Name
  UUID
  Avatar {
    URL
  }
}
`)

	// set any variables
	// req.Var("email", "ces365@163.com")
	// req.Var("password", "root")

	sess, _ := session.Get("session", c)
	token, _ := sess.Values["token"].(string)

	// set header fields
	req.Header.Set("Authentication", token)
	req.Header.Set("Cache-Control", "no-cache")

	var resp struct {
		Me struct {
			UUID   string
			Name   string
			Avatar struct {
				URL string
			}
		}
		Rules []struct {
			UUID        string
			Name        string
			Group       uint8
			Description string
			VldrGroup   uint16
			Values      string
			Bitwise     uint8
			Func        string
			Element     string
			CreateAt    uint
			UpdateAt    uint
		}
	}

	if err := request(req, &resp); err != nil {
		log.Println(err)
	}

	return c.Render(http.StatusOK, "rules-list.html", map[string]interface{}{
		"data": resp,
	})
}

func clusters(c echo.Context) error {
	req := graphql.NewRequest(`query index {
  me {
    ...UserInfo
  }
  clusters (first: 15){
    edges {
      node {
        UUID
        Alias
        Host
        IP
        Port
        Status
        CreateAt
      }
    }
  }
}
fragment UserInfo on User {
  Name
  UUID
  Avatar {
    URL
  }
}
`)

	// set any variables
	// req.Var("email", "ces365@163.com")
	// req.Var("password", "root")

	sess, _ := session.Get("session", c)
	token, _ := sess.Values["token"].(string)

	// set header fields
	req.Header.Set("Authentication", token)
	req.Header.Set("Cache-Control", "no-cache")

	var resp struct {
		Me struct {
			UUID   string
			Name   string
			Avatar struct {
				URL string
			}
		}
		Clusters struct {
			Edges []struct {
				Node struct {
					UUID     string
					Alias    string
					Host     string
					IP       string
					CreateAt uint
					Status   uint8
					Port     uint16
				}
			}
		}
	}

	if err := request(req, &resp); err != nil {
		log.Println(err)
	}

	return c.Render(http.StatusOK, "clusters-list.html", map[string]interface{}{
		"data": resp,
	})
}

func queries(c echo.Context) error {
	return c.Render(http.StatusOK, "queries-list.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func crons(c echo.Context) error {
	return c.Render(http.StatusOK, "crons-list.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func tasks(c echo.Context) error {
	return c.Render(http.StatusOK, "tasks-list.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func options(c echo.Context) error {
	return c.Render(http.StatusOK, "options-list.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func currencies(c echo.Context) error {
	return c.Render(http.StatusOK, "crypto-currencies.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func pagination(c echo.Context) error {
	return c.Render(http.StatusOK, "pagination.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func lookup(c echo.Context) error {
	return c.Render(http.StatusOK, "lookup.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func invoice(c echo.Context) error {
	return c.Render(http.StatusOK, "invoice.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func sample(c echo.Context) error {
	return c.Render(http.StatusOK, "sample-cards.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func login(c echo.Context) error {
	switch c.Request().Method {
	case "GET":
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"name": "Dolly!",
		})
	case "POST":
		req := graphql.NewRequest(`mutation ($email: String! $password: String!) {
		login (
			input: {
				Email: $email
				Password: $password
			}
		) {
			Me {
				UUID
				Name
				Phone
			}
			Token
		}
	}`)

		// set any variables
		req.Var("email", "ces365@163.com")
		req.Var("password", "root")

		// set header fields
		req.Header.Set("Cache-Control", "no-cache")

		var resp struct {
			Login struct {
				Me struct {
					UUID     string
					Email    string
					Password string
					Status   uint8
					Name     string
					Phone    uint64
					UpdateAt uint
					CreateAt uint
				}
				Token string
			}
		}

		if err := request(req, &resp); err != nil {
			log.Println(err)
		}

		sess, _ := session.Get("session", c)
		sess.Values["token"] = resp.Login.Token
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, "index.html")
	default:
		fmt.Println(c.Request().Method)
		return echo.NewHTTPError(http.StatusBadRequest, "Method not allowed.")
	}
}

func register(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func forgot(c echo.Context) error {
	return c.Render(http.StatusOK, "forgot-password.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func error400(c echo.Context) error {
	return c.Render(http.StatusOK, "400.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func error401(c echo.Context) error {
	return c.Render(http.StatusOK, "401.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func error402(c echo.Context) error {
	return c.Render(http.StatusOK, "402.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func error403(c echo.Context) error {
	return c.Render(http.StatusOK, "403.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func error404(c echo.Context) error {
	return c.Render(http.StatusOK, "404.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func error500(c echo.Context) error {
	return c.Render(http.StatusOK, "500.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func error503(c echo.Context) error {
	return c.Render(http.StatusOK, "503.html", map[string]interface{}{
		"name": "Dolly!",
	})
}

func about(c echo.Context) error {
	return c.Render(http.StatusOK, "about.html", map[string]interface{}{
		"name":  "Dolly!",
		"title": "page_title",
	})
}

func request(req *graphql.Request, resp interface{}) (err error) {
	option := graphql.WithHTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})
	// create a client (safe to share across requests)
	client := graphql.NewClient("https://127.0.0.1:4000/api/query", option)
	// define a Context for the request
	ctx := context.Background()
	err = client.Run(ctx, req, resp)
	return
}
