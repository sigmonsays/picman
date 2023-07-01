package core

import "fmt"

type ExifData struct {
	Values map[string]interface{}
}

func (me *ExifData) KeyExists(key string) bool {
	_, ok := me.Values[key]
	if !ok {
		return false
	}
	return true
}

// return first key to exist
func (me *ExifData) FindFirst(keys ...string) string {
	for _, k := range keys {
		if me.KeyExists(k) {
			return k
		}
	}
	return ""
}

func (me *ExifData) GetString(key string) (string, error) {
	v, ok := me.Values[key]
	if !ok {
		return "", fmt.Errorf("No such key: %s", key)
	}
	ret, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("key %s: not a string", key)
	}
	return ret, nil
}

func (me *ExifData) GetInt(key string) (int64, error) {
	v, ok := me.Values[key]
	if !ok {
		return 0, fmt.Errorf("No such key: %s", key)
	}
	ret, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("key %s: not a int", key)
	}
	return int64(ret), nil
}
