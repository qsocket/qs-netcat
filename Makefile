CURRET_DIR=$(shell pwd)
BUILD=go build
OUT_DIR=${CURRET_DIR}/build
BUILD_FLAGS=-trimpath -buildvcs --ldflags "-s -w -X github.com/qsocket/qs-netcat/config.Version=$$(git log --pretty=format:'v1.0.%at-%h' -n 1)" 
CGO_ENABLED=0
$(shell mkdir -p build/{windows,linux,darwin,android,ios,freebsd,openbsd,solaris,aix,illumos,dragonfly})

default: linux
windows:
	GOOS=windows GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/windows/qs-netcat-amd64.exe
	GOOS=windows GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/windows/qs-netcat-386.exe
	GOOS=windows GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/windows/qs-netcat-arm.exe
	GOOS=windows GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/windows/qs-netcat-arm64.exe
linux:
	GOOS=linux GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-amd64
	GOOS=linux GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-386
	GOOS=linux GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-arm
	GOOS=linux GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-arm64
	GOOS=linux GOARCH=mips ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-mips
	GOOS=linux GOARCH=mips64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-mips64
	GOOS=linux GOARCH=mips64le ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-mips64le
	GOOS=linux GOARCH=mipsle ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-mipsle
	GOOS=linux GOARCH=ppc64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-ppc64
	GOOS=linux GOARCH=ppc64le ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-ppc64le
	GOOS=linux GOARCH=s390x ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-s390x
freebsd:
	GOOS=freebsd GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/freebsd/qs-netcat-amd64
	GOOS=freebsd GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/freebsd/qs-netcat-386
	GOOS=freebsd GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/freebsd/qs-netcat-arm
	GOOS=freebsd GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/freebsd/qs-netcat-arm64
openbsd:
	GOOS=openbsd GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat-amd64
	GOOS=openbsd GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat-386
	GOOS=openbsd GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat-arm
	GOOS=openbsd GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat-arm64
	GOOS=openbsd GOARCH=mips64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat-mips64
netbsd:
	GOOS=netbsd GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/netbsd/qs-netcat-amd64
	GOOS=netbsd GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/netbsd/qs-netcat-386
	GOOS=netbsd GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/netbsd/qs-netcat-arm
	GOOS=netbsd GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/netbsd/qs-netcat-arm64
android: # android builds require native development kit
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/android/qs-netcat
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/android/qs-netcat32
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/android/qs-netcat-arm
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/android/qs-netcat-arm64
ios:
	GOOS=ios GOARCH=amd64 CGO_ENABLED=1 CC=${CURRET_DIR}/clangwrap.sh ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/ios/qs-netcat-amd64
	GOOS=ios GOARCH=arm CGO_ENABLED=1 CC=${CURRET_DIR}/clangwrap.sh ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/ios/qs-netcat-arm64
darwin:
	GOOS=darwin GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/darwin/qs-netcat-amd64
	GOOS=darwin GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/darwin/qs-netcat-arm64
solaris:
	GOOS=solaris GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/solaris/qs-netcat-amd64
illumos:
	GOOS=illumos GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/illumos/qs-netcat-amd64
aix:
	GOOS=aix GOARCH=ppc64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/aix/qs-netcat-ppc64
dragonfly:
	GOOS=dragonfly GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/dragonfly/qs-netcat-amd64


all: linux windows darwin freebsd openbsd netbsd solaris aix dragonfly illumos # android ios 
