package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"os-simulator-plan/internal/transport/realtime"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()

	manager := realtime.NewSessionManager()
	server := realtime.NewServer(manager)

	fmt.Printf("server listening on %s\n", *addr)
	if err := http.ListenAndServe(*addr, server.Handler()); err != nil {
		fmt.Fprintf(os.Stderr, "server failed: %v\n", err)
		os.Exit(1)
	}
}
