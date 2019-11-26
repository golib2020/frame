package f

import (
	"fmt"

	"github.com/golib2020/frame/storage"
)

const (
	storageInstancesPrefix = `storage`
	storageDefaultGroup    = `default`
)

//Storage 持久存储
func Storage(name ...string) storage.Storage {
	group := storageDefaultGroup
	if len(name) > 0 && name[0] != "" {
		group = name[0]
	}
	key := fmt.Sprintf("%s.%s", storageInstancesPrefix, group)
	return Instance().GetOrSetFunc(key, func() interface{} {
		return storageInit(key)
	}).(storage.Storage)
}

func storageInit(key string) storage.Storage {
	conf := Config().Sub(key)
	var s storage.Storage
	switch conf.GetString("driver") {
	case "local":
		conf.SetDefault("root", "/storage/")
		conf.SetDefault("host", "/")
		s = storage.NewLocal(
			conf.GetString("root"),
			conf.GetString("host"),
		)
	case "oss":
		s = storage.NewAliOSS(
			conf.GetString("root"),
			conf.GetString("host"),
			storage.WithSecretIdKey(conf.GetString("secret_id"), conf.GetString("secret_key")),
			storage.WithEndpoint(conf.GetString("endpoint")),
			storage.WithBucketName(conf.GetString("bucket_name")),
		)
	case "cos":
		s = storage.NewCOS(
			conf.GetString("root"),
			conf.GetString("host"),
			storage.WithSecretIdKey(conf.GetString("secret_id"), conf.GetString("secret_key")),
			storage.WithRegion(conf.GetString("region")),
			storage.WithBucketName(conf.GetString("bucket_name")),
		)
	}
	return s
}
