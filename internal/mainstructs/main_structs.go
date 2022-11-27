package mainstructs

type BannerStatStruct struct {
	ID         int
	SlotID     int
	BannerID   int
	SocGroupID int
	StatType   string
	RecDate    string
}

type GetBannerStruct struct {
	ID int
}

type BannerStruct struct {
	BannerID   int
	ShowCount  int
	ClickCount int
}

type BanditConfig struct {
	FullLearnigCount     int // количество запросов в режиме "полного обучения"
	PartialLearningCount int // количество запросов в режиме "чаcтичного обучения"
	FinalRandomPecent    int // вероятность случайного выбора после обучения (в процентах)
}

type LoggerConf struct {
	LogFile string
	Level   string
}

type HTTPConf struct {
	Host              string
	Port              string
	TimeOutSec        int
	ReadHeaderTimeout int
}

type DBConf struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUserName string
	DBPassward string
}

type MainConfig struct {
	Logger LoggerConf
	HTTP   HTTPConf
	DB     DBConf
	Bandit BanditConfig
}

type RabbitCFG struct {
	URI          string
	Exchange     string
	ExchangeType string
	Queue        string
	BindingKey   string
	ConsumerTag  string
	SleepSecond  int
	TimeOutSec   int
}

type StatSenderConfig struct {
	Logger LoggerConf
	Rabbit RabbitCFG
	DB     DBConf
}

type SlotBanner struct {
	SlotID   int
	BannerID int
}

type ForBannerClick struct {
	SlotID     int
	BannerID   int
	SocGroupID int
}

type ForGetBanner struct {
	SlotID     int
	SocGroupID int
}
