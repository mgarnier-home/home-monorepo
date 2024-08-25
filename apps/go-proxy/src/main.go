package main

import "mgarnier11/go-proxy/proxies"

func main() {
	println("Hello, World!")
	tcpProxy := &proxies.TCPProxy{ListenAddr: "localhost:8080", TargetAddr: "localhost:8081"}

	tcpProxy.ListenAddr = "localhost:8080"
	// proxies.TestUDP()
}
