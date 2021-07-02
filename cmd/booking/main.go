package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/booking/api"
	"github.com/booking/config"
	"github.com/booking/database"
	"github.com/booking/httpd"
	"github.com/booking/logger"
	"github.com/booking/repository"
	"github.com/booking/service"
)

const application = "BookingService"

func main() {
	ctx := context.Background()
	c := config.NewDefaults()

	l := logger.NewLogger(c).WithField("service", application)

	db, err := database.NewPGSQLClient(ctx, c)
	if err != nil {
		l.WithError(err).Error("error creating database")
		return
	}

	server := httpd.NewServer(c, l)

	rr, err := repository.NewRoomRepository(db, c.DBLog)
	if err != nil {
		l.WithError(err).Error("error creating room repository")
		return
	}
	rs := service.NewRoomService(c, rr, l)
	server.Add(api.NewRoomAPI(rs, l).WebService())

	mr, err := repository.NewMeetingRepository(db, c.DBLog)
	if err != nil {
		l.WithError(err).Error("error creating meeting repository")
		return
	}
	ms := service.NewBookingService(c, mr, rr, l)
	server.Add(api.NewBookingAPI(ms, l).WebService())

	server.Start(ctx)

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT)

	for {
		select {
		case <-sigChan:
			l.Info("interrupt signal - exiting")
			server.Stop()
			os.Exit(0)
		case <-server.Shutdown():
			l.Info("service shutting down")
			os.Exit(1)
		}
	}
}
