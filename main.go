package main

import (
	"github.com/qsocket/qs-netcat/config"
	"github.com/qsocket/qs-netcat/log"
	qsnetcat "github.com/qsocket/qs-netcat/pkg"
	"github.com/qsocket/qs-netcat/utils"
)

func main() {
	// Configure the options from the flags/config file
	opts, err := config.ConfigureOptions()
	if err != nil {
		log.Fatal(err)
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
		log.Error(err)
	}
}
