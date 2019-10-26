package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/it234/goapp/internal/pkg/config"
	models "github.com/it234/goapp/internal/pkg/models/common"
	"github.com/it234/goapp/internal/pkg/models/db"
	"github.com/it234/goapp/internal/pkg/models/sys"
	"github.com/it234/goapp/pkg/hash"
	"github.com/it234/goapp/pkg/random"

	"github.com/jinzhu/gorm"
)

const (
	DEFAULT_ADMIN_USER     = "admin"
	DEFAULT_ADMIN_PASSWORD = "123456"
)

func InitDB(config *config.Config) {
	var gdb *gorm.DB
	var err error
	if config.Gorm.DBType == "mysql" {
		config.Gorm.DSN = config.MySQL.DSN()
	} else if config.Gorm.DBType == "sqlite3" {
		config.Gorm.DSN = config.Sqlite3.DSN()
	}
	gdb, err = gorm.Open(config.Gorm.DBType, config.Gorm.DSN)
	if err != nil {
		panic(err)
	}
	gdb.SingularTable(true)
	if config.Gorm.Debug {
		gdb.LogMode(true)
		gdb.SetLogger(log.New(os.Stdout, "\r\n", 0))
	}
	gdb.DB().SetMaxIdleConns(config.Gorm.MaxIdleConns)
	gdb.DB().SetMaxOpenConns(config.Gorm.MaxOpenConns)
	gdb.DB().SetConnMaxLifetime(time.Duration(config.Gorm.MaxLifetime) * time.Second)
	db.DB = gdb
}

func Migration() {
	fmt.Println(db.DB.AutoMigrate(new(sys.Menu)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.Admins)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.RoleMenu)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.Role)).Error)
	fmt.Println(db.DB.AutoMigrate(new(sys.AdminsRole)).Error)
}

func AdminUserInit() {
	where := sys.Admins{UserName: DEFAULT_ADMIN_USER}
	user := sys.Admins{}
	notFound, err := models.First(&where, &user)
	if notFound || err != nil {
		where.Salt = random.RandString(8)
		where.Status = 1
		where.Password = hash.Md5String(where.Salt + DEFAULT_ADMIN_PASSWORD)
		err = models.Save(&where)
		if err != nil {
			panic(err)
		}
	}
}
