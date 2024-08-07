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
[docker-pulls]: https://img.shields.io/docker/pulls/qsocket/qsocket?logo=docker&label=docker%20pulls
[license]: https://raw.githubusercontent.com/qsocket/qs-netcat/master/LICENSE
[license-img]: https://img.shields.io/github/license/qsocket/qs-netcat.svg
[google-cloud-shell]: https://console.cloud.google.com/cloudshell/open?git_repo=https://github.com/qsocket/qs-netcat&tutorial=README.md
[workflow-img]: https://github.com/qsocket/qs-netcat/actions/workflows/main.yml/badge.svg
[workflow]: https://github.com/qsocket/qs-netcat/actions/workflows/main.yml
[qsrn]: https://www.qsocket.io/qsrn/

qs-netcat is a cross-platform networking utility which reads and writes E2E encrypted data across systems using the QSocket relay network ([QSRN][qsrn]).
It allows redirecting fully interactive PTY sessions with reverse connections, effectively allowing remote access to systems, forwarding traffic, and transferring files to and from systems under NAT networks or firewalls.

> [!WARNING]  
> This tool is in its early alpha development stage, featuring experimental functionality that may lack backwards compatibility, and users are advised to exercise caution and not use it in production environments.

> [!CAUTION]
> Due to the changes to the relay protocol, clients starting from version `v0.0.8-beta` will not be compatible with the older versions. You can still access your older instances using previous client versions `(>= v.0.0.7-beta)` until October 23, 2024. **After October 23, 2024, legacy relay support will end, and all out-of-date QSocket instances will become inaccessible!**



## Installation

[![Open in Cloud Shell](.github/img/cloud-shell.png)][google-cloud-shell]

