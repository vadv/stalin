package main

import (
	"flag"
	"log"
	"net"
	"runtime"
	"stalin/transport"
)

var (
	bindto            = flag.String("Bind", "0.0.0.0:5555", "Listen riemann events")
	resendaddr        = flag.String("ResendAdr", "", "Resend to next riemann")
	resendflushtime   = flag.Duration("ResenderFlush", 1000000000, "Flush every")
	graphiteaddr      = flag.String("GraphiteAdr", "", "Send graphite address, like 'graphite-1.production.infra.home:2023'")
	graphiteflushtime = flag.Duration("GraphiteFlush", 1000000000, "Flush every")
	opentsdbaddr      = flag.String("OpentsDBAdr", "", "Send opentsdb address, like 'opentsdb-1.production.infra.home:4242'")
	opentsdbflushtime = flag.Duration("OpentsDBFlush", 1000000000, "Flush every")
	pgconnstring      = flag.String("PgConnString", "", "Postgresql connection string, like: 'postgres://riemann_face_user@127.0.0.1/riemann_face?sslmode=disable'")
	pgquery           = flag.String("PgQuery", "SELECT UpdateProblems($1, $2, $3, $4, $5, $6, $7)", "Postgresql query string")
	pgpool            = flag.Int("PgPool", 40, "Connection pool for pg")
	statistictime     = flag.Int("StatTime", 0, "Print statistic (s)")
)

func main() {
	// run Forest, run
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()
	log.Println("Listen tcp", *bindto)
	lst, err := net.Listen("tcp", *bindto)
	if err != nil {
		log.Fatal(err)
	}
	t, err := transport.NewTransport(*graphiteaddr, *opentsdbaddr, *resendaddr, *pgconnstring, *pgquery, *pgpool, *graphiteflushtime, *opentsdbflushtime, *resendflushtime, *statistictime)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := lst.Accept()
		if err != nil {
			log.Println("Listener error:", err)
		}
		go t.HandleConn(conn)
	}
}
