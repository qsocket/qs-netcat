module github.com/qsocket/qs-netcat

go 1.19

require (
	github.com/alecthomas/kong v0.7.1
	github.com/briandowns/spinner v1.23.0
	github.com/creack/pty v1.1.18
	github.com/fatih/color v1.15.0
	github.com/mdp/qrterminal/v3 v3.1.1
	github.com/qsocket/conpty-go v0.0.0-20230315180542-d8f8596877dc
	github.com/qsocket/qsocket-go v0.0.4-beta.0.20231023185058-fb49646f6e34
	golang.org/x/term v0.13.0
)

require (
	github.com/google/uuid v1.3.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/qsocket/encrypted-stream v0.0.0-20231023165659-580d263e71f4 // indirect
	github.com/qsocket/go-srp v0.0.0-20230315175014-fb16dd9247df // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	rsc.io/qr v0.2.0 // indirect
)

replace github.com/qsocket/qsocket-go => ../../libs/qsocket-go
