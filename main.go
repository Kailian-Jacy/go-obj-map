package objDict

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"time"

	cache "github.com/patrickmn/go-cache"
)

var tagNames = []string{
	"none",
	"ignore",
}

type objDict struct {
	c *cache.Cache
}

type objIndex interface {
	interface{}
}

type objMeta struct {
	basicType  reflect.Type
	memberType []string
	memberName []string
	memberTags []string
}

func NewDict() objDict {
	o := cache.New(5*time.Minute, 10*time.Minute)
	return objDict{o}
}

func (o objDict) Set(k objIndex, v interface{}) {
	meta := getMeta(k)
	o.c.Set(meta.Hash(), v, cache.DefaultExpiration)
}

func (o objDict) Get(k objIndex) (interface{}, bool) {
	return o.c.Get(getMeta(k).Hash())
}

func getMeta(idx objIndex) (meta objMeta) {
	if idx == nil {
		return meta
	}
	meta.basicType = reflect.TypeOf(idx)
	switch meta.basicType.Kind() {
	case reflect.Struct:
		return structMeta(idx, meta)
	}
	return
}

func structMeta(idx objIndex, meta objMeta) objMeta {
	for i := 0; i < meta.basicType.NumField(); i++ {
		typeField := reflect.TypeOf(idx).Field(i)
		meta.memberType = append(meta.memberType, typeField.Type.String())
		meta.memberName = append(meta.memberName, typeField.Name)
		meta.memberTags = append(meta.memberTags, getTags(typeField))
	}
	return meta
}

func getTags(tf reflect.StructField) string {
	// TODO: Build Bitmap
	for _, str := range tagNames {
		if name, ok := tf.Tag.Lookup(str); ok {
			return name
		}
	}
	return "none"
}

func (meta objMeta) Hash() string {
	var metaStr string
	for idx := range meta.memberType {
		metaStr += meta.memberType[idx]
		metaStr += "+"
		metaStr += meta.memberName[idx]
		metaStr += "%"
	}
	h := fnv.New32a()
	h.Write([]byte(metaStr))
	return fmt.Sprint(h.Sum32())
}
