package config

import (
	"fmt"
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viper *viper.Viper, logrusLogger *logrus.Logger) *gorm.DB {
	username := viper.GetString("database.username")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	database := viper.GetString("database.name")
	idleConnection := viper.GetInt("database.pool.idle")
	maxConnection := viper.GetInt("database.pool.max")
	maxLifeTimeConnection := viper.GetInt("database.pool.lifetime")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)

	// Buat writer untuk GORM logger — arahkan output ke logrus
	gormWriter := log.New(logrusWriter{Logger: logrusLogger}, "", 0)

	// Setup GORM logger agar menampilkan query lengkap + value
	gormLogger := logger.New(
		gormWriter,
		logger.Config{
			SlowThreshold:             2 * time.Second, // query lambat
			LogLevel:                  logger.Info,     // tampilkan semua query
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false, // ❗ ubah ke false agar value ditampilkan
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		logrusLogger.Fatalf("failed to connect database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrusLogger.Fatalf("failed to connect database: %v", err)
	}

	sqlDB.SetMaxIdleConns(idleConnection)
	sqlDB.SetMaxOpenConns(maxConnection)
	sqlDB.SetConnMaxLifetime(time.Second * time.Duration(maxLifeTimeConnection))

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

// log.New() membutuhkan interface io.Writer, jadi kita implementasikan Write()
func (l logrusWriter) Write(p []byte) (n int, err error) {
	l.Logger.Trace(string(p))
	return len(p), nil
}
