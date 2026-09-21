package types

import (
	"log"
	"reflect"
)

var typeRegistry = make(map[string]reflect.Type)

func RegisterType(name string, datatype interface{}) {
	// Note: Cannot use logger here due to import cycle (logger imports types)
	log.Println("Registering type name:", name, reflect.TypeOf(datatype))
	typeRegistry[name] = reflect.TypeOf(datatype)
}

func CreateType(name string) interface{} {
	t := reflect.New(typeRegistry[name]).Interface()
	return t
}

func CreateSlice(name string) interface{} {
	tr := typeRegistry[name]
	if tr == nil {
		// Note: Cannot use logger here due to import cycle (logger imports types)
		log.Printf("cannot find type '%s' in registry: %v", name, typeRegistry)
		return nil
	}

	t := reflect.New(reflect.SliceOf(tr)).Interface()
	return t
}

func GetTypeNames() []string {
	names := make([]string, 0)
	for k, _ := range typeRegistry {
		names = append(names, k)
	}

	return names
}
