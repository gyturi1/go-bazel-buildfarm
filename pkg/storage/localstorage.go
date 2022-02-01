package storage

import (
	"context"
	"os"
	"path"
)

type LocalStorage struct {
	dir string
}

// It will store the files in a local FS directroy specified by env param (default): STORAGE_DIR (/tmp)
func LocalFS() Storage {
	dir := "/tmp"
	setIfPresentString(&dir, "STORAGE_DIR")
	return &LocalStorage{dir}
}

func (l *LocalStorage) Read(ctx context.Context, tenant string, name string) (FileContent, error) {
	n := path.Join(l.dir, tenant, name)
	return os.ReadFile(n)
}

func (l *LocalStorage) Query(ctx context.Context, tenant string, name string) ([]string, error) {
	n := path.Join(l.dir, tenant, name)
	es, err := os.ReadDir(n)
	if err != nil {
		return nil, err
	}
	var ret []string
	for _, e := range es {
		ret = append(ret, e.Name())
	}
	return ret, nil
}

func (l *LocalStorage) Write(ctx context.Context, tenant string, name string, settings FileContent, contentType string) error {
	d := path.Join(l.dir, tenant)
	err := os.MkdirAll(d, 0o700)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(d, name))
	if err != nil {
		return err
	}

	i, err := f.Write(settings)
	if err != nil || i != len(settings) {
		return err
	}
	return nil
}

func (l *LocalStorage) List(ctx context.Context, tenant string) ([]string, error) {
	d := path.Join(l.dir, tenant)
	fs, err := os.ReadDir(d)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, f := range fs {
		if !f.IsDir() {
			n := f.Name()
			res = append(res, n)
		}
	}

	return res, nil
}

func (l *LocalStorage) Delete(ctx context.Context, tenant string, name string) error {
	n := path.Join(l.dir, tenant, name)
	return os.Remove(n)
}
