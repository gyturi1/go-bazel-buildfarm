package ittest

import (
	"bufio"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gyturi1/go-bazel-buildfarm/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestCallingWithFakeStorage(t *testing.T) {
	os.Setenv(storage.BucketEnvParamName, ReadParam(t, "docker-compose/.env", "BUCKET"))
	os.Setenv(storage.EmulatorEnvParamName, "localhost:"+ReadParam(t, "docker-compose/.env", "FAKESTORAGE_PORT"))
	store, err := storage.GoogleBucket(nil)
	if err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile("../test-data/test_cred.json")
	if err != nil {
		t.Fatal(err)
	}
	ctx, c := context.WithTimeout(context.Background(), 2000*time.Millisecond)
	defer c()
	names := []string{"name1", "name2", "name4", "name6"}
	for _, n := range names {
		err = store.Write(ctx, "ten1", n, b, "application.json")
		if err != nil {
			t.Fatal(err)
		}
	}

	r, err := store.List(ctx, "ten1")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, names, r)
}

func ReadParam(t *testing.T, path string, paramName string) string {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		splitted := strings.Split(scanner.Text(), "=")
		if splitted[0] == paramName {
			return splitted[1]
		}
	}
	t.Logf("No matching entry found for: %s in %s", paramName, path)
	return ""
}
