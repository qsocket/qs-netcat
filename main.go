package main

import (
	"os"
	"time"

	"github.com/qsocket/qs-netcat/config"
	"github.com/qsocket/qs-netcat/log"
	qsnetcat "github.com/qsocket/qs-netcat/pkg"
	"github.com/qsocket/qs-netcat/utils"
	"github.com/qsocket/qsocket-go"
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
		go utils.WaitForExitSignal(os.Interrupt)
		firstRun := true
		log.Info("Listening for connections...")
		for {
			if !firstRun {
				time.Sleep(time.Duration(opts.ProbeInterval) * time.Second)
			} else {
				firstRun = false
			}
			err := qsnetcat.ProbeQSRN(opts)
			if err != nil {
				if err == qsocket.ErrAddressInUse {
					log.Fatal(err)
				}
				if err != qsocket.ErrConnRefused {
					log.Error(err)
				}
			}
		}

	}

	err = qsnetcat.Connect(opts)
	if err != nil {
		log.Error(err)
	}
}
