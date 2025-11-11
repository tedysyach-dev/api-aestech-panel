package config

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDb membuat koneksi MongoDB dengan konfigurasi pool, timeout, dan logging
func NewMongoDb(viper *viper.Viper, logrusLogger *logrus.Logger) *mongo.Database {
	databaseUrl := viper.GetString("mongo.url")
	username := viper.GetString("mongo.username")
	password := viper.GetString("mongo.password")
	database := fmt.Sprintf("%s_%s", viper.GetString("database.name"), viper.GetString("app.development"))
	idleConnection := viper.GetInt("mongo.pool.idle")
	maxConnection := viper.GetInt("mongo.pool.max")
	connectTimeout := viper.GetInt("mongo.pool.lifetime") // dalam detik

	clientOptions := options.Client().ApplyURI(databaseUrl)

	// Auth
	clientOptions.SetAuth(options.Credential{
		Username:   username,
		Password:   password,
		AuthSource: "admin",
	})

	// Pool configuration
	clientOptions.SetMaxPoolSize(uint64(maxConnection))
	clientOptions.SetMinPoolSize(uint64(idleConnection))

	// Monitoring command agar logrus bisa menangkap perintah ke MongoDB
	clientOptions.Monitor = &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			logrusLogger.WithFields(logrus.Fields{
				"command":  evt.CommandName,
				"database": evt.DatabaseName,
			}).Tracef("MongoDB command started: %v", evt.Command)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			logrusLogger.WithField("command", evt.CommandName).
				Tracef("MongoDB command succeeded in %v", evt.Duration)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			logrusLogger.WithField("command", evt.CommandName).
				Errorf("MongoDB command failed after %v: %v", evt.Duration, evt.Failure)
		},
	}

	// Context dengan timeout (default 10 detik)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(connectTimeout)*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrusLogger.Fatalf("failed to connect MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		logrusLogger.Fatalf("failed to ping MongoDB: %v", err)
	}

	return client.Database(database)
}
