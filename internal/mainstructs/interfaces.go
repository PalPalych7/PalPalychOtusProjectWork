package mainstructs

import "context"

type Server interface {
	Serve() error
	Stop() error
}

type Logger interface {
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
}

type Storage interface {
	Connect(ctx context.Context) error
	AddBannerSlot(ctx context.Context, slotID int, bannerID int) error
	DelBannerSlot(ctx context.Context, slotID int, bannerID int) error
	BannerClick(ctx context.Context, slotID int, bannerID int, socGroupID int) error
	GetBannerForSlot(ctx context.Context, slotID int, socGroupID int) (GetBannerStruct, error)
	GetBannerStat(ctx context.Context) ([]BannerStatStruct, error)
	ChangeSendStatID(ctx context.Context, ID int) error
	Close() error
}

type MyBandit interface {
	GetBannerNum(arrStruct []BannerStruct) int
}

type RabbitQueue interface {
	Start() error
	SendMess(myMes []byte) error
	Shutdown() error
}
