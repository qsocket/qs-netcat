package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	Version  = "?"
	usageStr = `
qs-netcat [-liC] [-e cmd] [-p port]
Version: %s
	-s <secret>  Secret. (e.g. password).
	-l           Listening server. [default: client]
	-g           Generate a Secret. (random)
	-C           Disable encryption.
	-t           Probe interval for QSRN. (5s)
	-T           Use TOR.
	-f <IP>      IPv4 address for port forwarding.
	-p <port>    Port to listen on or forward to.
	-i           Interactive login shell. (TTY) [Ctrl-e q to terminate]
	-e <cmd>     Execute command. [e.g. "bash -il" or "cmd.exe"]
	-pin         Enable certificate pinning on TLS connections.
	-v           Verbose output.
	-q           Quiet. No log output.

Example to forward traffic from port 2222 to 192.168.6.7:22:
  $ qs-netcat -s MyCecret -l -f 192.168.6.7 -p 22     # Server
  $ qs-netcat -s MyCecret -p 2222                     # Client
Example file transfer:
  $ qs-netcat -q -l -s MyCecret >warez.tar.gz         # Server
  $ qs-netcat -q -s MyCecret <warez.tar.gz            # Client
Example for a reverse shell:
  $ qs-netcat -s MyCecret -l -i                       # Server
  $ qs-netcat -s MyCecret -i                          # Client

`
)

// PrintUsageErrorAndDie ...
func PrintUsageErrorAndDie(err error) {
	color.Red("\n%s", err.Error())
	fmt.Printf(usageStr, Version)
	os.Exit(1)
}

// PrintHelpAndDie ...
func PrintHelpAndDie() {
	fmt.Printf(usageStr, Version)
	os.Exit(0)
}

// Main config struct for parsing the TOML file
type Options struct {
	Secret        string
	Execute       string
	ForwardAddr   string
	Port          int
	ProbeInterval int
	DisableTLS    bool
	Interactive   bool
	Listen        bool
	RandomSecret  bool
	CertPinning   bool
	Quiet         bool
	UseTor        bool
	Verbose       bool
	help          bool
}

// ConfigureOptions accepts a flag set and augments it with agentgo-server
// specific flags. On success, an options structure is returned configured
// based on the selected flags.
func ConfigureOptions(fs *flag.FlagSet, args []string) (*Options, error) {
	// Create empty options
	opts := &Options{}

	// Define flags
	fs.BoolVar(&opts.help, "h", false, "Prompt help")
	fs.BoolVar(&opts.help, "help", false, "Prompt help")
	fs.StringVar(&opts.Secret, "s", "", "Secret (e.g. password)")
	fs.StringVar(&opts.Execute, "e", "", "Execute command [e.g. \"bash -il\" or \"cmd.exe\"]")
	fs.StringVar(&opts.ForwardAddr, "f", "", "IPv4 address for port forwarding")
	fs.BoolVar(&opts.Listen, "l", false, "Listening server [default: client]")
	fs.BoolVar(&opts.RandomSecret, "g", false, "Generate a Secret (random)")
	fs.BoolVar(&opts.Interactive, "i", false, "Interactive login shell (TTY) [Ctrl-e q to terminate]")
	fs.IntVar(&opts.Port, "p", 0, "Port to listen on or forward to")
	fs.IntVar(&opts.ProbeInterval, "t", 5, "Probe interval for QSRN")
	fs.BoolVar(&opts.DisableTLS, "C", false, "Disable encryption")
	fs.BoolVar(&opts.UseTor, "T", false, "Use TOR")
	fs.BoolVar(&opts.CertPinning, "pin", false, "Enable certificate pinning on TLS connections")
	fs.BoolVar(&opts.Quiet, "q", false, "Quiet. No log outpu")
	fs.BoolVar(&opts.Verbose, "v", false, "Verbose mode")
	// Parse arguments and check for errors
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if opts.help {
		PrintHelpAndDie()
	}

	if !opts.RandomSecret && opts.Secret == "" {
		color.New(color.FgBlue).Add(color.Bold).Print("[>] ")
		print("Enter Secret (or press Enter to generate): ")
		n, _ := fmt.Scanln(&opts.Secret)
		if n == 0 {
			opts.RandomSecret = true
		}
		// print("\n==============================================\n")
	}

	if opts.ForwardAddr != "" && opts.Port == 0 {
		return nil, errors.New("Please specify a valid port number.")
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
	yellow := color.New(color.FgYellow).Add(color.Bold)
	red := color.New(color.FgRed).Add(color.Bold)
	blue := color.New(color.FgBlue).Add(color.Bold)
	mode := "client"
	if opts.Listen {
		mode = "server"
	}

	yellow.Printf("[#] %s\n", blue.Sprintf(".::Qsocket Netcat::."))
	yellow.Print("├──>")
	fmt.Printf(" Secret: %s\n", red.Sprintf(opts.Secret))
	yellow.Print("├──>")
	fmt.Printf(" Mode: %s\n", mode)
	yellow.Print("├──>")
	fmt.Printf(" TLS: %t\n", !opts.DisableTLS)
	yellow.Print("├──>")
	fmt.Printf(" TOR: %t\n", opts.UseTor)
	if opts.ForwardAddr != "" {
		yellow.Print("├──>")
		fmt.Printf(" Forward: %s\n", opts.ForwardAddr)
	}
	if opts.Port != 0 {
		yellow.Print("├──>")
		fmt.Printf(" Port: %d\n", opts.Port)
	}
	yellow.Print("└──>")
	if opts.Listen {
		fmt.Printf(" Probe Interval: %s\n", time.Second*time.Duration(opts.ProbeInterval))
	} else {
		fmt.Printf(" Probe Duration: %s\n", time.Second*time.Duration(opts.ProbeInterval))
	}
	print("\n")
}
