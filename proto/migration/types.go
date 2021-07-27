//go:generate type-migration -types=SyncStatus -target=github.com/prysmaticlabs/prysm/proto/prysm/v2 -src=github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1 -target-pkg=v2 -src-pkg=v1Alpha1 -out=types_migration.go -out-pkg=migration
package migration
