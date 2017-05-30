package main

import (
	"github.com/Sirupsen/logrus"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func main() {
	unixSocketPath := "/run/docker/plugins/ipamproxy.sock"
	httpReverseProxiedHost := os.Getenv("IPAM_HTTP_PROXY_HOST")
	requestLoggingEnabled := strings.ToLower(os.Getenv("REQUEST_LOGGING_ENABLED")) == "true"

	logger := logrus.New()
	log.SetOutput(logger.Writer())

	if httpReverseProxiedHost == "" {
		logger.Fatal("KEY=VALUE of IPAM_HTTP_PROXY_HOST=HOST must be be set when driver is installed, or set when installed driver is inactive")
	}

	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = httpReverseProxiedHost
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if requestLoggingEnabled {
			requestDump, err := httputil.DumpRequest(r, true)
			if err != nil {
				logger.Fatal(err)
			}
			logger.Info(string(requestDump))
		}

		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})

	unixListener, err := net.Listen("unix", unixSocketPath)
	if err != nil {
		logger.Fatal("could not bind to " + unixSocketPath + "; make sure it does not already exist, and can be created")
	}
	logger.Fatal(http.Serve(unixListener, nil))
}
