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
	rand.Seed(time.Now().UTC().UnixNano())
	port := start + rand.Intn(end-start)
	return net.JoinHostPort("localhost", strconv.Itoa(port))	
}
