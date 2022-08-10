CURRET_DIR=$(shell pwd)
BUILD=go build
OUT_DIR=${CURRET_DIR}/build
BUILD_FLAGS=-trimpath -buildvcs --ldflags "-s -w -X github.com/qsocket/qs-netcat/config.Version=$$(git log --pretty=format:'v1.0.%at-%h' -n 1)" 
CGO_ENABLED=0
$(shell mkdir -p build/{windows,linux,macos,android,freebsd,openbsd,solaris,aix})

default: linux
windows:
	GOOS=windows GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/windows/qs-netcat.exe
	GOOS=windows GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/windows/qs-netcat32.exe
linux:
	GOOS=linux GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat
	GOOS=linux GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat32
	GOOS=linux GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-arm
	GOOS=linux GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-arm64
	GOOS=linux GOARCH=mips ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-mips
	GOOS=linux GOARCH=mips64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-mips64
	GOOS=linux GOARCH=mips64le ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-mips64le
	GOOS=linux GOARCH=mipsle ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-mipsle
	GOOS=linux GOARCH=ppc64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/linux/qs-netcat-ppc64
freebsd:
	GOOS=freebsd GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/freebsd/qs-netcat
	GOOS=freebsd GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/freebsd/qs-netcat32
	GOOS=freebsd GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/freebsd/qs-netcat-arm
	GOOS=freebsd GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/freebsd/qs-netcat-arm64
openbsd:
	GOOS=openbsd GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat
	GOOS=openbsd GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat32
	GOOS=openbsd GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat-arm
	GOOS=openbsd GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/openbsd/qs-netcat-arm64
netbsd:
	GOOS=netbsd GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/netbsd/qs-netcat
	GOOS=netbsd GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/netbsd/qs-netcat32
	GOOS=netbsd GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/netbsd/qs-netcat-arm
	GOOS=netbsd GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/netbsd/qs-netcat-arm64
android: # android builds require native development kit
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/android/qs-netcat
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=386 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/android/qs-netcat32
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=arm ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/android/qs-netcat-arm
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/android/qs-netcat-arm64
macos:
	GOOS=darwin GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/darwin/qs-netcat
	GOOS=darwin GOARCH=arm64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/darwin/qs-netcat-arm64
solaris:
	GOOS=solaris GOARCH=amd64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/darwin/qs-netcat
aix:
	GOOS=aix GOARCH=ppc64 ${BUILD} ${BUILD_FLAGS} -o ${OUT_DIR}/aix/qs-netcat

all: windows linux freebsd openbsd netbsd macos solaris aix # qs-netcat-android 
