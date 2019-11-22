package txcos

type options struct {
	secretId   string
	secretKey  string
	region     string
	bucketName string
}

type Option func(opts *options)

func WithSecretIdKey(id, key string) Option {
	return func(opts *options) {
		opts.secretId = id
		opts.secretKey = key
	}
}

func WithRegion(data string) Option {
	return func(opts *options) {
		opts.region = data
	}
}

func WithBucketName(data string) Option {
	return func(opts *options) {
		opts.bucketName = data
	}
}
