package storage

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestNewStorage(t *testing.T) {
	ctrl := gomock.NewController(t)

	testData := map[string]struct {
		wal    WAL
		eng    Engine
		logger *zap.Logger

		err error
	}{
		"create storage without engine": {
			logger: zap.NewNop(),
			err:    errors.New("engine required"),
		},
		"create storage without logger": {
			err: errors.New("logger required"),
		},
		"create storage": {
			wal:    nil,
			eng:    NewMockEngine(ctrl),
			logger: zap.NewNop(),
			err:    nil,
		},
	}

	for name, data := range testData {
		t.Run(name, func(t *testing.T) {
			storage, err := NewStorage(data.logger, data.eng, data.wal)
			assert.Equal(t, data.err, err)
			if data.err != nil {
				assert.Nil(t, storage)
			} else {
				assert.NotNil(t, storage)
			}
		})
	}
}
