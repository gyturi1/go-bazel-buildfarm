package storage

import (
	"errors"
	"io/fs"
	"os"
	"strings"
	"testing"
)

func TestNewGoogleStorageBucketNameEmpty(t *testing.T) {
	os.Setenv(BucketEnvParamName, "")
	_, err := GoogleBucket(nil)
	if err != ErrNoBucket {
		t.Errorf("\n got:%v \n want: %v", err, ErrNoBucket)
	}
}

func TestNewGoogleStorageMissingCredFileEnvParam(t *testing.T) {
	os.Setenv(BucketEnvParamName, "bucketName")
	_, err := GoogleBucket(nil)
	if err == nil || err.Error() != ErrNoCredFilePath.Error() {
		t.Errorf("\n got: %v \n want: %v", err, ErrNoCredFilePath)
	}
}

func TestNewGoogleStorageMissingCredFile(t *testing.T) {
	os.Setenv(BucketEnvParamName, "bucketName")
	os.Setenv(CredJSONEnvParamName, "no_file")
	_, err := GoogleBucket(nil)
	if errors.Is(err, &fs.PathError{}) {
		t.Errorf("\n got: %v \n want: %v", err, os.ErrNotExist)
	}
}

func TestNewGoogleStorageInvalidJSON(t *testing.T) {
	bn := "bucketName"
	os.Setenv(BucketEnvParamName, bn)
	os.Setenv(CredJSONEnvParamName, "test-data/test_cred_invalid.json")
	_, err := GoogleBucket(nil)
	if err == nil || !strings.HasPrefix(err.Error(), ClientCreationErrPrefix) {
		t.Errorf("\n got: %v \n want: %s: ...", err, ClientCreationErrPrefix)
	}
}

func TestNewGoogleStorageValid(t *testing.T) {
	bn := "bucketName"
	os.Setenv(BucketEnvParamName, bn)
	os.Setenv(CredJSONEnvParamName, "test-data/test_cred.json")
	g, err := GoogleBucket(nil)
	if g == nil || err != nil {
		t.Errorf("\n The GoogleStorage should have been created got: %v, %v \n want: %v, %v", g, err, "a not nil result", nil)
	}
}

func Test_apply(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want Config
	}{
		{
			name: "empty options",
			args: args{options: nil},
			want: Config{defaultBucket: "", defaultTimeout: 2000, defaultCredJSONPath: "", encoder: DefaultEncoder, decoder: DefaultDecoder},
		},
		{
			name: "setting default timeout",
			args: args{options: []Option{WithDefaultTimeout(3000)}},
			want: Config{defaultBucket: "", defaultTimeout: 3000, defaultCredJSONPath: "", encoder: DefaultEncoder, decoder: DefaultDecoder},
		},
		{
			name: "setting default timeout, and default bucket",
			args: args{options: []Option{WithDefaultTimeout(3000), WithDefaultBucket("testBucket")}},
			want: Config{defaultBucket: "testBucket", defaultTimeout: 3000, defaultCredJSONPath: "", encoder: DefaultEncoder, decoder: DefaultDecoder},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := apply(tt.args.options...); !ConfigEqual(got, tt.want) {
				t.Errorf("apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ConfigEqual(got, want Config) bool {
	ret := true
	ret = ret && (got.defaultBucket == want.defaultBucket)
	ret = ret && (got.defaultTimeout == want.defaultTimeout)
	ret = ret && (got.defaultCredJSONPath == want.defaultCredJSONPath)
	ret = ret && (got.fileSuffix == want.fileSuffix)
	// ret = ret && (&got.encoder == &want.encoder)
	// ret = ret && (&got.decoder == &want.decoder)
	return ret
}
