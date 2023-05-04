package internal

import (
	"math/rand"
	"net"
	"strconv"
)

const (
	startPort = 50000
	endPort   = 60000
)

func RandomLocalAddr() string {
	port := startPort + rand.Intn(endPort-startPort)
	return net.JoinHostPort("localhost", strconv.Itoa(port))
}
