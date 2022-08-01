CURRET_DIR=$(shell pwd)
BUILD=go build
OUT_DIR=${CURRET_DIR}/build
BUILD_FLAGS=-trimpath -buildvcs --ldflags "-s -w -X main.Version=$$(git log --pretty=format:'v1.0.%at-%h' -n 1)" 
CGO_ENABLED=0
$(shell mkdir -p build/{windows,linux,macos,android,freebsd,openbsd,solaris})

default: qs-netcat-linux
qs-netcat-win:
	GOOS=windows GOARCH=amd64 ${BUILD} -o ${OUT_DIR}/windows/qs-netcat.exe
	GOOS=windows GOARCH=386 ${BUILD} -o ${OUT_DIR}/windows/qs-netcat32.exe
qs-netcat-linux:
	GOOS=linux GOARCH=amd64 ${BUILD} -o ${OUT_DIR}/linux/qs-netcat
	GOOS=linux GOARCH=386 ${BUILD} -o ${OUT_DIR}/linux/qs-netcat32
	GOOS=linux GOARCH=arm ${BUILD} -o ${OUT_DIR}/linux/qs-netcat-arm
	GOOS=linux GOARCH=arm64 ${BUILD} -o ${OUT_DIR}/linux/qs-netcat-arm64
	GOOS=linux GOARCH=mips ${BUILD} -o ${OUT_DIR}/linux/qs-netcat-mips
	GOOS=linux GOARCH=mips64 ${BUILD} -o ${OUT_DIR}/linux/qs-netcat-mips64
	GOOS=linux GOARCH=mips64le ${BUILD} -o ${OUT_DIR}/linux/qs-netcat-mips64le
	GOOS=linux GOARCH=mipsle ${BUILD} -o ${OUT_DIR}/linux/qs-netcat-mipsle
qs-netcat-freebsd:
	GOOS=freebsd GOARCH=amd64 ${BUILD} -o ${OUT_DIR}/freebsd/qs-netcat
	GOOS=freebsd GOARCH=386 ${BUILD} -o ${OUT_DIR}/freebsd/qs-netcat32
	GOOS=freebsd GOARCH=arm ${BUILD} -o ${OUT_DIR}/freebsd/qs-netcat-arm
	GOOS=freebsd GOARCH=arm64 ${BUILD} -o ${OUT_DIR}/freebsd/qs-netcat-arm64
qs-netcat-openbsd:
	GOOS=openbsd GOARCH=amd64 ${BUILD} -o ${OUT_DIR}/openbsd/qs-netcat
	GOOS=openbsd GOARCH=386 ${BUILD} -o ${OUT_DIR}/openbsd/qs-netcat32
	GOOS=openbsd GOARCH=arm ${BUILD} -o ${OUT_DIR}/openbsd/qs-netcat-arm
	GOOS=openbsd GOARCH=arm64 ${BUILD} -o ${OUT_DIR}/openbsd/qs-netcat-arm64
qs-netcat-netbsd:
	GOOS=netbsd GOARCH=amd64 ${BUILD} -o ${OUT_DIR}/netbsd/qs-netcat
	GOOS=netbsd GOARCH=386 ${BUILD} -o ${OUT_DIR}/netbsd/qs-netcat32
	GOOS=netbsd GOARCH=arm ${BUILD} -o ${OUT_DIR}/netbsd/qs-netcat-arm
	GOOS=netbsd GOARCH=arm64 ${BUILD} -o ${OUT_DIR}/netbsd/qs-netcat-arm64
qs-netcat-android: # android builds require native development kit
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=amd64 ${BUILD} -o ${OUT_DIR}/android/qs-netcat
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=386 ${BUILD} -o ${OUT_DIR}/android/qs-netcat32
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=arm ${BUILD} -o ${OUT_DIR}/android/qs-netcat-arm
	CC=$NDK_ROOT/21.3.6528147/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android30-clang GOOS=android GOARCH=arm64 ${BUILD} -o ${OUT_DIR}/android/qs-netcat-arm64
qs-netcat-mac:
	GOOS=darwin GOARCH=amd64 ${BUILD} -o ${OUT_DIR}/darwin/qs-netcat
	GOOS=darwin GOARCH=arm64 ${BUILD} -o ${OUT_DIR}/darwin/qs-netcat-arm64
qs-netcat-solaris:
	GOOS=solaris GOARCH=amd64 ${BUILD} -o ${OUT_DIR}/darwin/qs-netcat
all: qs-netcat-win qs-netcat-linux qs-netcat-freebsd qs-netcat-openbsd qs-netcat-netbsd qs-netcat-mac qs-netcat-solaris # qs-netcat-android 
