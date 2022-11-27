package integration_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	_ "github.com/jackc/pgx/stdlib" // justifying
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
)

type mySuite struct {
	suite.Suite
	ctx       context.Context
	client    http.Client
	hostName  string
	DBConnect *sql.DB
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

type GetBannerStruct struct {
	ID int
}

var (
	err               error
	bodyRaw           []byte
	req               *http.Request
	resp              *http.Response
	countRec          int
	myGetBannerStruct GetBannerStruct
)

func (s *mySuite) CheckCountRec(myQueryText string, expCount int) {
	mySQLRows, err := s.DBConnect.QueryContext(s.ctx, myQueryText)
	s.Require().NoError(err)
	defer mySQLRows.Close()
	mySQLRows.Next()
	err = mySQLRows.Scan(&countRec)
	fmt.Println("countRec=", countRec)
	s.Require().NoError(err)
	s.Require().Equal(expCount, countRec)
}

func (s *mySuite) SetupSuite() {
	s.client = http.Client{
		Timeout: time.Second * 5,
	}

	s.hostName = "http://mainSevice:5000/"
	s.ctx = context.Background()
	myStr := "postgres://otusfinalproj:otusfinalproj@postgres_db:5432/otusfinalproj?sslmode=disable" // через докер
	// myStr := "postgres://otusfinalproj:otusfinalproj@localhost:5432/otusfinalproj?sslmode=disable" // локально

	s.DBConnect, err = sql.Open("postgres", myStr)
	if err == nil {
		err = s.DBConnect.PingContext(s.ctx)
	}
	s.Require().NoError(err)

	s.CheckCountRec("select count(*) RC from banner", 20)
	_, err = s.DBConnect.ExecContext(s.ctx, "delete from slot_banner")
	s.Require().NoError(err)
	s.CheckCountRec("select count(*) RC from slot_banner", 0)

	_, err = s.DBConnect.ExecContext(s.ctx, "delete from banner_stat")
	s.Require().NoError(err)
	_, err = s.DBConnect.ExecContext(s.ctx, "update send_stat_max_id set banner_stat_id=0")
	s.Require().NoError(err)
	s.CheckCountRec("select count(*) RC from banner_stat", 0)
}

func (s *mySuite) TearDownSuite() {
	mySQL := `
		delete from slot_banner;
		delete from banner_stat;
		update send_stat_max_id set banner_stat_id=0;
	`
	_, err = s.DBConnect.ExecContext(s.ctx, mySQL)
	s.Require().NoError(err)
	s.DBConnect.Close()
}

func (s *mySuite) SendRequest(myMethodName string, myAnyStruct interface{}) []byte {
	bodyRaw, err = json.Marshal(myAnyStruct)
	s.Require().NoError(err)
	req, err = http.NewRequestWithContext(s.ctx, http.MethodPost, s.hostName+myMethodName, bytes.NewBuffer(bodyRaw))
	s.Require().NoError(err)

	resp, err = s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()
	bodyRaw, err = ioutil.ReadAll(resp.Body)
	s.Require().NoError(err)
	return bodyRaw
}

func (s *mySuite) AddSlotBanner(mySlotBanner SlotBanner) {
	// добавление баннера к слоту
	bodyRaw = s.SendRequest("AddBannerSlot", mySlotBanner)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) DelSlotBanner(mySlotBanner SlotBanner) { // удалени баннера из слота
	bodyRaw = s.SendRequest("DelBannerSlot", mySlotBanner)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) GetBannerForSlot(mySlotSoc ForGetBanner) GetBannerStruct { // получения баннера для показа в слоте
	bodyRaw = s.SendRequest("GetBannerForSlot", mySlotSoc)
	s.Require().NotEmpty(bodyRaw)
	err = json.Unmarshal(bodyRaw, &myGetBannerStruct)
	s.Require().NoError(err)
	return myGetBannerStruct
}

func (s *mySuite) BannerClick(myBannerClick ForBannerClick) { // клик по баннеру
	bodyRaw = s.SendRequest("BannerClick", myBannerClick)
	s.Require().Empty(bodyRaw)
}

func (s *mySuite) Test1AddBanner() {
	for i := 1; i <= 10; i++ {
		s.AddSlotBanner(SlotBanner{1, i})
	}
	// к слоту привязано 10 баннеров
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1", 10)
	// привязан баннер с id=1
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1 and banner_id=1", 1)
	s.AddSlotBanner(SlotBanner{1, 1})
	// после повторной попытке ничего не изменилось
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=1 and banner_id=1", 1)
	fmt.Println("finish Test1AddBanner")
}

func (s *mySuite) Test2DelBanner() {
	fmt.Println("start TestDelSlotBanner")
	//  добавим баннер к слоту
	s.AddSlotBanner(SlotBanner{2, 2})
	// убедимся что он добавился
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2 and banner_id=2", 1)
	// отвяжем баннер от слота
	s.DelSlotBanner(SlotBanner{2, 2})
	// убедимся что отвязался
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2 and banner_id=2", 0)
	fmt.Println("finish TestDelSlotBanner")
}

func (s *mySuite) Test3GetBannerForSlot() {
	// убедимся, что к слоту 2 не првязан ни один баннер
	s.CheckCountRec("select count(*) RC from slot_banner where slot_id=2", 0)
	// поскольку к слоту 2 не првязан ни один баннер должен вернуть 0
	myGetBannerStruct = s.GetBannerForSlot(ForGetBanner{2, 1})
	s.Require().Equal(GetBannerStruct{}, myGetBannerStruct)

	// добавим во второй слот баннер с ID=3
	s.AddSlotBanner(SlotBanner{2, 3})
	// теперь должен вернуть ID=3 (так как это единственный баннер
	myGetBannerStruct = s.GetBannerForSlot(ForGetBanner{2, 1})
	s.Require().Equal(3, myGetBannerStruct.ID)
	// убедимся что этот показ отразился в статистике (1 раз)
	s.CheckCountRec("select count(*) RC from banner_stat where stat_type='S' and slot_id=2 and banner_id=3", 1)
}

func (s *mySuite) Test4BannerClick() {
	// убедимся, что к в слоте 1 для баннера 2 для соц группы 3 ещё не было кликов
	mySQL := `
	  select count(*) RC 
	  from banner_stat 
	  where stat_type='C' 
		and slot_id=1 
		and banner_id=2 
		and soc_group_id=3`
	s.CheckCountRec(mySQL, 0)
	//  кликнем в слоте 1 на баннер 2 для соц группы 3
	s.BannerClick(ForBannerClick{1, 2, 3})
	// убедимся, что теперь сохранился 1 клик
	s.CheckCountRec(mySQL, 1)
}

func (s *mySuite) Test5SendMessages() {
	mySQL := `
	  select count(*) RC 
      from banner_stat 
	  where id>(
		select max(banner_stat_id) 
		from send_stat_max_id
	)`
	// убедимся, что есть неотправленные сообщения (2)
	s.CheckCountRec(mySQL, 2)

	// уснём, чтобы дождаться отправки статистики
	time.Sleep(time.Second * 30)

	// убедимся, что теперь не осталоьс не отправленных сообщений
	s.CheckCountRec(mySQL, 0)
}

func TestService(t *testing.T) {
	suite.Run(t, new(mySuite))
}
