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

// ScanAndAutoConfig 扫描带有@autoconfig标签的字段, 并完成其配置
func (av *AutoValueOfBeanUtil) ScanAndAutoConfig(ptr interface{}) (err error) {
	// 仅支持指针类型结构体
	if t := reflect.TypeOf(ptr); t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer objects are supported")
	}

	return av.scanAndAutoConfig(ptr, "")
}

// ScanAndAutoValue 扫描带有@autovalue标签的字段, 并完成其配置
func (av *AutoValueOfBeanUtil) ScanAndAutoValue(cprefix string, ptr interface{}) (err error) {
	refVal := reflect.ValueOf(ptr)
	if refVal.Kind() != reflect.Pointer || refVal.Elem().Kind() != reflect.Struct {
		return errors.New("the input object must be a pointer struct")
	}
	if err = av.setBeanValue(cprefix, refVal); nil != err {
		return
	}
	return err
}

// scanAndAutoConfig 扫描自动配置类并配置
func (av *AutoValueOfBeanUtil) scanAndAutoConfig(ptr interface{}, prefix string) (err error) {
	var fieldVals map[string]string
	if fieldVals = reflectutil.GetTagValues(ipakku.PAKKUTAG_AUTOCONFIG, ptr); len(fieldVals) == 0 {
		return
	}

	// 获得配置类
	for field, cprefix := range fieldVals {
		if len(prefix) > 0 {
			cprefix = prefix + "." + cprefix
		}
		if err = av.doConfigField(ptr, cprefix, field); nil != err {
			return
		}
	}
	return err
}

// doConfigField 配置ptr内的某个字段
func (av *AutoValueOfBeanUtil) doConfigField(ptr interface{}, cprefix, fieldName string) (err error) {
	var fvalue reflect.Value
	if fvalue, err = reflectutil.GetStructFieldRefValue(ptr, fieldName); nil != err {
		logs.Infof("> AutoConfig %s [err=%s] \r\n", fieldName, err.Error())
		return
	}
	logs.Infof("> AutoConfig %s [%s] \r\n", fieldName, fvalue.Type().String())

	// 创建对象 & 赋值
	var newValue reflect.Value
	if fvalue.Type().Kind() == reflect.Ptr {
		newValue = reflect.New(fvalue.Type().Elem())
	} else {
		newValue = reflect.New(fvalue.Type())
	}
	if err = av.setBeanValue(cprefix, newValue); nil != err {
		return
	}

	// 回写值
	fValElem := reflect.NewAt(fvalue.Type(), unsafe.Pointer(fvalue.UnsafeAddr())).Elem()
	if fvalue.Type().Kind() == reflect.Ptr {
		fValElem.Set(newValue)
	} else {
		fValElem.Set(newValue.Elem())
	}
	return
}

// setBeanValue 结构赋值
func (av *AutoValueOfBeanUtil) setBeanValue(cprefix string, ptr reflect.Value) (err error) {
	var tagvals = make(map[string]string)
	if tagvals = reflectutil.GetTagValues(ipakku.PAKKUTAG_CONFIG_VALUE, ptr); len(tagvals) == 0 {
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

	if nil == err {
		// 继续扫描匿名类
		err = av.scanAndAutoConfigAnonymous(ptr, cprefix)
	}
	if nil == err {
		// 继续扫描嵌套的自动配置类
		err = av.scanAndAutoConfig(ptr, cprefix)
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
		return newUnsupportedTypeErr(v.Type())

	} else if vKind == reflect.Array {
		return newUnsupportedTypeErr(v.Type())

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
	} else if vKind == reflect.Struct {
		if v.Type() == reflect.TypeOf(o) {
			v.Set(reflect.ValueOf(o))
		} else {
			return newUnsupportedTypeErr(v.Type())
		}
	} else {
		return newUnsupportedTypeErr(v.Type())
	}
	return nil
}

// setFeildValue4Array 设置字段值-数组类型
func (av *AutoValueOfBeanUtil) setFeildValue4Array(v reflect.Value, in []interface{}) error {
	if len(in) == 0 {
		return nil
	}
	vType := v.Type().Elem()
	obj := &utypes.Object{}
	vKind := vType.Kind()
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
		return newUnsupportedTypeErr(vType)

	} else if vKind == reflect.Array {
		return newUnsupportedTypeErr(vType)

	} else if vKind == reflect.Slice {
		return newUnsupportedTypeErr(vType)

	} else if vKind == reflect.Interface {
		v.Set(reflect.ValueOf(in))

	} else {
		return fmt.Errorf("data types that do not support automatic configuration: %s", vType.Name())
	}

	return nil
}

// scanAndAutoConfigAnonymous 扫描匿名嵌套类并配置
func (av *AutoValueOfBeanUtil) scanAndAutoConfigAnonymous(ptr interface{}, cprefix string) (err error) {
	var fields []reflect.StructField
	if fields = reflectutil.GetAnonymousOrNoneTypeNameField(ptr); len(fields) == 0 {
		return
	}
	for i := 0; i < len(fields); i++ {
		if err = av.doConfigField(ptr, cprefix, fields[i].Name); nil != err {
			return
		}
	}
	return
}

func newUnsupportedTypeErr(k reflect.Type) error {
	return fmt.Errorf("data types that do not support automatic configuration: %s", k.Name())
}
