package alioss

type options struct {
	secretId   string
	secretKey  string
	endpoint   string
	bucketName string
}

type Option func(*options)

func WithSecretIdKey(id, key string) Option {
	return func(opts *options) {
		opts.secretId = id
		opts.secretKey = key
	}
}

func WithEndpoint(data string) Option {
	return func(o *options) {
		o.endpoint = data
	}
}

func WithBucketName(data string) Option {
	return func(o *options) {
		o.bucketName = data
	}
}
