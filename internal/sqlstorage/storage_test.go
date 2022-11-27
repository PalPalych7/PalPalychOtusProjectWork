package sqlstorage

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
		var getBannerStruct ms.GetBannerStruct
		var bannerStuct []ms.BannerStatStruct
		ctx := context.Background()

		m := new(mocks.Storage)

		m.On("Connect", ctx).Return(nil)
		err = m.Connect(ctx)
		require.NoError(t, err)
		m.AssertExpectations(t)

		m.On("AddBannerSlot", ctx, 1, 1).Return(nil)
		err = m.AddBannerSlot(ctx, 1, 1)
		require.NoError(t, err)
		m.AssertExpectations(t)

		m.On("DelBannerSlot", ctx, 1, 1).Return(nil)
		err = m.DelBannerSlot(ctx, 1, 1)
		require.NoError(t, err)
		m.AssertExpectations(t)

		m.On("BannerClick", ctx, 1, 1, 1).Return(nil)
		err = m.BannerClick(ctx, 1, 1, 1)
		require.NoError(t, err)
		m.AssertExpectations(t)

		m.On("GetBannerForSlot", ctx, 1, 1).Return(ms.GetBannerStruct{}, nil)
		getBannerStruct, err = m.GetBannerForSlot(ctx, 1, 1)
		require.NoError(t, err)
		require.Equal(t, ms.GetBannerStruct{}, getBannerStruct)
		m.AssertExpectations(t)

		m.On("GetBannerStat", ctx).Return(nil, nil)
		bannerStuct, err = m.GetBannerStat(ctx)
		require.NoError(t, err)
		require.Empty(t, bannerStuct)
		m.AssertExpectations(t)

		m.On("ChangeSendStatID", ctx, 1).Return(nil)
		err = m.ChangeSendStatID(ctx, 1)
		require.NoError(t, err)
		m.AssertExpectations(t)

		m.On("Close").Return(nil)
		err = m.Close()
		require.NoError(t, err)
		m.AssertExpectations(t)
	})
}
