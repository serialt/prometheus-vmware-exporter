package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"prometheus-vmware-exporter/controller"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
)

var (
	listen   = ":9879"
	host     = ""
	username = ""
	password = ""
	logLevel = "info"
	config   = ""
)

func env(key, def string) string {
	if x := os.Getenv(key); x != "" {
		return x
	}
	return def
}

func init() {
	flag.StringVar(&listen, "listen", env("ESX_LISTEN", listen), "listen port")
	flag.StringVar(&host, "host", env("ESX_HOST", host), "URL ESX host ")
	flag.StringVar(&username, "username", env("ESX_USERNAME", username), "User for ESX")
	flag.StringVar(&password, "password", env("ESX_PASSWORD", password), "password for ESX")
	flag.StringVar(&logLevel, "log", env("ESX_LOG", logLevel), "Log level must be, debug or info")
	flag.StringVar(&config, "config", env("CONFIG", config), "[EXPERIMENTAL] Path to config yaml file that can enable TLS or authentication.")
	flag.Parse()
	controller.RegistredMetrics()
	collectMetrics()
}

func collectMetrics() {
	go func() {
		slog.Debug("Start collect host metrics")
		controller.NewVmwareHostMetrics(host, username, password)
		slog.Debug("End collect host metrics")
	}()
	go func() {
		slog.Debug("Start collect datastore metrics")
		controller.NewVmwareDsMetrics(host, username, password)
		slog.Debug("End collect datastore metrics")
	}()
	go func() {
		slog.Debug("Start collect VM metrics")
		controller.NewVmwareVmMetrics(host, username, password)
		slog.Debug("End collect VM metrics")
	}()
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		collectMetrics()
	}
	h := promhttp.Handler()
	h.ServeHTTP(w, r)
}

func main() {
	if host == "" {
		slog.Error("Yor must configured systemm env ESX_HOST or key -host")
	}
	if username == "" {
		slog.Error("Yor must configured system env ESX_USERNAME or key -username")
	}
	if password == "" {
		slog.Error("Yor must configured system env ESX_PASSWORD or key -password")
	}
	msg := fmt.Sprintf("Exporter start on port %s", listen)
	slog.Info(msg)
	http.HandleFunc("/metrics", handler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>VMware Exporter</title></head>
			<body>
			<h1>VMware Exporter</h1>
			<p><a href="` + "/metrics" + `">Metrics</a></p>
			</body>
			</html>`))
	})
	flags := &web.FlagConfig{
		WebListenAddresses: &([]string{listen}),
		WebConfigFile:      &config,
	}
	server := &http.Server{Addr: listen}
	if err := web.ListenAndServe(server, flags, slog.Default()); err != nil {
		slog.Error("Service start failed", "err", err)
		os.Exit(1)
	}

}
