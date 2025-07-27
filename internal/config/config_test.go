package config

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	t.Run("should load config from yaml file", func(t *testing.T) {
		// given
		path := filepath.Join("testdata", "test.yaml")

		// when
		cfg, err := Load(path)

		// then
		assert.NoErrorf(t, err, "Load(%s) should not return error on valid file", path)
		assert.Equalf(t, time.Minute*1+time.Second*30, cfg.RefreshInterval, "Load(%s) should correctly parse refresh interval", path)
		assert.Equalf(t, "blue", cfg.ColumnBordersColor, "Load(%s) should correctly parse column borders color", path)
		assert.Equalf(t, true, cfg.ShowColumnSeparator, "Load(%s) should correctly parse show column separator", path)
	})

	t.Run("should return error if file not exist", func(t *testing.T) {
		// given
		path := filepath.Join("testdata", "not_exist.yaml")

		// when
		_, err := Load(path)

		// then
		assert.Errorf(t, err, "Load(%s) should return error on not existing file", path)
	})

	t.Run("should return error if file has invalid format", func(t *testing.T) {
		// given
		path := filepath.Join("testdata", "invalid.yaml")

		// when
		_, err := Load(path)

		// then
		assert.Errorf(t, err, "Load(%s) should return error on invalid file format", path)
	})
}
