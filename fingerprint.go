package puid

import (
	"os"
	"strconv"
)

// a var here so we can override in tests
var getHostname = func() string {
	s, _ := os.Hostname()
	return s
}

var getPid = func() int64 {
	return int64(os.Getpid())
}

func getDefaultFingerprint() []byte {
	host := getHostname()
	if host == "" {
		host = "localhost"
	}

	return CreateFingerprint(host, getPid())
}

// Creates a BLOCK size []byte to be used as the fingerprint section
// of the ids. The default implementations create the fingerprint
// from the hostname and pid, but you can use whatever you like.
// This function ensures the created fingerprint is suitable for
// use in the puid. That is, it consists only of base36 characters.
func CreateFingerprint(str string, num int64) []byte {
	// we need BLOCK size bytes
	// most implementations use hostname and pid
	// BLOCK/2 bytes from the pid
	// BLOCK/2 bytes from the hostname
	// I am going to do it the other way around
	// to make the append easier
	half := BLOCK / 2
	fp := make([]byte, half, BLOCK) //allocate a full cap half len buffer

	// string bit
	for i := range str {
		fp[i%half] += str[i]
	}
	// now normalize that to base36
	base36convert(fp)

	// add num, for 2 digits we need to clamp the number to a max of BASE^(BLOCK/2)
	n := num % MAX_HALF_INT
	// which mean we might need to pad with '0' (character not \0)
	l := strlen(n)
	for l < half {
		fp = append(fp, '0')
		l++
	}
	return strconv.AppendInt(fp, n, BASE)
}
