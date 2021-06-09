package cased

import (
	"io/ioutil"
	"log"
)

// Logger ...
var Logger = log.New(ioutil.Discard, "[Cased] ", log.LstdFlags)
