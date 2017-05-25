package puid

import "testing"

func Test_Base36Normalize(t *testing.T) {
	var base36test = make([]byte, 72)
	for i := range base36test {
		base36test[i] = byte(i)
	}
	var expected = "0123456789abcdefghijklmnopqrstuvwxyz0123456789abcdefghijklmnopqrstuvwxyz"
	base36convert(base36test)
	if string(base36test) != expected {
		t.Errorf("base36 normalization not what I expected: \n\t%s", base36test)
	}
}
