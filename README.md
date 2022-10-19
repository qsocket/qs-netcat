<div align="center">
  <img src=".github/img/banner.png">
  <br>
  <br>


  [![GitHub All Releases][release-img]][release]
  [![Build][workflow-img]][workflow]
  [![Issues][issues-img]][issues]
  [![Go Report Card][go-report-img]][go-report]
  ![Docker Pulls][docker-pulls]
  [![License: MIT][license-img]][license]
</div>

[go-report]: https://goreportcard.com/report/github.com/qsocket/qs-netcat
[go-report-img]: https://goreportcard.com/badge/github.com/qsocket/qs-netcat
[release]: https://github.com/qsocket/qs-netcat/releases
[release-img]: https://img.shields.io/github/v/release/qsocket/qs-netcat
[downloads]: https://github.com/qsocket/qs-netcat/releases
[downloads-img]: https://img.shields.io/github/downloads/qsocket/qs-netcat/total?logo=github
[issues]: https://github.com/qsocket/qs-netcat/issues
[issues-img]: https://img.shields.io/github/issues/qsocket/qs-netcat?color=red
[docker-pulls]: https://img.shields.io/docker/pulls/qsocket/qs-netcat?logo=docker&label=docker%20pulls
[license]: https://raw.githubusercontent.com/qsocket/qs-netcat/master/LICENSE
[license-img]: https://img.shields.io/github/license/qsocket/qs-netcat.svg
[google-cloud-shell]: https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/qsocket/qs-netcat&tutorial=README.md
[workflow-img]: https://github.com/qsocket/qs-netcat/actions/workflows/main.yml/badge.svg
[workflow]: https://github.com/qsocket/qs-netcat/actions/workflows/main.yml
[qsrn]: https://github.com/qsocket/qsrn

qs-netcat is a cross-platform networking utility which reads and writes data across systems using the [QSRN](qsrn). 
It allows redirecting true PTY session with reverse connections effectively backdooring systems, creating TCP proxies, and transfering files to/from systems under NAT networks.

## Installation

[![Open in Cloud Shell](.github/img/cloud-shell.png)](google-cloud-shell)

|    **Tool**   |                 **Build From Source**                |       **Docker Image**      |                     **Binary Release**                    |
|:-------------:|:----------------------------------------------------:|:---------------------------:|:---------------------------------------------------------:|
| **qs-netcat** | ```go install github.com/qsocket/qs-netcat@latest``` | [Download](#docker-install) | [Download](release) |

---

qs-netcat supports 10 architectures and 12 operating systems, check **Supported Platforms** below for detailed table. 
<details>
<summary>Supported Platforms</summary>

|  **Platform** | **AMD64** | **386** | **ARM** | **ARM64** | **MIPS** | **MIPS64** | **MIPS64LE** | **PPC64** | **PPC64LE** | **S390X** |
|:-------------:|:---------:|:-------:|:-------:|:---------:|:--------:|:----------:|:------------:|:---------:|:-----------:|:---------:|
|   **Linux**   |     ✅     |    ✅    |    ✅    |     ✅     |     ✅    |      ✅     |       ✅      |     ✅     |      ✅      |     ✅     |
|   **Darwin**  |     ✅     |    ❌    |    ❌    |     ✅     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
|  **Windows**  |     ✅     |    ✅    |    ✅    |     ✅     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
|  **OpenBSD**  |     ✅     |    ✅    |    ✅    |     ✅     |     ❌    |      ✅     |       ❌      |     ❌     |      ❌      |     ❌     |
|   **NetBSD**  |     ✅     |    ✅    |    ✅    |     ✅     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
|  **FreeBSD**  |     ✅     |    ✅    |    ✅    |     ✅     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
|  **Android**  |     ✅     |    ✅    |    ✅    |     ✅     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
|    **IOS**    |     ✅     |    ❌    |    ❌    |     ✅     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
|  **Solaris**  |     ✅     |    ❌    |    ❌    |     ❌     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
|  **Illumos**  |     ✅     |    ❌    |    ❌    |     ❌     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
| **Dragonfly** |     ✅     |    ❌    |    ❌    |     ❌     |     ❌    |      ❌     |       ❌      |     ❌     |      ❌      |     ❌     |
|    **AIX**    |     ❌     |    ❌    |    ❌    |     ❌     |     ❌    |      ❌     |       ❌      |     ✅     |      ❌      |     ❌     |

</details>


### Docker Install

[![Docker](http://dockeri.co/image/egee/qsocket)](https://hub.docker.com/r/egee/qsocket/)

```
docker pull egee/qsocket
docker run -it egee/qsocket
```

## Usage

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
### Examples
- SSH from *Workstation B* to *Workstation A* through any firewall/NAT
```bash
$ qs-netcat -f "localhost" -p 22 -l     # Workstation A
$ qsocket ssh root@qsocket.io        # Workstation B
```
- Log in to Workstation A from Workstation B through any firewall/NAT
```bash
$ qs-netcat -l -i   # Workstation A
$ qs-netcat -i      # Workstation B
```
- Transfer files from *Workstation B* to *Workstation A*
```bash
$ qs-netcat -q -s MySecret -l > file.txt     # Workstation A
$ qs-netcat -q -s MySecret < file.txt        # Workstation B
```

<details>
<summary>SSH connection over QSRN</summary>

https://user-images.githubusercontent.com/1161307/171013513-95f18734-233d-45d3-aaf5-d6aec687db0e.mov

</details>

<details>
<summary>RDP connection over QSRN</summary>

https://user-images.githubusercontent.com/1161307/171013513-95f18734-233d-45d3-aaf5-d6aec687db0e.mov

</details>