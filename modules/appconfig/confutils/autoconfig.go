package confutils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"github.com/wup364/pakku/ipakku"
	"github.com/wup364/pakku/utils/logs"
	"github.com/wup364/pakku/utils/reflectutil"
	"github.com/wup364/pakku/utils/utypes"
)

// NewAutoValueOfBeanUtil 自动配置工具
func NewAutoValueOfBeanUtil(config ipakku.IConfig) *AutoValueOfBeanUtil {
	return &AutoValueOfBeanUtil{config: config}
}

// AutoValueOfBeanUtil 自动配置工具
type AutoValueOfBeanUtil struct {
	config ipakku.IConfig
}

// AutoValueOfBean 根据struct的描述配置自动初始化值
func (av *AutoValueOfBeanUtil) AutoValueOfBean(ptr interface{}) (err error) {
	var fieldVals map[string]string
	if fieldVals, err = reflectutil.GetTagValues(ipakku.PAKKUTAG_AUTOCONFIG, ptr); nil != err || len(fieldVals) == 0 {
		return
	}
	// 仅支持指针类型结构体
	if t := reflect.TypeOf(ptr); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer objects are supported")
	}
	// 获得配置类
	for field, cprefix := range fieldVals {
		var fvalue reflect.Value
		if fvalue, err = reflectutil.GetStructFieldRefValue(ptr, field); nil != err {
			logs.Infof("> AutoConfig %s [err=%s] \r\n", field, err.Error())
			return
		}
		logs.Infof("> AutoConfig %s [%s] \r\n", field, fvalue.Type().String())
		if fvalue.Type().Kind() != reflect.Ptr {
			err = fmt.Errorf("only pointer objects are supported, field: %s", field)
			break
		}
		// 创建对象
		newValue := reflect.New(fvalue.Type().Elem())
		if err = av.setBeanValue(cprefix, newValue); nil != err {
			return
		}
		//
		reflect.NewAt(fvalue.Type(), unsafe.Pointer(fvalue.UnsafeAddr())).Elem().Set(newValue)
	}
	return err
}

// setBeanValue 结构赋值
func (av *AutoValueOfBeanUtil) setBeanValue(cprefix string, ptr reflect.Value) (err error) {
	var tagvals = make(map[string]string)
	if tagvals, err = reflectutil.GetTagValues(ipakku.PAKKUTAG_CONFIG_VALUE, ptr); nil != err || len(tagvals) == 0 {
		return
	}
	//
	for fieldName, confPath := range tagvals {
		defaultVal := ""
		configKey := confPath
		if dfValIndex := strings.Index(confPath, ":"); dfValIndex > -1 {
			configKey = confPath[:dfValIndex]
			defaultVal = strings.TrimSpace(confPath[dfValIndex+1:])
		}
		if len(configKey) == 0 {
			configKey = fieldName
		}
		if len(cprefix) > 0 {
			configKey = cprefix + "." + configKey
		}

		var confVal utypes.Object
		vv := ptr.Elem().FieldByName(fieldName)
		if confVal = av.config.GetConfig(configKey); confVal.IsNill() && len(defaultVal) > 0 {
			confVal = utypes.NewObject(defaultVal)
		}
		if err = av.setFeildValue(reflect.NewAt(vv.Type(), unsafe.Pointer(vv.UnsafeAddr())).Elem(), confVal); nil != err {
			return
		}
		logs.Debugf(">  setBeanValue %s <= %s[value=%v] \r\n", fieldName, configKey, vv)
	}
	return err
}

// setFeildValue 设置字段值
func (av *AutoValueOfBeanUtil) setFeildValue(v reflect.Value, o utypes.Object) error {
	obj := o.GetVal()
	if obj == nil {
		return nil
	}
	vKind := v.Type().Kind()
	if vKind == reflect.Int || vKind == reflect.Int8 || vKind == reflect.Int16 || vKind == reflect.Int32 || vKind == reflect.Int64 {
		v.SetInt(o.ToInt64(0))

	} else if vKind == reflect.Uint || vKind == reflect.Uint8 || vKind == reflect.Uint16 || vKind == reflect.Uint32 || vKind == reflect.Uint64 {
		v.SetUint(o.ToUint64(0))

	} else if vKind == reflect.Float32 || vKind == reflect.Float64 {
		v.SetFloat(o.ToFloat64(0))

	} else if vKind == reflect.String {
		v.SetString(o.ToString(""))

	} else if vKind == reflect.Bool {
		v.SetBool(o.ToBool(false))

	} else if vKind == reflect.Map {
		return ipakku.ErrUnsupported

	} else if vKind == reflect.Array {
		return ipakku.ErrUnsupported

	} else if vKind == reflect.Slice {
		if objType := reflect.TypeOf(obj); objType.Kind() == reflect.Slice {
			objKind := objType.Elem().Kind()
			vKind := v.Type().Elem().Kind()
			if objKind == vKind {
				v.Set(reflect.ValueOf(obj))
				return nil

			} else if objKind == reflect.Interface {
				return av.setFeildValue4Array(v, obj.([]interface{}))

			}
		}
	} else {
		return ipakku.ErrUnsupported
	}
	return nil
}

// setFeildValue4Array 设置字段值-数组类型
func (av *AutoValueOfBeanUtil) setFeildValue4Array(v reflect.Value, in []interface{}) error {
	if len(in) == 0 {
		return nil
	}
	obj := &utypes.Object{}
	vKind := v.Type().Elem().Kind()
	if vKind == reflect.Int {
		res := make([]int, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToInt(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Int8 {
		res := make([]int8, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToInt8(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Int16 {
		res := make([]int16, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToInt16(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Int32 {
		res := make([]int32, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToInt32(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Int64 {
		res := make([]int64, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToInt64(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Uint {
		res := make([]uint, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToUint(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Uint8 {
		res := make([]uint8, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToUint8(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Uint16 {
		res := make([]uint16, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToUint16(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Uint32 {
		res := make([]uint32, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToUint32(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Uint64 {
		res := make([]uint64, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToUint64(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Float32 {
		res := make([]float32, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToFloat32(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Float64 {
		res := make([]float64, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToFloat64(0)
		}
		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.String {
		res := make([]string, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToString("")
		}

		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Bool {
		res := make([]bool, len(in))
		for i := 0; i < len(in); i++ {
			res[i] = obj.SetVal(in[i]).ToBool(false)
		}

		v.Set(reflect.ValueOf(res))

	} else if vKind == reflect.Map {
		return ipakku.ErrUnsupported

	} else if vKind == reflect.Array {
		return ipakku.ErrUnsupported

	} else if vKind == reflect.Slice {
		return ipakku.ErrUnsupported

	} else if vKind == reflect.Interface {
		v.Set(reflect.ValueOf(in))

	} else {
		return ipakku.ErrUnsupported
	}

	return nil
}
