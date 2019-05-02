module github.com/mia0x75/venus

replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.36.0

	golang.org/x/build => github.com/golang/build v0.0.0-20190228010158-44b79b8774a7
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190227175134-215aa809caaf
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190221220918-438050ddec5e
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190227174305-5b3e6a55c961
	golang.org/x/net => github.com/golang/net v0.0.0-20190227160552-c95aed5357e7
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190226205417-e64efc72b421
	golang.org/x/perf => github.com/golang/perf v0.0.0-20190124201629-844a5f5b46f4
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190226215855-775f8194d0f9
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20181108054448-85acf8d2951c
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190227232517-f0a709d59f0f
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.1.0
	google.golang.org/appengine => github.com/golang/appengine v1.4.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190227213309-4f5b463f9597
	google.golang.org/grpc => github.com/grpc/grpc-go v1.19.0
)

require (
	github.com/gorilla/sessions v1.1.3
	github.com/labstack/echo-contrib v0.0.0-20190220224852-7fa08ffe9442
	github.com/labstack/echo/v4 v4.1.4
	github.com/machinebox/graphql v0.2.2
	github.com/pkg/errors v0.8.1 // indirect
	github.com/radovskyb/watcher v1.0.6
	github.com/sirupsen/logrus v1.4.1
	github.com/toolkits/file v0.0.0-20160325033739-a5b3c5147e07
	github.com/toolkits/nux v0.0.0-20190312004434-44d006618852
	github.com/toolkits/slice v0.0.0-20141116085117-e44a80af2484 // indirect
	github.com/toolkits/sys v0.0.0-20170615103026-1f33b217ffaf // indirect
)
