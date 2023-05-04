package internal

import (
	"math/rand"
	"net"
	"strconv"
	"time"
)

const (
	startPort = 50000
	endPort   = 60000
)

func RandomLocalAddr() string {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	port := startPort + r.Intn(endPort-startPort)
	return net.JoinHostPort("localhost", strconv.Itoa(port))
}
