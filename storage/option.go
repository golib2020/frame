package storage

type options struct {
	secretId   string
	secretKey  string
	bucketName string
	endpoint   string //oss
	region     string //cos
}

type Option func(opts *options)

func WithSecretIdKey(id, key string) Option {
	return func(opts *options) {
		opts.secretId = id
		opts.secretKey = key
	}
}

func WithBucketName(data string) Option {
	return func(opts *options) {
		opts.bucketName = data
	}
}

func WithEndpoint(data string) Option {
	return func(o *options) {
		o.endpoint = data
	}
}

func WithRegion(data string) Option {
	return func(opts *options) {
		opts.region = data
	}
}
