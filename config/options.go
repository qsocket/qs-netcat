package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/qsocket/qs-netcat/utils"
	"github.com/sirupsen/logrus"
)

var Version = "?"

const (
	UsageExamples = `
Example to forward traffic from port 2222 to 192.168.6.7:22:
	$ qs-netcat -s MyCecret -l -f 192.168.6.7:22        # Server
	$ qs-netcat -s MyCecret -f :2222                    # Client
Example file transfer:
	$ qs-netcat -q -l -s MyCecret >warez.tar.gz         # Server
	$ qs-netcat -q -s MyCecret <warez.tar.gz            # Client
Example for a reverse shell:
	$ qs-netcat -s MyCecret -l -i                       # Server
	$ qs-netcat -s MyCecret -i                          # Client
`
)

// Main config struct for parsing the TOML file
type Options struct {
	UUID          string `help:"UUID form of the qsocket secret." name:"uuid" hidden:""`
	Secret        string `help:"Secret (e.g. password)." name:"secret" short:"s"`
	Execute       string `help:"Execute command [e.g. \"bash -il\" or \"cmd.exe\"]" name:"exec" short:"e"`
	ForwardAddr   string `help:"IP:PORT for traffic forwarding." name:"forward" short:"f"`
	ProbeInterval int    `help:"Probe interval for connecting QSRN." name:"probe" short:"n" default:"5"`
	DisableTLS    bool   `help:"Disable TLS encryption." name:"no-tls" short:"C"`
	Interactive   bool   `help:"Execute with a PTY shell." name:"interactive" short:"i"`
	Listen        bool   `help:"Server mode. (listen for connections)" name:"listen" short:"l"`
	RandomSecret  bool   `help:"Generate a Secret. (random)" name:"generate" short:"g"`
	CertPinning   bool   `help:"Enable certificate pinning on TLS connections." name:"pin" short:"K"`
	Quiet         bool   `help:"Quiet mode. (no stdout)" name:"quiet" short:"q"`
	UseTor        bool   `help:"Use TOR for connecting QSRN." name:"tor" short:"T"`
	Verbose       bool   `help:"Verbose mode." name:"verbose" short:"v"`
	Version       kong.VersionFlag
}

func HelpPrompt(options kong.HelpOptions, ctx *kong.Context) error {
	err := kong.DefaultHelpPrinter(options, ctx)
	if err != nil {
		return err
	}

	_, err = ctx.Stdout.Write([]byte(UsageExamples))
	return err
}

// ConfigureOptions accepts a flag set and augments it with agentgo-server
// specific flags. On success, an options structure is returned configured
// based on the selected flags.
func ConfigureOptions() (*Options, error) {

	// If QS_ARGS exists overwrite the given arguments.
	qsArgs := os.Getenv("QS_ARGS")
	args := os.Args[1:]
	if qsArgs != "" {
		args = strings.Split(qsArgs, " ")
	}

	// Parse arguments and check for errors
	opts := &Options{}

	parser, err := kong.New(
		opts,
		kong.Help(HelpPrompt),
		kong.UsageOnError(),
		kong.Vars{"version": Version},
		kong.ConfigureHelp(kong.HelpOptions{
			Summary: true,
		}),
	)
	if err != nil {
		return nil, err
	}

	_, err = parser.Parse(args)
	if err != nil {
		return nil, err
	}

	// Generate random secret
	if !opts.Listen && opts.RandomSecret {
		print(utils.RandomString(20))
		os.Exit(0)
	}

	if !opts.RandomSecret && (opts.Secret == "" && opts.UUID == "") {
		color.New(color.FgBlue).Add(color.Bold).Print("[>] ")
		print("Enter Secret (or press Enter to generate): ")
		n, _ := fmt.Scanln(&opts.Secret)
		if n == 0 {
			opts.RandomSecret = true
		}
	}

	if opts.Verbose {
		logrus.SetLevel(logrus.TraceLevel) // Show all the shit!
	}

	if opts.Quiet {
		logrus.SetLevel(logrus.FatalLevel) // Show nothing!
	}

	return opts, nil
}

func (opts *Options) Summarize() {
	if opts.Quiet {
		return
	}
	yellow := color.New(color.FgYellow)
	byellow := color.New(color.FgYellow).Add(color.Bold)
	red := color.New(color.FgRed)
	blue := color.New(color.FgBlue)
	mode := "client"
	if opts.Listen {
		mode = "server"
	}

	byellow.Printf("[#] %s\n", blue.Sprintf(".::Qsocket Netcat::."))
	yellow.Print(" ├──>")
	fmt.Printf(" Secret: %s\n", red.Sprintf(opts.Secret))
	yellow.Print(" ├──>")
	fmt.Printf(" Mode: %s\n", mode)
	yellow.Print(" ├──>")
	fmt.Printf(" TLS: %t\n", !opts.DisableTLS)
	yellow.Print(" ├──>")
	fmt.Printf(" TOR: %t\n", opts.UseTor)
	if opts.ForwardAddr != "" {
		yellow.Print(" ├──>")
		fmt.Printf(" Forward: %s\n", opts.ForwardAddr)
	}
	yellow.Print(" └──>")
	if opts.Listen {
		fmt.Printf(" Probe Interval: %s\n", time.Second*time.Duration(opts.ProbeInterval))
	} else {
		fmt.Printf(" Probe Duration: %s\n", time.Second*time.Duration(opts.ProbeInterval))
	}
	print("\n")
}
