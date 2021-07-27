//go:generate type-migration -types=SyncStatus,Genesis,Attestation -target-relative=../prysm/v2 -target=github.com/prysmaticlabs/prysm/proto/prysm/v2 -src=github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1 -target-pkg=v2 -src-pkg=v1Alpha1 -out=type_migration.go -out-pkg=migration
package migration
