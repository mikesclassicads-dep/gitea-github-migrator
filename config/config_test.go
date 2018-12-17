package config

import (
	"github.com/jinzhu/configor"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadConfig(t *testing.T) {
	assert.NoError(t, configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&Config, "../config.example.yml"), "Parse example config failed.")
}
