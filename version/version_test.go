package version

import (
	"strings"
	"testing"
)

func BenchmarkGetHumanVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GetHumanVersion()
	}
}

func TestExtractVersion(t *testing.T) {

	expectedVersion := strings.TrimSpace(Version)
	expectedPrerelease := strings.TrimSpace(VersionPrerelease)
	if Version != expectedVersion {
		t.Errorf("Version mismatch. Expected: %s, Got: %s", expectedVersion, Version)
	}

	if VersionPrerelease != expectedPrerelease {
		t.Errorf("VersionPrerelease mismatch. Expected: %s, Got: %s", expectedPrerelease, VersionPrerelease)
	}
}

func TestGetHumanVersion(t *testing.T) {

	expectedResult := Version + "-" + VersionPrerelease
	if VersionMetadata != "" {
		expectedResult += "+" + VersionMetadata
	}

	result := GetHumanVersion()

	if result != expectedResult {
		t.Errorf("Unexpected result. Expected: %s, Got: %s", expectedResult, result)
	}

	if len(result) != len(strings.TrimSpace(result)) {
		t.Errorf("Unexpected characters in the result. Result: %s", result)
	}
}
