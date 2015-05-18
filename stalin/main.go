package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"stalin/config"
	. "stalin/log"
	_ "stalin/plugins/graphite"
	_ "stalin/plugins/http"
	_ "stalin/plugins/opentsdb"
	_ "stalin/plugins/sandbox"
	_ "stalin/plugins/store"
	_ "stalin/plugins/tcp"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU() * 4)

	configPath := flag.String("config", filepath.FromSlash("/etc/stalin.toml"), "Config file.")
	pluginUsage := flag.String("plugin-usage", "", "Print plugin usage and exit. Use 'all' for all")
	flag.Parse()

	if *pluginUsage != "" {
		fmt.Printf(config.DescriptionForPlugin(*pluginUsage))
		os.Exit(1)
	}

	pipeline, err := config.ReadConfig(*configPath)
	if err != nil {
		Log.Error("Config error: %v", err)
		os.Exit(2)
	}
	pipeline.RunAll()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		http.ListenAndServe("0.0.0.0:6060", nil)
	}()

	for {
		<-c
		// todo: hup, usr1
		os.Exit(1)
	}

}
