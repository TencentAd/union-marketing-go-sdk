package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/tencentad/union-marketing-go-sdk/api/sdk"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/ams"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/config"
	"github.com/tencentad/union-marketing-go-sdk/pkg/sdk/ocean_engine"
)

var (
	instance *manager
)

func init() {
	instance = &manager{
		impl: make(map[sdk.MarketingPlatformType]sdk.MarketingSDK),
	}
}

type manager struct {
	impl     map[sdk.MarketingPlatformType]sdk.MarketingSDK
	serveMux *http.ServeMux
}

// Register 注册对应平台的实现
func Register(platform sdk.MarketingPlatformType, config *config.Config) error {
	switch platform {
	case sdk.AMS:
		instance.impl[platform] = ams.NewAMSService(config)
		return nil
	case sdk.OceanEngine:
		instance.impl[platform] = ocean_engine.NewOceanEngineService(config)
		return nil
	default:
		return fmt.Errorf("not support platform = %s", platform)
	}
}

// GetPlatformList 获取注册的平台列表
func GetPlatformList() []sdk.MarketingPlatformType {
	keys := make([]sdk.MarketingPlatformType, 0, len(instance.impl))
	for k := range instance.impl {
		keys = append(keys, k)
	}
	return keys
}

// GetImpl 获取对应平台的实现
func GetImpl(platform sdk.MarketingPlatformType) (sdk.MarketingSDK, error) {
	if instance.impl[platform] == nil {
		return nil, fmt.Errorf("not register platform = %s", platform)
	}
	return instance.impl[platform], nil
}

// Call 调用对应的方法, 方便web直接传入string，直接调用
func Call(platform sdk.MarketingPlatformType, method string, input string) (string, error) {
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
