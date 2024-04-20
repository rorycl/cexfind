#!/bin/bash

# build many go binary output targets
# https://freshman.tech/snippets/go/cross-compile-go-programs/

if [ "$#" -ne 2 ]; then
    echo "please provide the file target to compile, followed by the"
    echo "base name of the output binary as arguments"
	exit 1
fi

TARGET=$1
BASEBIN=$2

THISDIR=$(dirname "$0")

LINUX='linux:amd64:linux-amd64'
WIN='windows:amd64:win-amd64.exe'
MACAMD='darwin:amd64:darwin-amd64'
MACARM='darwin:arm64:darwin-arm64'

for II in $LINUX $WIN $MACAMD $MACARM; do
	os=$(echo $II | cut -d":" -f1)
	arch=$(echo $II | cut -d":" -f2)
	suffix=$(echo $II | cut -d":" -f3)
	# echo $os $arch $suffix;
	GOOS=${os} GOARCH=${arch} go build -o ${THISDIR}/${BASEBIN}-${suffix} ${TARGET}
done
