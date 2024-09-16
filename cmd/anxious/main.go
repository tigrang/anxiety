package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/tigrang/anxiety"
)

func main() {
	var notifyFlag bool
	var appFlag string
	var proxyUrlFlag string
	var proxyBindFlag string
	var buildRouteFlag string
	var cmdFlag string
	var connectTimeoutFlag int

	flag.BoolVar(&notifyFlag, "notify", false, "notify proxy to trigger build")
	flag.StringVar(&appFlag, "app", "", "path to app")
	flag.StringVar(&proxyUrlFlag, "proxy", "http://localhost:3000", "url app is listening on to forward requests")
	flag.StringVar(&buildRouteFlag, "buildroute", "/internal/build", "path to trigger builds (must be the same when --notify is used)")
	flag.StringVar(&proxyBindFlag, "proxybind", "localhost:9000", "the addr for error proxy to listen on")
	flag.StringVar(&cmdFlag, "cmd", "./build", "patht to script that will build app")
	flag.IntVar(&connectTimeoutFlag, "timeout", 30, "the number of seconds to wait for proxy to be available")
	flag.Parse()

	if notifyFlag {
		notify(proxyBindFlag, buildRouteFlag, connectTimeoutFlag)
		return
	}

	startServer(appFlag, cmdFlag, proxyUrlFlag, proxyBindFlag, buildRouteFlag)
}

func waitForConnection(addr string, connectTimeout int) (err error) {
	attempt := 0
	var conn net.Conn
	for attempt < connectTimeout {
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			time.Sleep(time.Second)
			attempt = attempt + 1
			continue
		}

		conn.Close()
		return nil
	}

	return err
}

func notify(proxybind, proxyRoute string, connectTimeout int) {
	if err := waitForConnection(proxybind, connectTimeout); err != nil {
		panic(err)
	}

	_, err := http.Get("http://" + proxybind + proxyRoute)
	if err != nil {
		panic(err)
	}
}

func startServer(app string, cmd string, proxyUrl string, proxyBind string, buildPath string) {
	h := anxiety.NewProxy(app, cmd, proxyUrl, buildPath)

	fmt.Println("Listening on " + proxyBind)
	fmt.Println("Proxy to " + proxyUrl)

	err := http.ListenAndServe(proxyBind, h)
	if err != nil {
		panic(err)
	}
}
