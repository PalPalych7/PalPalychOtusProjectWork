package sqlstorage

import (
	"context"
	"database/sql"

	ms "github.com/PalPalych7/OtusProjectWork/internal/mainstructs"
	_ "github.com/jackc/pgx/stdlib" // justifying
	_ "github.com/lib/pq"
)

type Storage struct {
	DBConf    ms.DBConf
	DBConnect *sql.DB
	MyBandit  ms.MyBandit
}

func New(myDBConf ms.DBConf, myBandit ms.MyBandit) *Storage {
	return &Storage{
		DBConf: myDBConf, MyBandit: myBandit,
	}
}

func rowsToStruct(rows *sql.Rows) ([]ms.BannerStruct, error) {
	var myBannerList []ms.BannerStruct
	var bannerID, ShowCount, ClickCount int
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&bannerID, &ShowCount, &ClickCount); err != nil {
			return nil, err
		}
		myBanner := ms.BannerStruct{
			BannerID:   bannerID,
			ShowCount:  ShowCount,
			ClickCount: ClickCount,
		}
		myBannerList = append(myBannerList, myBanner)
	}
	return myBannerList, nil
}

func rowsToStat(rows *sql.Rows) ([]ms.BannerStatStruct, error) {
	var myBannerList []ms.BannerStatStruct
	var id, slotID, bannerID, socGroupID int
	var statType, recDate string
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &slotID, &bannerID, &socGroupID, &statType, &recDate); err != nil {
			return nil, err
		}
		myBannerList = append(myBannerList, ms.BannerStatStruct{
			ID:         id,
			SlotID:     slotID,
			BannerID:   bannerID,
			SocGroupID: socGroupID,
			StatType:   statType,
			RecDate:    recDate,
		})
	}
	return myBannerList, nil
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	myStr := "postgres://" + s.DBConf.DBUserName + ":" + s.DBConf.DBPassward + "@"
	myStr += s.DBConf.DBHost + ":" + s.DBConf.DBPort + "/" + s.DBConf.DBName + "?sslmode=disable"
	s.DBConnect, err = sql.Open("postgres", myStr)
	if err == nil {
		err = s.DBConnect.PingContext(ctx)
	}
	return err
}

func (s *Storage) AddBannerSlot(ctx context.Context, slotID int, bannerID int) error {
	query := `
			insert into slot_banner(slot_id,  banner_id)
			values($1, $2)
		`
	_, err := s.DBConnect.ExecContext(ctx, query, slotID, bannerID)
	return err
}

func (s *Storage) DelBannerSlot(ctx context.Context, slotID int, bannerID int) error {
	query := `
			delete from slot_banner
			where slot_id = $1 
			  and  banner_id=$2
		`
	_, err := s.DBConnect.ExecContext(ctx, query, slotID, bannerID)
	return err
}

func (s *Storage) BannerClick(ctx context.Context, slotID int, bannerID int, socGroupID int) error {
	query := `
			insert into banner_stat(slot_id,  banner_id, soc_group_id, stat_type)
			values($1, $2, $3, 'C')
		`
	_, err := s.DBConnect.ExecContext(ctx, query, slotID, bannerID, socGroupID)
	return err
}

func (s *Storage) GetBannerForSlot(ctx context.Context, slotID int, socGroupID int) (ms.GetBannerStruct, error) {
	queryStat := `
		select  sb.banner_id, count(distinct bs_s.id) show_count, count(distinct bs_c.id) click_count
		from slot_banner sb
		left join banner_stat bs_s
			on sb.slot_id=bs_s.slot_id
			and sb.banner_id=bs_s.banner_id
			and bs_s.soc_group_id=$1
			and bs_s.stat_type = 'S'
		left join banner_stat bs_c
			on sb.slot_id=bs_c.slot_id
			and sb.banner_id=bs_c.banner_id
			and bs_c.soc_group_id=$1
			and bs_c.stat_type = 'C'
		where sb.slot_id=$2 
		group by sb.banner_id;
	`
	myStat, errStat := s.DBConnect.QueryContext(ctx, queryStat, socGroupID, slotID)
	if errStat != nil {
		return ms.GetBannerStruct{}, errStat
	}

	myBannerList, errStruct := rowsToStruct(myStat)
	if errStruct != nil {
		return ms.GetBannerStruct{}, errStruct
	}

	if len(myBannerList) == 0 {
		return ms.GetBannerStruct{}, nil
	}

	arrNum := s.MyBandit.GetBannerNum(myBannerList)
	myBannerID := myBannerList[arrNum].BannerID
	query := `
					insert into banner_stat(slot_id,  banner_id, soc_group_id, stat_type)
					values($1, $2, $3, 'S')
	`
	_, err := s.DBConnect.ExecContext(ctx, query, slotID, myBannerID, socGroupID)
	return ms.GetBannerStruct{ID: myBannerID}, err
}

func (s *Storage) GetBannerStat(ctx context.Context) ([]ms.BannerStatStruct, error) {
	queryStat := `
		select id, slot_id, banner_id, soc_group_id, stat_type, rec_date
		from banner_stat
		where id>(select max(banner_stat_id) from send_stat_max_id)	
	`
	myStat, errStat := s.DBConnect.QueryContext(ctx, queryStat)
	if errStat != nil {
		return nil, errStat
	}

	myBannerStatList, errStruct := rowsToStat(myStat)
	if errStruct != nil {
		return nil, errStruct
	}
	return myBannerStatList, nil
}

func (s *Storage) ChangeSendStatID(ctx context.Context, id int) error {
	query := `
			update send_stat_max_id
			set banner_stat_id = $1
		`
	_, err := s.DBConnect.ExecContext(ctx, query, id)
	return err
}

func (s *Storage) Close() error {
	err := s.DBConnect.Close()
	return err
}
