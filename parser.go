package objDict

import "reflect"

func structMeta(idx objIndex, meta objMeta) objMeta {
	for i := 0; i < meta.basicType.NumField(); i++ {
		typeField := reflect.TypeOf(idx).Field(i)
		meta.memberType = append(meta.memberType, typeField.Type.String())
		meta.memberName = append(meta.memberName, typeField.Name)
		meta.memberTags = append(meta.memberTags, getTags(typeField))
		meta.memberValue = append(meta.memberValue, reflect.ValueOf(idx).Field(i).Interface())
	}
	return meta
}
