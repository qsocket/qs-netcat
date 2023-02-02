#!/bin/bash 
## ANSI Colors (FG & BG)
RED="$(printf '\033[31m')" GREEN="$(printf '\033[32m')" YELLOW="$(printf '\033[33m')" BLUE="$(printf '\033[34m')"
MAGENTA="$(printf '\033[35m')" CYAN="$(printf '\033[36m')" WHITE="$(printf '\033[37m')" BLACK="$(printf '\033[30m')"
REDBG="$(printf '\033[41m')" GREENBG="$(printf '\033[42m')" YELLOWBG="$(printf '\033[43m')" BLUEBG="$(printf '\033[44m')"
MAGENTABG="$(printf '\033[45m')" CYANBG="$(printf '\033[46m')" WHITEBG="$(printf '\033[47m')" BLACKBG="$(printf '\033[40m')"
RESET="$(printf '\e[0m')"

## Globals
RELEASE_DIR="`pwd`/release"
BUILD_DIR="`pwd`/build"
ERR_LOG="/dev/null"
[[ ! -z $VERBOSE ]] && ERR_LOG="`tty`"

print_status() {
    echo ${YELLOW}"[*] ${RESET}${1}"
}

print_progress() {
	[[ ! -z "${VERBOSE}" ]] && return
    echo -n ${YELLOW}"[*] ${RESET}${1}"
	n=${#1}
	printf %$((70-$n))s |tr " " "."
}

print_warning() {
  echo -n ${YELLOW}"[!] ${RESET}${1}"
}

print_error() {
  echo ${RED}"[-] ${RESET}${1}"
}

print_fatal() {
  echo -e ${RED}"[!] $1\n${RESET}"
  kill -10 $$
}

print_good() {
  echo ${GREEN}"[+] ${RESET}${1}"
}

print_verbose() {
  if [[ ! -z "${VERBOSE}" ]]; then
    echo ${WHITE}"[*] ${RESET}${1}"
  fi
}

print_ok(){
	[[ -z "${VERBOSE}" ]] && echo -e " [${GREEN}OK${RESET}]"
}

print_fail(){
	[[ -z "${VERBOSE}" ]] && echo -e " [${RED}FAIL${RESET}]"
}

must_exist() {
  for i in "$@"; do
		command -v $i >$ERR_LOG || print_fatal "$i not installed! Exiting..."
  done
}

one_must_exist() {
	command -v $1 >$ERR_LOG || command -v $2 >$ERR_LOG || print_fatal "Neither $1 nor $2 installed! Exiting..."
}

## Handle SININT
exit_on_signal_SIGINT () {
  echo ""
	print_error "Script interrupted!"
  clean_exit
}

exit_on_signal_SIGTERM () {
	echo ""
  print_error "Script interrupted!"
	clean_exit
}

trap exit_on_signal_SIGINT SIGINT
trap exit_on_signal_SIGTERM SIGTERM


# Remove all artifacts and exit...
clean_exit() {
	[[ -d "$RELEASE_DIR" ]] && rm -rf "$RELEASE_DIR" &>$ERR_LOG
    kill -10 $$
}

# Expects <platform> <architecture> and creates a tar.gz archive after compressing the binary with UPX. 
# $1 = <linux>
# $2 = <amd64>
package_release_binary() {
    local bin_suffix=""
    [[ $1 == "windows" ]] && bin_suffix=".exe"
    cp "$BUILD_DIR/$1/qs-netcat-${2}${bin_suffix}" "$BUILD_DIR/qs-netcat${bin_suffix}" &>$ERR_LOG || return 1
    print_verbose "Compressing $BUILD_DIR/$1/qs-netcat-$2${bin_suffix}..."
    upx -q --best "$BUILD_DIR/qs-netcat${bin_suffix}" &>$ERR_LOG # Ignore errors...
    print_verbose "Packaging $BUILD_DIR/$1/qs-netcat-$2${bin_suffix}..."
    tar -C "$BUILD_DIR" -czvf "$RELEASE_DIR/qs-netcat_$1_$arc.tar.gz" "./qs-netcat${bin_suffix}" &>$ERR_LOG || return 1
    return 0
}


[[ ! -d $BUILD_DIR ]] && print_fatal "Could not find build firectory! Exiting..."
print_status "Initiating..."
print_status "Release Date: `date`"
echo ""
mkdir -p release
declare -a arcs=("amd64" "386" "arm" "arm64" "mips" "mips64" "mips64le" "mipsle" "ppc64" "ppc64le" "s390x")
for arc in "${arcs[@]}"
do
    print_progress "Packaging linux-$arc binary"
    package_release_binary "linux" $arc && print_ok || print_fail
done

declare -a arcs=("amd64" "386" "arm" "arm64")
for arc in "${arcs[@]}"
do
    print_progress "Packaging windows-$arc binary"
    package_release_binary "windows" $arc && print_ok || print_fail
done

declare -a arcs=("amd64" "arm64")
for arc in "${arcs[@]}"
do
    print_progress "Packaging darwin-$arc binary"
    package_release_binary "darwin" $arc && print_ok || print_fail
done

# declare -a arcs=("amd64" "arm64")
# for arc in "${arcs[@]}"
# do
#     echo -n "[*] Packaging ios-$arc binary -> "
#     cp "$BUILD_DIR/ios/qs-netcat-$arc" "$BUILD_DIR/qs-netcat"
#     upx -q --best "$BUILD_DIR/qs-netcat" &>/dev/null
#     tar -C "$BUILD_DIR" -czvf "$RELEASE_DIR/qs-netcat_ios_$arc.tar.gz" "./qs-netcat"
# done

declare -a arcs=("amd64" "386" "arm" "arm64")
for arc in "${arcs[@]}"
do
    print_progress "Packaging android-$arc binary"
    package_release_binary "android" $arc && print_ok || print_fail
done

declare -a arcs=("amd64" "386" "arm" "arm64")
for arc in "${arcs[@]}"
do
    print_progress "Packaging freebsd-$arc binary"
    package_release_binary "freebsd" $arc && print_ok || print_fail
done

declare -a arcs=("amd64" "arm" "arm64" "mips64")
for arc in "${arcs[@]}"
do
    print_progress "Packaging openbsd-$arc binary"
    package_release_binary "openbsd" $arc && print_ok || print_fail
done

declare -a arcs=("amd64" "386" "arm" "arm64")
for arc in "${arcs[@]}"
do
    print_progress "Packaging netbsd-$arc binary"
    package_release_binary "netbsd" $arc && print_ok || print_fail
done

# Special distro cases...
print_progress "Packaging android APK"
cp "$BUILD_DIR/android/qs-netcat.apk" "$RELEASE_DIR/" && print_ok || print_fail

print_progress "Packaging solaris-amd64 binary"
package_release_binary "netbsd" "amd64" && print_ok || print_fail

print_progress "Packaging illumos-amd64 binary"
package_release_binary "illumos" "amd64" && print_ok || print_fail

print_progress "Packaging dragonfly-amd64 binary"
package_release_binary "dragonfly" "amd64" && print_ok || print_fail

print_progress "Packaging aix-ppc64 binary"
package_release_binary "aix" "ppc64" && print_ok || print_fail

print_good "All done!"

cd $RELEASE_DIR
echo -e "\n\`\`\`"
sha1sum *
echo -e "\`\`\`\n"
