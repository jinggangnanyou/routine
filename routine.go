package routine

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var goroutineLocalStorage sync.Map

func Set(inV map[string]interface{}) {
	if inV == nil {
		return
	}
	goroutineID := Goid()
	if oldV, exists := goroutineLocalStorage.Load(goroutineID); exists {
		tmp := oldV.(map[string]interface{})
		for k, v := range inV {
			tmp[k] = v
		}
		goroutineLocalStorage.Store(goroutineID, tmp)
	} else {
		goroutineLocalStorage.Store(goroutineID, inV)
	}
}

func Get() (map[string]interface{}, bool) {
	goroutineID := Goid()
	v, ok := goroutineLocalStorage.Load(goroutineID)
	if ok {
		if v != nil {
			return v.(map[string]interface{}), ok
		} else {
			return nil, ok
		}
	} else {
		return nil, ok
	}
}

func Del() {
	goroutineID := Goid()
	goroutineLocalStorage.Delete(goroutineID)
}

func Goid() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)
	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return -1
	}
	return int64(id)
}

func deepCopyMap(originalMap map[string]interface{}) map[string]interface{} {
	if originalMap == nil {
		return nil
	}
	newMap := make(map[string]interface{})
	for key, value := range originalMap {
		valueType := reflect.TypeOf(value)
		if valueType.Kind() == reflect.Map {
			nestedMap := value.(map[string]interface{})
			newMap[key] = deepCopyMap(nestedMap)
		} else {
			newMap[key] = value
		}
	}
	return newMap
}

func Copy() map[string]interface{} {
	copy := make(map[string]interface{})
	oldV, exists := Get()
	if exists {
		copy = deepCopyMap(oldV)
	}
	return copy
}

func Inherit(in map[string]interface{}) {
	Set(in)
}

func Goroutine(f func()) {
	oldV := Copy()
	go func() {
		defer Del()
		Inherit(oldV)
		f()
	}()
}
