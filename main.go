package main

import (
	"github.com/qsocket/qs-netcat/config"
	qsnetcat "github.com/qsocket/qs-netcat/pkg"
	"github.com/qsocket/qs-netcat/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	// Configure the options from the flags/config file
	opts, err := config.ConfigureOptions()
	if err != nil {
		logrus.Error(err)
	}
	if opts.RandomSecret {
		opts.Secret = utils.RandomString(20)
	}
	opts.Summarize()

	if opts.Listen {
		qsnetcat.StartProbingQSRN(opts)
		return
	}

	err = qsnetcat.Connect(opts)
	if err != nil {
		logrus.Error(err)
	}
}
