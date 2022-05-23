package providers

import (
	"testing"
	"time"

	"github.com/mcuadros/go-defaults"
	"gotest.tools/assert"
)

func TestConfigurationDefault(t *testing.T) {
	c := Option{}
	defaults.SetDefaults(&c)
	assert.Equal(t, c.RetryCount, 3)
	assert.Equal(t, c.RetryInterval, 1*time.Second)
	assert.Equal(t, c.RequestTimeout, 3*time.Second)

	c = Option{RetryCount: 10}
	defaults.SetDefaults(&c)
	assert.Equal(t, c.RetryCount, 10)
	assert.Equal(t, c.RetryInterval, 1*time.Second)
	assert.Equal(t, c.RequestTimeout, 3*time.Second)
}
