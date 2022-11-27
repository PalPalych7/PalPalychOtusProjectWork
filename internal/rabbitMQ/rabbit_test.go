package rabbitmq

import (
	"context"
	"testing"

	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
	"github.com/PalPalych7/OtusProjectWork/mocks"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("main", func(t *testing.T) {
		var err error

		_, err = New(context.Background(), ms.RabbitCFG{URI: "amqp://guest:guest@rabbitmq:5672/"})
		require.NoError(t, err)

		m := new(mocks.RabbitQueue)

		m.On("Start").Return(nil)
		err = m.Start()
		require.NoError(t, err)
		m.AssertExpectations(t)

		m.On("SendMess", []byte("loko")).Return(nil)
		err = m.SendMess([]byte("loko"))
		require.NoError(t, err)
		m.AssertExpectations(t)

		m.On("Shutdown").Return(nil)
		err = m.Shutdown()
		require.NoError(t, err)
		m.AssertExpectations(t)
	})
}
