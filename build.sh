#!/bin/bash
archs=(amd64 arm64 ppc64le ppc64 s390x)

for arch in ${archs[@]}
do
	env GOOS=linux GOARCH=${arch} go build -o prepnode_${arch}
done
