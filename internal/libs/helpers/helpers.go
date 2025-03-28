package helpers

import (
	"runtime"
)

const (
	DefaultDepth = 1
	DefaultIndex = 0
)

// GetFunctionName Возвращает имяПакета.ИмяФункции.
func GetFunctionName(depthList ...int) string { //nolint:unused // helper func
	var depth int

	if depthList == nil {
		depth = DefaultDepth
	} else {
		depth = depthList[DefaultIndex]
	}

	function, _, _, ok := runtime.Caller(depth)
	if !ok {
		return "Не удалось получить имя функции"
	}

	return runtime.FuncForPC(function).Name()
}
