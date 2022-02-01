package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/rs/zerolog"
)

const (
	EmulatorEnvParamName    = "STORAGE_EMULATOR_HOST"
	CredJSONEnvParamName    = "CLOUD_STORAGE_CREDENTIALS_JSON" //nolint:gosec
	TimeoutEnvParamName     = "TIMEOUT_MILLISECOND"
	BucketEnvParamName      = "CLOUD_STORAGE_BUCKET"
	ClientCreationErrPrefix = "could not create storage client"
)

var (
	ErrNoBucket       = fmt.Errorf("empty bucket name, please specify bucket with %s", BucketEnvParamName)
	ErrNoCredFilePath = fmt.Errorf("empty path, please specify google API cradantial json file path with %s", CredJSONEnvParamName)
)

type Config struct {
	defaultBucket       string
	defaultTimeout      int
	defaultCredJSONPath string
	fileSuffix          string
	encoder             FilenameEncoder
	decoder             FilenameDecoder
}

type Option = func(*Config)

// use this as the dafult timeout if no setting from env parameter TIMEOUT_MILLISECOND.
func WithDefaultTimeout(milliseconds int) Option {
	return func(sc *Config) {
		sc.defaultTimeout = milliseconds
	}
}

// use this as the default bucket name if no seetings from env parameter CLOUD_STORAGE_BUCKET.
func WithDefaultBucket(b string) Option {
	return func(sc *Config) {
		sc.defaultBucket = b
	}
}

// Filename encoder and decoder.
func WithFileSuffix(s string) Option {
	return func(sc *Config) {
		sc.fileSuffix = s
	}
}

// the filename encoder, usually on google sotrage u want to encode the filename with slash so google storage rander them as a directory structure.
func WithEncoder(f FnEncoder) Option {
	return func(sc *Config) {
		sc.encoder = f
	}
}

// the opposite of filename encoder.
func WithDecoder(f FnDecoder) Option {
	return func(sc *Config) {
		sc.decoder = f
	}
}

func apply(options ...Option) Config {
	conf := &Config{
		defaultBucket:       "",
		defaultCredJSONPath: "",
		defaultTimeout:      2000,
		fileSuffix:          "",
		encoder:             DefaultEncoder,
		decoder:             DefaultDecoder,
	}
	for _, o := range options {
		o(conf)
	}
	return *conf
}

type GoogleStorage struct {
	client          *storage.Client
	bucket          string
	timeout         time.Duration
	logger          *zerolog.Logger
	fileSuffix      string
	fileNameEncoder FilenameEncoder
	filenameDecoder FilenameDecoder
}

// Creates a google storage client with the specified options. Env params (defaults if not specified as an option): CLOUD_STORAGE_CREDENTIALS_JSON (""), TIMEOUT_MILLISECOND (2000), CLOUD_STORAGE_BUCKET ("")
// for testing use: STORAGE_EMULATOR_HOST
func GoogleBucket(logger *zerolog.Logger, o ...Option) (Storage, error) {
	ctx := context.Background()
	sc := apply(o...)
	var (
		credJSONFilePath = ""
		b                = sc.defaultBucket
		timeoutMillis    = sc.defaultTimeout
	)

	setIfPresentString(&credJSONFilePath, CredJSONEnvParamName)
	setIfPresentString(&b, BucketEnvParamName)
	setIfPresentInt(&timeoutMillis, TimeoutEnvParamName)

	if len(b) == 0 {
		return nil, ErrNoBucket
	}

	var opt []option.ClientOption = []option.ClientOption{}
	if isPresent(EmulatorEnvParamName) {
		opt = append(opt, option.WithoutAuthentication())
	} else {
		if len(credJSONFilePath) == 0 {
			return nil, ErrNoCredFilePath
		}
		_, err := os.Stat(credJSONFilePath)
		if err != nil {
			return nil, err
		}
		opt = append(opt, option.WithCredentialsFile(credJSONFilePath))
	}

	client, err := storage.NewClient(ctx, opt...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", ClientCreationErrPrefix, err)
	}

	t := time.Duration(timeoutMillis) * time.Millisecond

	return &GoogleStorage{client, b, t, logger, sc.fileSuffix, sc.encoder, sc.decoder}, nil
}

func (g *GoogleStorage) Read(ctx context.Context, tenant string, name string) (res FileContent, err error) {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	f, err := g.fileNameEncoder(tenant, name, g.fileSuffix)
	if err != nil {
		return nil, err
	}

	r, err := g.client.Bucket(g.bucket).Object(f).NewReader(ctx)
	if err != nil {
		return nil, err
	}

	res, err = ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	err = r.Close()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (g *GoogleStorage) Query(ctx context.Context, tenant string, name string) (res []string, err error) {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	prefix := tenant + Separator + name
	i := g.client.Bucket(g.bucket).Objects(ctx, &storage.Query{Prefix: prefix})
	for {
		oa, err := i.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		res = append(res, oa.Name)
	}
	return res, nil
}

func (g *GoogleStorage) Write(ctx context.Context, tenant string, name string, settings FileContent, contentType string) error {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	f, err := g.fileNameEncoder(tenant, name, g.fileSuffix)
	if err != nil {
		return err
	}
	o := g.client.Bucket(g.bucket).Object(f)
	w := o.NewWriter(ctx)

	w.ContentType = contentType

	_, err = w.Write(settings)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (g *GoogleStorage) List(ctx context.Context, tenant string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()
	oi := g.client.Bucket(g.bucket).Objects(ctx, &storage.Query{Prefix: tenant})

	var names []string
	for {
		attrs, err := oi.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		f, err := g.filenameDecoder(attrs.Name, g.fileSuffix)
		if err != nil {
			return nil, err
		}
		names = append(names, f)
	}
	return names, nil
}

func (g *GoogleStorage) Delete(ctx context.Context, tenant string, name string) error {
	ctx, cancel := context.WithTimeout(ctx, g.timeout)
	defer cancel()

	f, err := g.fileNameEncoder(tenant, name, g.fileSuffix)
	if err != nil {
		return err
	}
	return g.client.Bucket(g.bucket).Object(f).Delete(ctx)
}

func setIfPresentString(field *string, envParam string) {
	if s, present := os.LookupEnv(envParam); present {
		(*field) = s
	}
}

func setIfPresentInt(field *int, envParam string) {
	if s, present := os.LookupEnv(envParam); present {
		i, err := strconv.Atoi(s)
		if err != nil {
			panic(fmt.Errorf("could not parse int from: %s, err: %w", envParam, err))
		}
		(*field) = i
	}
}

func isPresent(envParam string) bool {
	_, present := os.LookupEnv(envParam)
	return present
}
