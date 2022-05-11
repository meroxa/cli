package global

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalTurbineJSConfig(t *testing.T) {
	os.Setenv("MEROXA_USE_LOCAL_TURBINE_JS", "true")
	result := GetLocalTurbineJSSetting()
	assert.Equal(t, "true", result)
}
