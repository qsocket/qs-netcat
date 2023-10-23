package utils

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"

	"github.com/fatih/color"
)

func PrintFatal(format string, a ...any) {
	yellow := color.New(color.FgRed).Add(color.Bold)
	yellow.Print("[!] ")
	fmt.Printf(format, a...)
}

func PrintStatus(format string, a ...any) {
	yellow := color.New(color.FgYellow).Add(color.Bold)
	yellow.Print("[*] ")
	fmt.Printf(format, a...)
}

func CaclChecksum(data []byte, base uint) uint {
	checksum := uint(0)
	for _, n := range data {
		checksum += uint(n)
	}
	return checksum % base
}

func RandomString(n int) string {
	// rand.Seed(time.Now().UTC().UnixMicro())
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func WaitForExitSignal(sig os.Signal) {
	for {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, sig)
		<-sigChan
		print("\n")
		PrintFatal("Exiting...\n")
		os.Exit(0)
	}
}
