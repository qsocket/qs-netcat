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
[qsrn]: https://www.qsocket.io/qsrn/

qs-netcat is a cross-platform networking utility which reads and writes data across systems using the [QSRN][qsrn].
It allows redirecting true PTY sessions with reverse connections, effectively allowing remote access to systems, creating TCP proxies, and transferring files to and from systems under NAT networks or firewalls.

## Installation

[![Open in Cloud Shell](.github/img/cloud-shell.png)][google-cloud-shell]

|    **Tool**   |                 **Build From Source**                |       **Docker Image**      |                     **Binary Release**                    |
|:-------------:|:----------------------------------------------------:|:---------------------------:|:---------------------------------------------------------:|
| **qs-netcat** | ```go install github.com/qsocket/qs-netcat@latest``` | [Download](#docker-install) | [Download](release) |

---

qs-netcat supports 10 architectures and 12 operating systems, following table contains detailed list of all **Supported Platforms**. 
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

```bash
docker pull qsocket/qs-netcat
docker run -it qsocket/qs-netcat -h
```

## Usage

```
Usage: qs-netcat

Flags:
  -h, --help              Show context-sensitive help.
  -s, --secret=STRING     Secret (e.g. password).
  -e, --exec=STRING       Execute command [e.g. "bash -il" or "cmd.exe"]
  -f, --forward=STRING    IP:PORT for traffic forwarding.
  -n, --probe=5           Probe interval for connecting QSRN.
  -C, --no-tls            Disable TLS encryption.
  -i, --interactive       Execute with a PTY shell.
  -l, --listen            Server mode. (listen for connections)
  -g, --generate          Generate a Secret. (random)
  -K, --pin               Enable certificate pinning on TLS connections.
  -q, --quiet             Quiet mode. (no stdout)
  -T, --tor               Use TOR for connecting QSRN.
  -v, --verbose           Verbose mode.
      --version

Example to forward traffic from port 2222 to 192.168.6.7:22:
	$ qs-netcat -s MyCecret -l -f 192.168.6.7:22        # Server
	$ qs-netcat -s MyCecret -f :2222                    # Client
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
$ qs-netcat -f "localhost" -p 22 -l  # Workstation A
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

https://user-images.githubusercontent.com/17179401/221196823-5c6e3a66-3b06-410a-9d2b-33efd101428a.mp4

---

<details>
<summary>RDP connection over QSRN</summary>

https://user-images.githubusercontent.com/17179401/213314447-65ecaf43-89fd-48bd-a242-3345f6baf185.mov

</details>


<details>
<summary>ADB access over QSRN</summary>

https://user-images.githubusercontent.com/17179401/216651601-6ddc8ddf-7248-4c2b-bd77-00f00f773c80.mov

</details>