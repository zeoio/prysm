//go:generate type-migration -types=PeersResponse -out=out.go  -target-relative=../b -target=github.com/prysmaticlabs/prysm/tools/type-migration/b -src=github.com/prysmaticlabs/prysm/tools/type-migration/a -target-pkg=b -src-pkg=a -out-pkg=c
package c
