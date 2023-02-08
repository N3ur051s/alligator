package model

import (
	"fmt"
	"time"

	. "alligator/pkg/config"

	"github.com/labstack/echo/v4"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"alligator/pkg/utils/log"
)

var db *gorm.DB

type Model struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func Init(opts Options) {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/alligator?charset=utf8mb4&parseTime=True&loc=Local",
		opts.Db.User, opts.Db.Passwd, opts.Db.Ip, opts.Db.Port)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.NewOrmWriter(),
			logger.Config{
				LogLevel: logger.Error,
			},
		),
	})

	if err != nil {
		log.Error(err)
	}

	sqlDB, err := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(20)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func orderAndPaginate(c echo.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := cast.ToInt(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}
		pageSize := 10
		reqPageSize := c.QueryParam("page_size")
		if reqPageSize != "" {
			pageSize = cast.ToInt(reqPageSize)
		}
		offset := (page - 1) * pageSize

		return db.Order("name").Offset(offset).Limit(pageSize)
	}
}

func totalPage(total int64, pageSize int) int64 {
	n := total / int64(pageSize)
	if total%int64(pageSize) > 0 {
		n++
	}
	return n
}

type Pagination struct {
	Total       int64 `json:"total"`
	PerPage     int   `json:"per_page"`
	CurrentPage int   `json:"current_page"`
	TotalPages  int64 `json:"total_pages"`
}

type DataList struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination,omitempty"`
}

func GetListWithPagination(models interface{},
	c echo.Context, totalRecords int64) (result DataList) {

	page := cast.ToInt(c.QueryParam("page"))
	if page == 0 {
		page = 1
	}

	result = DataList{}

	result.Data = models

	pageSize := 10
	reqPageSize := c.QueryParam("page_size")
	if reqPageSize != "" {
		pageSize = cast.ToInt(reqPageSize)
	}

	result.Pagination = Pagination{
		Total:       totalRecords,
		PerPage:     pageSize,
		CurrentPage: page,
		TotalPages:  totalPage(totalRecords, pageSize),
	}

	return
}
