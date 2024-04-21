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

LINUX='linux:0:amd64:linux-amd64'
WIN='windows:0:amd64:win-amd64.exe'
MACAMD='darwin:0:amd64:darwin-amd64'
MACARM='darwin:0:arm64:darwin-arm64'

for II in $LINUX $WIN $MACAMD $MACARM; do
	os=$(echo $II | cut -d":" -f1)
	cgo=$(echo $II | cut -d":" -f2)
	arch=$(echo $II | cut -d":" -f3)
	suffix=$(echo $II | cut -d":" -f4)
	# echo $os $arch $suffix;
	GOOS=${os} GOARCH=${arch} CGO_ENABLED=${cgo} go build -o ${THISDIR}/${BASEBIN}-${suffix} ${TARGET}
done
