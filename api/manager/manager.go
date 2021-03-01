package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"git.code.oa.com/tme-server-component/kg_growth_open/api/sdk"
)

var (
	instance *manager
)

func init() {
	instance = &manager{
		impl:           make(map[string]sdk.MarketingSDK),
	}
}

type manager struct {
	impl           map[string]sdk.MarketingSDK
	serveMux       *http.ServeMux
}

// Register 注册对应平台的实现
func Register(platform string, impl sdk.MarketingSDK) {
	instance.impl[platform] = impl
}

// GetImpl 获取对应平台的实现
func GetImpl(platform string) sdk.MarketingSDK {
	return instance.impl[platform]
}

// Call 调用对应的方法, 方便web直接传入string，直接调用
func Call(platform string, method string, input string) (string, error) {
	if impl, ok := instance.impl[platform]; !ok {
		return "", fmt.Errorf("platform[%s] not register", platform)
	} else {
		return call(impl, method, input)
	}
}

func call(impl sdk.MarketingSDK, method string, input string) (string, error) {
	mt := reflect.ValueOf(impl).MethodByName(method)

	if mt.Type().NumIn() != 1 {
		return "", fmt.Errorf("method argument length must be 1")
	}

	inputType := mt.Type().In(0)
	kind := inputType.Kind()
	if kind == reflect.Ptr {
		inputType = inputType.Elem()
	}

	inputStruct := reflect.New(inputType)
	if err := json.Unmarshal([]byte(input), inputStruct.Interface()); err != nil {
		return "", err
	}
	if kind != reflect.Ptr {
		inputStruct = inputStruct.Elem()
	}

	values := mt.Call([]reflect.Value{inputStruct})
	if len(values) != 2 {
		return "", fmt.Errorf("method output length not 2")
	}

	if !values[1].IsNil() {
		return "", values[1].Interface().(error)
	}

	output, err := json.Marshal(values[0].Interface())
	if err != nil {
		return "", err
	}

	return string(output), nil
}
