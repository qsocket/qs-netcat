package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/qsocket/qs-netcat/config"
	qsnetcat "github.com/qsocket/qs-netcat/pkg"
	"github.com/qsocket/qs-netcat/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	// Configure the options from the flags/config file
	opts, err := config.ConfigureOptions(fs, os.Args[1:])
	if err != nil {
		config.PrintUsageErrorAndDie(err)
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
