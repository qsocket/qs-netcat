package main

import (
	"flag"
	"os"
	"path/filepath"
	"qsutils/config"
	qsutils "qsutils/pkg"
	"qsutils/utils"

	"github.com/sirupsen/logrus"
)

var Version = "v1.0"

func main() {

	// Create a FlagSet and sets the usage
	fs := flag.NewFlagSet(filepath.Base(os.Args[0]), flag.ExitOnError)
	// Configure the options from the flags/config file
	opts, err := config.ConfigureOptions(fs, os.Args[1:])
	if err != nil {
		config.PrintUsageErrorAndDie(err)
	}
	config.SetVersion(Version)
	if opts.RandomSecret {
		opts.Secret = utils.RandomString(20)
	}
	opts.Summarize()

	if opts.Listen {
		qsutils.StartProbingQSRN(opts)
		return
	}

	if opts.Port != 0 {
		qsutils.ServeToLocal(opts)
		return
	}

	err = qsutils.Connect(opts)
	if err != nil {
		logrus.Error(err)
	}
}
