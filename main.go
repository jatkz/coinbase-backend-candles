package main

import (
	"github.com/natefinch/lumberjack"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	// Create a new logger with log rotation
	logger := &lumberjack.Logger{
		Filename:   "/var/log/coinbase-etl.log",
		MaxSize:    100, // megabytes
		MaxAge:     7,   // days
		LocalTime:  true,
		Compress:   true,
		MaxBackups: 5,
	}

    // Configure logrus to write to the log file
    log := logrus.New()
    log.SetFormatter(&logrus.JSONFormatter{})
    log.SetOutput(logger)

    // Log a message using logrus
    log.Info("Hello from Go!")

	c := cron.New()

    // Define a cron job that runs every minute
    c.AddFunc("* * * * *", func() {
        log.Info("Running cron job 1...")
    })

    // Define a cron job that runs every hour
    c.AddFunc("0 * * * *", func() {
        log.Info("Running cron job 2...")
    })

    // Start the cron scheduler
    c.Start()

}
