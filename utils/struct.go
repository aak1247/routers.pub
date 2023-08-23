package utils

import (
	"reflect"
	"strings"
)

type CanHasError struct {
	Errors []error `gorm:"-" json:"-"`
}

func (e *CanHasError) HasError() bool {
	return e.Errors != nil && len(e.Errors) != 0
}

func (e *CanHasError) AddError(err error) {
	if e.Errors == nil {
		e.Errors = make([]error, 0)
	}
	e.Errors = append(e.Errors, err)
}

func (e *CanHasError) AddErrors(errs []error) {
	if e.Errors == nil {
		e.Errors = make([]error, 0)
	}
	e.Errors = append(e.Errors, errs...)
}

func (e *CanHasError) GetErrors() []error {
	return e.Errors
}

func (e *CanHasError) GetError() error {
	if !e.HasError() {
		return nil
	}
	return e.Errors[0]
}

func GetColumnNameOfDef(name string, tag string) string {
	// 从定义中获取列名
	var columnName = ""
	if tag != "" {
		columnName = GetColumnNameFromTag(tag)
	}
	if columnName == "" {
		columnName = CamelToSnake(name)
	}
	return columnName
}

func GetColumnNameFromTag(tag string) string {
	// 从tag中获取列名
	var columnName = ""
	tags := strings.Split(tag, ";")
	for _, t := range tags {
		if strings.HasPrefix(t, "column:") && len(t) > 7 {
			columnName = strings.Split(t, ":")[1]
		}
	}
	return columnName
}

// Patch语义的更新，只更新不为nil的字段（不需要实体属性为指针）
// valueDef: 值定义的结构体(是结构体不是指针)
// entity: 实体结构体(是结构体不是指针)
func GetUpdatesByDefAndEntity(valueDef any, entity any) map[string]interface{} {
	updates := make(map[string]interface{})
	// 通过反射遍历valueDef，如果不为nil则更新
	f := reflect.TypeOf(valueDef)
	t_e := reflect.TypeOf(entity)
	v := reflect.ValueOf(valueDef)
	for i := 0; i < f.NumField(); i++ {
		// 如果为nil则跳过(其他类型字段的零值也会被跳过)
		if v.Field(i).IsZero() {
			continue
		}
		name := f.Field(i).Name
		tag := f.Field(i).Tag.Get("gorm")
		colName := GetColumnNameOfDef(name, tag)
		if type_t_e, ok := t_e.FieldByName(name); ok && // 如果目标实体中存在该字段
			type_t_e.Type != f.Field(i).Type && // 且类型不一致
			f.Field(i).Type.Kind() == reflect.Pointer && // 且valueDef中为指针类型
			type_t_e.Type == f.Field(i).Type.Elem() { // 且目标实体中的类型为实体中类型的对应指针类型
			updates[colName] = v.Field(i).Elem().Interface() // 则借引用后用值更新目标实体中的元素
		} else if ok {
			if !f.Field(i).Type.AssignableTo(type_t_e.Type) &&
				!f.Field(i).Type.ConvertibleTo(type_t_e.Type) {
				// 如果目标实体中的值类型不可赋值，则跳过该字段
				continue
			}
			updates[colName] = v.Field(i).Interface() // 否则更新目标实体中的字段
		}
	}
	return updates
}

// Patch语义的更新，只更新不为nil的字段（不需要实体属性为指针）
// entity: 实体结构体(是结构体不是指针)
// valueDef: 值定义的结构体(是结构体不是指针)
func UpdateEntityByDef[T any](entity T, valueDef interface{}) T {
	// 通过反射遍历valueDef，如果不为nil则更新
	f := reflect.TypeOf(valueDef)
	t_e := reflect.TypeOf(entity)
	v := reflect.ValueOf(valueDef)
	v_e := reflect.ValueOf(&entity).Elem()
	for i := 0; i < f.NumField(); i++ {
		// 如果为nil则跳过(其他类型字段的零值也会被跳过)
		if v.Field(i).IsZero() {
			continue
		}
		name := f.Field(i).Name
		if type_te, ok := t_e.FieldByName(name); ok && // 如果目标实体中存在该字段
			type_te.Type != f.Field(i).Type && // 且类型不一致
			f.Field(i).Type.Kind() == reflect.Pointer && // 且valueDef中为指针类型
			type_te.Type == f.Field(i).Type.Elem() { // 且目标类型为实体中类型的对应指针类型
			// 则解引用后用值更新目标实体中的值
			v_e.FieldByName(name).Set(v.Field(i).Elem())
		} else if ok && v_e.FieldByName(name).CanSet() {
			switch v_e.FieldByName(name).Kind() {
			case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Slice:
				if v_e.FieldByName(name).IsNil() {
					// 如果目标实体中的值为nil，则需要先初始化
					v_e.FieldByName(name).Set(reflect.New(type_te.Type.Elem()))
				}
			}
			if !f.Field(i).Type.AssignableTo(type_te.Type) {
				if f.Field(i).Type.ConvertibleTo(type_te.Type) {
					// 如果目标实体中的值类型不可赋值，则跳过该字段
					v_e.FieldByName(name).Set(v.Field(i).Convert(type_te.Type))
				}
				continue
			}
			// set entity中的字段
			v_e.FieldByName(name).Set(v.Field(i))
		}
	}
	return entity
}
