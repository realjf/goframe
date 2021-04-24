package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/realjf/goframe/config"
	"github.com/realjf/goframe/internal/middleware"
	"github.com/realjf/goframe/pkg/exception"
	"github.com/realjf/goframe/pkg/utils"
)

var (
	DbConn *sql.DB
)

type MySQLDriver struct {
	Username     string
	Password     string
	DbName       string
	Host         string
	Port         string
	Charset      string
	MaxOpenConns int
}

func NewMysql(config config.IConfig) *MySQLDriver {
	configData := config.GetConfigData()
	return &MySQLDriver{
		Username:     configData.Mysql.Username,
		Password:     configData.Mysql.Password,
		Host:         configData.Mysql.Host,
		Port:         utils.ToString(configData.Mysql.Port),
		Charset:      configData.Mysql.Charset,
		DbName:       configData.Mysql.Dbname,
		MaxOpenConns: configData.Mysql.MaxOpenConns,
	}
}

func (db *MySQLDriver) Init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&timeout=10s",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.DbName,
		db.Charset)

	var err error
	DbConn, err = sql.Open("mysql", dataSourceName)
	exception.CheckError(err, 3000)
	DbConn.SetMaxIdleConns(2)
	DbConn.SetMaxOpenConns(db.MaxOpenConns)
	DbConn.SetConnMaxLifetime(time.Second * time.Duration(60))

	middleware.Logger.Logger.Info("init db ...")
}