|    **Tool**   |                 **Build From Source**                |       **Docker Image**      |                     **Binary Release**                    |
|:-------------:|:----------------------------------------------------:|:---------------------------:|:---------------------------------------------------------:|
| **qs-netcat** | ```go install github.com/qsocket/qs-netcat@master``` | [Download](#docker-install) | [Download](release) |

---

qs-netcat supports 10 architectures and 12 operating systems, following table contains detailed list of all **Supported Platforms**. 

<details>
<summary>Supported Platforms</summary>

- ✅ `Supported`
- 🚧 `In progress`
- ❌ `Unsupported`

| **ARCH/OS** | **Linux** | **MacOS** | **Windows** | **Android** | **IOS** | **FreeBSD** | **OpenBSD** | **NetBSD** | **Solaris** | **Illumos** | **Dragonfly** | **AIX** |
|:-----------:|:---------:|:---------:|:-----------:|:-----------:|:-------:|:-----------:|:-----------:|:----------:|:-----------:|:-----------:|:-------------:|:-------:|
|  **AMD64**  |     ✅     |     ✅     |      ✅      |      ✅      |    🚧    |      ✅      |      ✅      |      ✅     |      ✅      |      ✅      |       ✅       |    ❌    |
|  **ARM64**  |     ✅     |     ✅     |      ✅      |      ✅      |    🚧    |      ✅      |      ✅      |      ❌     |      ❌      |      ❌      |       ❌       |    ❌    |
|   **386**   |     ✅     |     ❌     |      ✅      |      ✅      |    ❌    |      ✅      |      ✅      |      ✅     |      ❌      |      ❌      |       ❌       |    ❌    |
|  **ARM32**  |     ✅     |     ❌     |      ✅      |      ✅      |    ❌    |      ✅      |      ✅      |      ✅     |      ❌      |      ❌      |       ❌       |    ❌    |
| **RISCV64** |     🚧     |     ❌     |      ❌      |      ❌      |    ❌    |      🚧      |      ❌      |      ❌     |      ❌      |      ❌      |       ❌       |    ❌    |
|  **MIPS64** |     ✅     |     ❌     |      ❌      |      ❌      |    ❌    |      ❌      |      ❌      |      ❌     |      ❌      |      ❌      |       ❌       |    ❌    |
|  **MIPS32** |     ✅     |     ❌     |      ❌      |      ❌      |    ❌    |      ❌      |      ❌      |      ❌     |      ❌      |      ❌      |       ❌       |    ❌    |
|  **MIPSLE** |     ✅     |     ❌     |      ❌      |      ❌      |    ❌    |      ❌      |      ❌      |      ❌     |      ❌      |      ❌      |       ❌       |    ❌    |
|  **PPC64**  |     ✅     |     ❌     |      ❌      |      ❌      |    ❌    |      ❌      |      🚧      |      ❌     |      ❌      |      ❌      |       ❌       |    🚧    |
| **PPC64LE** |     ✅     |     ❌     |      ❌      |      ❌      |    ❌    |      ❌      |      ❌      |      ❌     |      ❌      |      ❌      |       ❌       |    ❌    |
|  **S390X**  |     ✅     |     ❌     |      ❌      |      ❌      |    ❌    |      ❌      |      ❌      |      ❌     |      ❌      |      ❌      |       ❌       |    ❌    |

</details>


### Docker Install

[![Docker](http://dockeri.co/image/qsocket/qsocket)](https://hub.docker.com/r/qsocket/qsocket/)

```bash
docker pull qsocket/qsocket:latest
docker run -it qsocket -h
```

## Usage

```
Usage: qs-netcat

Flags:
  -h, --help                 Show context-sensitive help.
  -s, --secret=STRING        Secret (e.g. password).
  -e, --exec=STRING          Execute command [e.g. "bash -il" or "cmd.exe"]
  -f, --forward=STRING       IP:PORT for traffic forwarding.
  -x, --socks=STRING         User socks proxy address for connecting QSRN.
      --cert-fp=STRING       Hex encoded TLS certificate fingerprint for validation.
  -n, --probe=5              Probe interval for connecting QSRN.
  -C, --plain                Disable all encryption.
      --e2e                  Use E2E encryption. (default:true)
  -i, --interactive          Execute with a PTY shell.
  -l, --listen               Server mode. (listen for connections)
  -g, --generate             Generate a Secret. (random)
  -K, --pin                  Enable certificate pinning on TLS connections.
  -q, --quiet                Quiet mode. (no stdout)
  -T, --tor                  Use TOR for connecting QSRN.
      --qr                   Generate a QR code with given stdin and print on the terminal.
  -v, --verbose              Verbose mode.
      --in-pipe=IN-PIPE
      --out-pipe=OUT-PIPE
      --version

Example to forward traffic from port 2222 to 192.168.6.7:22:
  $ qs-netcat -s MyCecret -f 2222:192.168.6.7:22
Example file transfer:
	$ qs-netcat -l -s MyCecret > warez.tar.gz         # Server
	$ qs-netcat -s MyCecret < warez.tar.gz            # Client
Example for a reverse shell:
	$ qs-netcat -s MyCecret -l -i                     # Server
	$ qs-netcat -s MyCecret -i                        # Client
```

### Examples
- Log in to Workstation A from Workstation B through any firewall/NAT
```bash
qs-netcat -l -i   # Workstation A
qs-netcat -i      # Workstation B
```

- SSH from *Workstation A* to *Workstation B* by port forwarding through any firewall/NAT
```bash
qs-netcat -l                    # Workstation B
qs-netcat -f "22:localhost:22"  # Workstation A
ssh user@localhost              # Workstation A
```

- Transfer files from *Workstation B* to *Workstation A* using smart pipes
```bash
qs-netcat -s MySecret -l > file.txt     # Workstation A
qs-netcat -s MySecret < file.txt        # Workstation B
```

- Port forward. Access 192.168.6.7:80 on Workstation A's private LAN from Workstation B:
```bash
qs-netcat -l                  # Workstation A
qs-netcat -f 192.168.6.7:80   # Workstation B
```

- Execute any command (nc -e style) on *Workstation A*
```bash
qs-netcat -l                         # Workstation A
qs-netcat -e "echo hello_world; id"  # Workstation B
```
- Access entirety of Workstation A's private LAN (Sock4/4a/5 proxy)
```bash
qs-netcat -l                    # Workstation A
qs-netcat -f "22:localhost:22"  # Workstation B
ssh -D 9090 root@localhost      # Workstation B
# Access www.google.com via Workstation A's private LAN from your Workstation B:
curl --socks4a 127.1:9090 http://www.google.com
```

- Mount a remote folder of Workstation A using sshfs and qs-netcat
```bash
qs-netcat -l                    # Workstation A
qs-netcat -f "22:localhost:22"  # Workstation B
sudo sshfs -o allow_other,default_permissions root@localhost:/remote_dir /mnt/local_dir # Workstation B
```

---

**Crypto / Security Mumble Jumble**
- The connections are end-2-end encrypted. This means from User-2-User (and not just to the Relay Network). The Relay Network relays only (encrypted) data to and from the Users.
- The QSocket uses [SRP](https://en.wikipedia.org/wiki/Secure_Remote_Password_protocol) for ensuring [perfect forward secrecy](https://en.wikipedia.org/wiki/Forward_secrecy). This means that the session keys are always different, and recorded session traffic cannot be decrypted by the third parties even if the user secret is known.
- The session key is 256 bit and ephemeral. It is freshly generated for every session and generated randomly (and is not based on the password).
- A brute force attack against weak secrets requires a new TCP connection for every guess. But QSRN contains a strong load balancer which is limiting the consecutive connection attempts.
- Do not use stupid passwords like 'password123'. Malice might pick the same (stupid) password by chance and connect. If in doubt use *qs-netcat -g* to generate a strong one. Alice's and Bob's password should at least be strong enough so that Malice can not guess it by chance while Alice is waiting for Bob to connect.
- If Alice shares the same password with Bob and Charlie and either one of them connects then Alice can not tell if it is Bob or Charlie who connected.
- Assume Alice shares the same password with Bob and Malice. When Alice stops listening for a connection then Malice could start to listen for the connection instead. Bob (when opening a new connection) can not tell if he is connecting to Alice or to Malice.
- We did not invent SRP. It's a well-known protocol, and it is well-analyzed and trusted by the community. 


https://user-images.githubusercontent.com/17179401/224060762-e0f121f6-431b-4eb5-8833-4a5d533003de.mp4

---

<details>
<summary>RDP connection over QSRN</summary>

https://github.com/qsocket/qs-netcat/assets/17179401/af46c8fb-cb33-483a-b5c1-9142843da2bd

</details>


<details>
<summary>ADB access over QSRN</summary>

https://user-images.githubusercontent.com/17179401/216651601-6ddc8ddf-7248-4c2b-bd77-00f00f773c80.mov
    
</details>
