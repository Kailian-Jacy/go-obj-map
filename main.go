package objDict

import (
	"fmt"
	"hash/fnv"
	"reflect"
	"strings"
	"time"

	"github.com/boljen/go-bitmap"
	"github.com/patrickmn/go-cache"
)

const (
	ignore int = 1
)

type ObjDict struct {
	c *cache.Cache
}

type objIndex interface {
	interface{}
}

type objMeta struct {
	basicType   reflect.Type
	memberType  []string
	memberName  []string
	memberTags  []bitmap.Bitmap
	memberValue []interface{}
}

func New() ObjDict {
	o := cache.New(5*time.Minute, 10*time.Minute)
	return ObjDict{o}
}

func (o ObjDict) Set(k objIndex, v interface{}) {
	meta := getMeta(k)
	o.c.Set(meta.Hash(), v, cache.DefaultExpiration)
}

func (o ObjDict) Get(k objIndex) (interface{}, bool) {
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
	// TODO: Support other types.
	default:
		return meta
	}
}

func getTags(tf reflect.StructField) bitmap.Bitmap {
	// TODO: Build Bitmap
	bm := bitmap.New(1)
	if str, ok := tf.Tag.Lookup("objmap"); ok {
		if strings.Contains(str, "ignore") {
			bm.Set(ignore, true)
		}
	}
	return bm
}

func (meta objMeta) Hash() string {
	// TODO: promote hash func:
	// Change hash type for less collision.
	// Promote MetaStr composer for better performance and better support for special chars.
	var metaStr string
	for idx, tag := range meta.memberTags {
		if tag.Get(ignore) {
			continue
		}
		metaStr += meta.memberType[idx]
		metaStr += "+"
		metaStr += meta.memberName[idx]
		metaStr += "%"
		metaStr += meta.memberName[idx]
		metaStr += "?"
		metaStr += fmt.Sprint(meta.memberValue[idx])
		metaStr += "&"
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(metaStr))
	return fmt.Sprint(h.Sum32())
}
