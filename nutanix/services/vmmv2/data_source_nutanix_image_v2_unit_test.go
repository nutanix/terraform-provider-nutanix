package vmmv2

import (
	"testing"

	import5 "github.com/nutanix-core/ntnx-api-golang-sdk-internal/vmm-go-client/v17/models/vmm/v4/content"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func newSha1Checksum(t *testing.T, hexDigest string) *import5.OneOfImageChecksum {
	t.Helper()
	chk := import5.NewOneOfImageChecksum()
	sha1 := import5.NewImageSha1Checksum()
	sha1.HexDigest = utils.StringPtr(hexDigest)
	if err := chk.SetValue(*sha1); err != nil {
		t.Fatalf("failed to build sha1 OneOfImageChecksum: %v", err)
	}
	return chk
}

func newSha256Checksum(t *testing.T, hexDigest string) *import5.OneOfImageChecksum {
	t.Helper()
	chk := import5.NewOneOfImageChecksum()
	sha256 := import5.NewImageSha256Checksum()
	sha256.HexDigest = utils.StringPtr(hexDigest)
	if err := chk.SetValue(*sha256); err != nil {
		t.Fatalf("failed to build sha256 OneOfImageChecksum: %v", err)
	}
	return chk
}

func TestFlattenOneOfImageChecksum_Sha1(t *testing.T) {
	t.Parallel()

	got := flattenOneOfImageChecksum(newSha1Checksum(t, "abc123"))

	if len(got) != 1 {
		t.Fatalf("expected exactly one element, got %d", len(got))
	}
	if got[0]["hex_digest"] != "abc123" {
		t.Fatalf("unexpected hex_digest: %v", got[0]["hex_digest"])
	}
	if got[0]["object_type"] != "sha1" {
		t.Fatalf("expected object_type=sha1 to avoid perpetual diff, got %v", got[0]["object_type"])
	}
}

func TestFlattenOneOfImageChecksum_Sha256(t *testing.T) {
	t.Parallel()

	got := flattenOneOfImageChecksum(newSha256Checksum(t, "deadbeef"))

	if len(got) != 1 {
		t.Fatalf("expected exactly one element, got %d", len(got))
	}
	if got[0]["hex_digest"] != "deadbeef" {
		t.Fatalf("unexpected hex_digest: %v", got[0]["hex_digest"])
	}
	if got[0]["object_type"] != "sha256" {
		t.Fatalf("expected object_type=sha256 to avoid perpetual diff, got %v", got[0]["object_type"])
	}
}

func TestFlattenOneOfImageChecksum_NilInput(t *testing.T) {
	t.Parallel()

	if got := flattenOneOfImageChecksum(nil); got != nil {
		t.Fatalf("expected nil for nil input, got %v", got)
	}
}

func TestFlattenOneOfImageChecksum_UnknownObjectType(t *testing.T) {
	t.Parallel()

	chk := import5.NewOneOfImageChecksum()
	chk.ObjectType_ = utils.StringPtr("vmm.v4.content.SomeFutureChecksum")

	if got := flattenOneOfImageChecksum(chk); got != nil {
		t.Fatalf("expected nil for unknown object_type, got %v", got)
	}
}

func TestFlattenOneOfImageChecksum_EmptyObjectType(t *testing.T) {
	t.Parallel()

	chk := import5.NewOneOfImageChecksum()

	if got := flattenOneOfImageChecksum(chk); got != nil {
		t.Fatalf("expected nil for empty object_type, got %v", got)
	}
}
