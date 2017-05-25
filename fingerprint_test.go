package puid

import (
	"testing"
)

// This test is mostly for code coverage, the internal algorithm and
// output is not really relevant. The main test case is testing an intput string
// shorter than BLOCK/2 and an int < BASE^(BLOCK/2)
func Test_FingerprintCreation(t *testing.T) {
	fp := CreateFingerprint("\x00", 1)
	// assuming we keep BLOCK at 4
	// "\x00" results in "00" with my base36 normalisation
	// 1 results in a "01"
	expected := "0001"
	if string(fp) != expected {
		t.Errorf("fingerprint not created as expected. Actual `%s` Expected `%s`", fp, expected)
	}
}

func Test_CustomFingerprint(t *testing.T) {
	fp := "abcd"
	g := WithFingerprintBytes([]byte(fp))
	validateFingerprint(t, g, fp)
	// now a generated fingerprint
	g = WithFingerprint("\x00", 1)
	fp = "0001"
	validateFingerprint(t, g, fp)
}

func Test_CustomLongFingerprintShouldBeTruncated(t *testing.T) {
	fp := "acbdefg"
	g := WithFingerprintBytes([]byte(fp))

	// the fingerprint should be truncated:
	if len(g.fingerprint) != BLOCK {
		t.Error("long fingerprint not truncated")
	}

	// the truncated version should still work
	validateFingerprint(t, g, fp[0:BLOCK])
}

func validateFingerprint(t *testing.T, g *Generator, fp string) {
	id := g.New()
	// fingerprint is at 9 + BLOCK => 1 + BLOCK*2
	if id[9+BLOCK:9+BLOCK*2] != fp {
		t.Errorf("unexpected fingerprint in id: `%s`", id)
	}
	// and the next one should have the same fingerprint.
	id = g.New()
	if id[9+BLOCK:9+BLOCK*2] != fp {
		t.Errorf("unexpected fingerprint in id: `%s`", id)
	}
}

func Test_NilFingerprintCausesPanic(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("we should have panic'd on a nil Fingerprint")
		}
	}()
	WithFingerprintBytes(nil)
}

func Test_BadFingerprintCausesPanic(t *testing.T) {
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error("we should have panic'd on a bad Fingerprint")
			}
		}()
		WithFingerprintBytes([]byte{'a', 'b', 'c', 0x00})
	}()
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Error("we should have panic'd on a bad Fingerprint")
			}
		}()
		// should do the same if we do it when we create a brand new generator
		NewGenerator(&Options{
			Fingerprint: []byte{'a', 'b', 'c', 0x00},
		})
	}()
}

func Test_MissingHostnameStillWorks(t *testing.T) {
	host, pid := getHostname, getPid
	getHostname = func() string { return "" }
	getPid = func() int64 { return 0 }
	fp1 := getDefaultFingerprint()
	getHostname, getPid = host, pid
	// with the current algorithm using "localhost" (the fallback)
	// as hostname and 0 as pid, will give: `ur00` as fingerprint
	// but this should be exactly the same as calling CreateFingerprint("localhost", 0)
	fp2 := CreateFingerprint("localhost", 0)
	if string(fp1) != string(fp2) {
		t.Errorf("unexpected default fingerprint `%s`, expected `%s`", fp1, fp2)
	}
}
