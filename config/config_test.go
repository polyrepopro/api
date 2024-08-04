package config

import (
	"testing"

	"github.com/polyrepopro/api/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {
	monitoring.Setup()

	type args struct {
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "default test",
			args:    args{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := GetConfig()

			if err != nil && tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.GreaterOrEqual(t, len(cfg.Workspaces), 1)
		})
	}
}
