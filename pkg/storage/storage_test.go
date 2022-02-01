package storage

import (
	"testing"
)

const fileSuffix = ".json"

var encodeFileNameTestData = []struct {
	testName   string
	tenant     string
	name       string
	fileSuffix string
	out        string
	ee         error
}{
	{"valid1", "ten1", "name1", "", "ten1/name1", nil},
	{"valid1", "ten1", "name1", fileSuffix, "ten1/name1.json", nil},
	{"valid2", "t", "n", fileSuffix, "t/n.json", nil},
	{"valid2", "t", "n", "", "t/n", nil},
	{"valid2", "t", "n.akarmi", "", "t/n.akarmi", nil},
	{"error1", "", "", "", "", ErrTenantEmpty},
	{"error2", "adfdafda", "", "", "", ErrNameEmpty},
}

var decodeFileNameTestData = []struct {
	testName   string
	fileName   string
	fileSuffix string
	out        string
	ee         error
}{
	{"valid1", "ten1/name1", "", "name1", nil},
	{"valid1", "ten1/name1.akarmi", "", "name1.akarmi", nil},
	{"valid1", "ten1/name1.json", fileSuffix, "name1", nil},
	{"valid2", "t/n", "", "n", nil},
	{"valid2", "t/n.json", fileSuffix, "n", nil},
	{"error1", "", "", "", ErrFileNameEmpty},
	{"error2", "wrong.extension", "", "", ErrFileName("wrong.extension")},
	{"error3", "noseparatornoextension", "", "", ErrFileName("noseparatornoextension")},
	{"error4", "multiple/separator/noextension", "", "", ErrFileName("multiple/separator/noextension")},
	{"error5", "multiple/separator/wrong.ext", "", "", ErrFileName("multiple/separator/wrong.ext")},
}

func TestEncodeFileName(t *testing.T) {
	for _, tt := range encodeFileNameTestData {
		t.Run(tt.testName, func(t *testing.T) {
			s, err := encodeFileName(tt.tenant, tt.name, tt.fileSuffix)
			good := s == tt.out && err == tt.ee
			if !good {
				t.Errorf("\n got: %q, %v \n want: %q, %v", s, err, tt.out, tt.ee)
			}
		})
	}
}

func TestDecodeFileName(t *testing.T) {
	for _, tt := range decodeFileNameTestData {
		t.Run(tt.testName, func(t *testing.T) {
			s, err := decodeFileName(tt.fileName, tt.fileSuffix)
			// Here err.Error() to comapre the error strings
			good := s == tt.out && errorsEqual(err, tt.ee)
			if !good {
				t.Errorf("\n got %q, %v \n want %q, %v", s, err, tt.out, tt.ee)
			}
		})
	}
}

func errorsEqual(e, ee error) bool {
	if e == nil {
		return ee == nil
	}
	if ee == nil {
		return e == nil
	}
	return e.Error() == ee.Error()
}
