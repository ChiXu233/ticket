package database

import (
	"github.com/gin-gonic/gin"
	log "github.com/wonderivan/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	config "ticket-service/conf"
	"ticket-service/database/model"
)

var ormDB Database

type OrmDB struct {
	*gorm.DB
}

type Database interface {
	GetEntityByID(table string, id int, entity interface{}) error
	GetEntityForUpdate(table string, id int, entity interface{}) error
	AssertRowExist(table string, filter map[string]interface{}, params model.QueryParams, entity interface{}) error
	ListEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, entities interface{}) error
	CountEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, count *int64) error
	CountAllEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, count *int64) error
	GetEntityPluck(table string, filter map[string]interface{}, params model.QueryParams, column string, cols interface{}) error
	CreateEntity(table string, entity interface{}) error
	BatchCreateEntity(table string, entities interface{}) error
	SaveEntity(table string, updater interface{}) error
	UpdateEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, updater interface{}) error
	DeleteEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, mode interface{}) error

	DeleteEntity(mode interface{}) error

	DeleteUnscopedEntityByFilter(table string, filter map[string]interface{}, params model.QueryParams, mode interface{}) error

	Begin() (Database, error)
	Commit() error
	Rollback() error
}

func InitDB() error {
	//const DSN = "gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local"
	DBConfig := config.Conf.DB
	DSN := DBConfig.User + ":" + DBConfig.Password + "@tcp(" + DBConfig.Host + ":" + DBConfig.Port + ")/" + DBConfig.Name + "?charset=utf8&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       DSN,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})

	if err != nil {
		return err
	}

	if config.Conf.DB.InitTable {
		initTable(db)
	}

	switch config.Conf.APP.Mode {
	case gin.ReleaseMode:
		db.Logger = db.Config.Logger.LogMode(logger.Error)
	case gin.TestMode:
		db.Logger = db.Config.Logger.LogMode(logger.Warn)
	case gin.DebugMode:
		db.Logger = db.Config.Logger.LogMode(logger.Info)
	}
	ormDB = &OrmDB{
		DB: db,
	}
	return nil
}

func GetDatabase() Database {
	return ormDB
}

// SetMockDatabase for unit test
func SetMockDatabase(mockDB Database) {
	ormDB = mockDB
}

func initTable(db *gorm.DB) {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		log.Error("init table[%s] error.[%s]", model.TableNameUser, err.Error())
	}

}

func (db *OrmDB) Begin() (Database, error) {
	tx := db.DB.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	return &OrmDB{DB: tx}, nil
}

func (db *OrmDB) Commit() error {
	tx := db.DB.Commit()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

func (ormdb *OrmDB) Rollback() error {
	tx := ormdb.DB.Rollback()
	if err := tx.Error; err != nil {
		return err
	}
	return nil
}

// NewSqliteDatabase for unit test
//func NewSqliteDatabase() (*gorm.DB, error) {
//	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
//		NowFunc: func() time.Time {
//			return time.Now().UTC()
//		},
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//	})
//	if err != nil {
//		log.Error("failed to connect sqlite database")
//		return nil, err
//	}
//	initTable(db)
//	return db, nil
//}

// //NewPostgresDatabase for unit test
//func NewPostgresDatabase(host, user, password, dbName string, port int) (*gorm.DB, error) {
//	sqlConnection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s TimeZone=Asia/Shanghai",
//		host, port, user, dbName, password)
//	db, err := gorm.Open(postgres.Open(sqlConnection), &gorm.Config{
//		DisableForeignKeyConstraintWhenMigrating: true,
//		IgnoreRelationshipsWhenMigrating:         true,
//		NamingStrategy: schema.NamingStrategy{
//			SingularTable: true,
//		},
//		Logger: logger.Default.LogMode(logger.Silent),
//	})
//	if err != nil {
//		log.Error("failed to connect postgres database")
//		return nil, err
//	}
//	return db, nil
//}
