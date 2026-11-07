package config

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}

func (s *ConfigSuite) TestNewConfigTest() {
	s.Suite.T().Run("NewConfig успешный с полным конфигом", func(t *testing.T) {
		ctx := context.Background()

		dir := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("bot:\n  token: test-token\n"), 0o600))
		require.NoError(t, os.WriteFile(filepath.Join(dir, ".env"), []byte("BOT_TOKEN=test_token\n"), 0o600))

		newConfig, err := NewConfig(ctx, nil, dir)
		require.NoError(t, err)

		assert.NotNil(t, newConfig, "полный конфиг не пустой")
		assert.NotEmpty(t, newConfig.Token, "Token не пустой")
	})
}
