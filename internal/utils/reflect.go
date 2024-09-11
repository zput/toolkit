package utils

import "reflect"

func RemoveTags(input interface{}) interface{} {
	origValue := reflect.ValueOf(input)
	origType := origValue.Type()

	// 检查输入是否是指针
	if origType.Kind() == reflect.Ptr {
		origValue = origValue.Elem()
		origType = origValue.Type()
	}

	if origType.Kind() != reflect.Struct {
		panic("input must be a struct or a pointer to a struct")
	}

	// 创建一个新的结构体类型，没有tags
	numFields := origType.NumField()
	fields := make([]reflect.StructField, numFields)
	for i := 0; i < numFields; i++ {
		origField := origType.Field(i)
		fields[i] = reflect.StructField{
			Name: origField.Name,
			Type: origField.Type,
		}
	}

	// 使用新的结构体类型创建一个新的实例
	newType := reflect.StructOf(fields)
	newValue := reflect.New(newType).Elem()

	// 复制原始结构体的字段值
	for i := 0; i < numFields; i++ {
		newValue.Field(i).Set(origValue.Field(i))
	}

	// 返回指针类型的新结构体实例
	return newValue.Addr().Interface()
}
