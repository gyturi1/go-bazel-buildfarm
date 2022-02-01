package storage

import (
	"context"
	"fmt"
	"strings"
)

const (
	Separator = "/"
)

var (
	ErrTenantEmpty   = fmt.Errorf("tenant can not be empty")
	ErrNameEmpty     = fmt.Errorf("name can not be empty")
	ErrFileNameEmpty = fmt.Errorf("filename can not be empty, nothing to decode")
	ErrFileName      = func(f string) error {
		return fmt.Errorf("filename [%s] not conforms to '{tenant}/{name}.json' format", f)
	}
)

type (
	// encodes the filename for the storage accepts tenant, name and fileSuffix
	FnEncoder = func(string, string, string) (string, error)
	// deocdes the filename, must be the opposite of FnEncoder
	FnDecoder = func(string, string) (string, error)
)

type (
	FileContent     []byte
	FilenameEncoder FnEncoder
	FilenameDecoder FnDecoder
)

type Storage interface {
	Query(ctx context.Context, tenant, name string) ([]string, error)
	Read(ctx context.Context, tenant, name string) (FileContent, error)
	Write(ctx context.Context, tenant, name string, content FileContent, contentType string) error
	List(ctx context.Context, tenant string) ([]string, error)
	Delete(ctx context.Context, tenant, name string) error
}

var (
	DefaultEncoder = encodeFileName
	DefaultDecoder = decodeFileName
)

func encodeFileName(t, n, fileSuffix string) (string, error) {
	if len(t) == 0 {
		return "", ErrTenantEmpty
	}
	if len(n) == 0 {
		return "", ErrNameEmpty
	}
	return t + Separator + n + fileSuffix, nil
}

func decodeFileName(f, fileSuffix string) (string, error) {
	if len(f) == 0 {
		return "", ErrFileNameEmpty
	}
	x := strings.Replace(f, fileSuffix, "", 1)
	s := strings.Split(x, Separator)
	if len(s) == 2 {
		return s[1], nil
	}
	return "", ErrFileName(f)
}
