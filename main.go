package main

import (
	"io"
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
	if err != nil || opts == nil {
		log.Fatal(err)
		return
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
				switch err {
				case qsocket.ErrPeerNotFound, io.EOF:
					log.Debug(err)
				case qsocket.ErrServerCollision:
					log.Fatal(err)
				default:
					log.Error(err)
				}
			}
		}

	}

	err = qsnetcat.Connect(opts)
	if err != nil && err != io.EOF {
		log.Error(err)
	}
}
