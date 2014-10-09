package main

import (
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"stalin/plugins"
	_ "stalin/plugins/elastic"
	_ "stalin/plugins/graphite"
	_ "stalin/plugins/http"
	_ "stalin/plugins/opentsdb"
	_ "stalin/plugins/pgsql"
	_ "stalin/plugins/router"
	_ "stalin/plugins/tcp"
	"syscall"
)

func main() {

	// run Forest, run
	runtime.GOMAXPROCS(runtime.NumCPU())

	configPath := flag.String("config", filepath.FromSlash("/etc/stalin.toml"), "Config file.")
	flag.Parse()
	plugins.Run(*configPath)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	for {
		<-c
		// todo: hup, usr1
		os.Exit(1)
	}

}
