package main

import (
	"errors"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/jim-minter/mitm-proxy/pkg/config"
	"github.com/jim-minter/mitm-proxy/pkg/linux"
	"github.com/jim-minter/mitm-proxy/pkg/proxy"
)

func checkRoot() error {
	if os.Getuid() != 0 {
		return errors.New("must run as root")
	}

	return nil
}

func run(log *logrus.Entry) error {
	if err := checkRoot(); err != nil {
		return err
	}

	if err := config.Read(log); err != nil {
		return err
	}

	p, err := proxy.NewProxy(log)
	if err != nil {
		return err
	}

	go p.Run()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	if err := linux.EnableIPForwarding(log); err != nil {
		return err
	}

	if err := linux.EnableIPTables(log); err != nil {
		return err
	}

	<-ch

	if err := linux.DisableIPTables(log); err != nil {
		return err
	}

	return nil
}

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	log := logrus.NewEntry(logrus.StandardLogger())

	if err := run(log); err != nil {
		log.Fatal(err)
	}
}
