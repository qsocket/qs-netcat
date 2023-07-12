package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fatih/color"
	"github.com/qsocket/qs-netcat/log"
	"github.com/qsocket/qs-netcat/utils"
)

var Version = "?"

const (
	USAGE_EAMPLES = `
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
	DEFAULT_E2E_CIPHER = "SRP-AES-GCM-256-E2E (Prime: 4096)"
)

var (
	ForwardAddrRgx = regexp.MustCompile(`([0-9]{1,5}:|)(?:[0-9]{1,3}\.){3}[0-9]{1,3}:[0-9]{1,5}`)
)

// Main config struct for parsing the TOML file
type Options struct {
	Secret        string `help:"Secret (e.g. password)." name:"secret" short:"s"`
	Execute       string `help:"Execute command [e.g. \"bash -il\" or \"cmd.exe\"]" name:"exec" short:"e"`
	ForwardAddr   string `help:"IP:PORT for traffic forwarding." name:"forward" short:"f"`
	SocksAddr     string `help:"User socks proxy address for connecting QSRN." name:"socks" short:"x"`
	ProbeInterval int    `help:"Probe interval for connecting QSRN." name:"probe" short:"n" default:"5"`
	DisableEnc    bool   `help:"Disable all encryption." name:"plain" short:"C"`
	End2End       bool   `help:"Use E2E encryption. (default:true)" name:"e2e" default:"true"`
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
	print(USAGE_EAMPLES)
	return nil
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

	if opts.ForwardAddr != "" && !ForwardAddrRgx.MatchString(opts.ForwardAddr) {
		return nil, errors.New("invalid forward address")
	}

	if opts.UseTor {
		opts.SocksAddr = "socks5://127.0.0.1:9050"
	}

	if !opts.RandomSecret && opts.Secret == "" {
		color.New(color.FgBlue).Add(color.Bold).Print("[>] ")
		print("Enter Secret (or press Enter to generate): ")
		n, _ := fmt.Scanln(&opts.Secret)
		if n == 0 {
			opts.RandomSecret = true
		}
	}

	if opts.Verbose {
		log.SetLevel(log.LOG_LEVEL_TRACE) // Show all the shit!
	}

	if opts.Quiet {
		log.SetLevel(log.LOG_LEVEL_FATAL) // Show nothing!
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
	fmt.Printf(" TOR: %t\n", opts.UseTor)
	if opts.ForwardAddr != "" {
		yellow.Print(" ├──>")
		fmt.Printf(" Forward: %s\n", opts.ForwardAddr)
	}
	yellow.Print(" ├──>")
	if opts.Listen {
		fmt.Printf(" Probe Interval: %s\n", time.Second*time.Duration(opts.ProbeInterval))
	} else {
		fmt.Printf(" Probe Duration: %s\n", time.Second*time.Duration(opts.ProbeInterval))
	}
	yellow.Print(" └──>")
	if opts.DisableEnc {
		fmt.Printf(" Encryption: %s\n", red.Sprintf("DISABLED"))
	} else {
		if opts.End2End {
			fmt.Printf(" Encryption: %s\n", DEFAULT_E2E_CIPHER)
		} else {
			fmt.Println(" Encryption: TLS (v1.2)")
		}
	}

	print("\n")
}
