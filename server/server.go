package main

import (
	"DSLab1-209319K/common"
	"io"
	"net/http"
	"net/rpc"
)

func main() {
	// create a `*College` object
	market := common.NewMarket()

	// register `market` object with `rpc.DefaultServer`
	rpc.Register(market)

	// register an HTTP handler for RPC communication on `http.DefaultServeMux` (default)
	// registers a handler on the `rpc.DefaultRPCPath` endpoint to respond to RPC messages
	// registers a handler on the `rpc.DefaultDebugPath` endpoint for debugging
	rpc.HandleHTTP()

	// sample test endpoint
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		io.WriteString(res, "RPC SERVER LIVE!")
	})

	// listen and serve default HTTP server
	http.ListenAndServe(":9000", nil)
}
