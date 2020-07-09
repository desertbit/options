/*
 * ORBIT - Interlink Remote Applications
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Sebastian Borchers <sebastian[at]desertbit.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package options

import (
	"errors"
	"reflect"
)

func SetDefaults(v, def interface{}) (err error) {
	// Ensure the type of v is a struct pointer.
	vType := reflect.TypeOf(v)
	if vType == nil || vType.Kind() != reflect.Ptr {
		return errors.New("value must be a struct pointer")
	}
	vVal := reflect.ValueOf(v).Elem()
	if vVal.Kind() != reflect.Struct {
		return errors.New("value must be a struct pointer")
	}

	// Ensure the type of def is a struct.
	defType := reflect.TypeOf(def)
	if defType == nil || defType.Kind() != reflect.Struct {
		return errors.New("def must be a struct value")
	}

	// Ensure both value and def have to the same type.
	if vType.Elem().Name() != defType.Name() {
		return errors.New("value and def have different struct types")
	}

	return setStruct(vVal, reflect.ValueOf(def))
}

func setStruct(vStr, defStr reflect.Value) (err error) {
	// Recursively compare the structs.
	for i := 0; i < vStr.NumField(); i++ {
		vf := vStr.Field(i)
		df := defStr.Field(i)

		// Skip unexported fields, etc.
		if !vf.CanSet() || !vf.CanInterface() {
			continue
		}

		// If the field is a struct again, set defaults on it as well.
		vfStr := toStructOrStructPointer(vf)
		if vfStr.IsValid() {
			err = setStruct(vfStr, toStructOrStructPointer(df))
			if err != nil {
				return
			}
			continue
		}

		// If the struct field is the zero value, set the default value.
		if vf.IsZero() {
			vf.Set(df)
		}
	}

	return
}

// StripDefaults sets every field of the struct pointer v to its zero value, if it matches
// the same field of the struct def. In case a struct field is itself a struct, the stripping
// continues recursively.
func StripDefaults(v, def interface{}) (err error) {
	// Ensure the type of v is a struct pointer.
	vType := reflect.TypeOf(v)
	if vType == nil || vType.Kind() != reflect.Ptr {
		return errors.New("value must be a struct pointer")
	}
	vVal := reflect.ValueOf(v).Elem()
	if vVal.Kind() != reflect.Struct {
		return errors.New("value must be a struct pointer")
	}

	// Ensure the type of def is a struct.
	defType := reflect.TypeOf(def)
	if defType == nil || defType.Kind() != reflect.Struct {
		return errors.New("def must be a struct value")
	}

	// Ensure both value and def have to the same type.
	if vType.Elem().Name() != defType.Name() {
		return errors.New("value and def have different struct types")
	}

	return stripStruct(vVal, reflect.ValueOf(def))
}

// stripStruct sets every field of the struct vStr to the zero value, if its value matches
// the same field of defStr.
func stripStruct(vStr, defStr reflect.Value) (err error) {
	// Recursively compare the structs.
	for i := 0; i < vStr.NumField(); i++ {
		vf := vStr.Field(i)
		df := defStr.Field(i)

		// Skip unexported fields, etc.
		if !vf.CanSet() || !vf.CanInterface() {
			continue
		}

		// If the field is a struct again, strip it as well.
		vfStr := toStructOrStructPointer(vf)
		if vfStr.IsValid() {
			err = stripStruct(vfStr, toStructOrStructPointer(df))
			if err != nil {
				return
			}
			continue
		}

		// Compare the two fields and zero out the field of the value struct,
		// if it is equal to the default field.
		if reflect.DeepEqual(vf.Interface(), df.Interface()) {
			// Same as default field, set to zero value.
			vf.Set(reflect.Zero(vf.Type()))
		}
	}

	return
}

// toStructOrStructPointer checks, if the given value is either a struct or a pointer
// to a struct and tries to convert it to a struct value.
// Whether the conversion was successful or not can be checked by calling IsValid() on the
// returned reflection value. When true, the conversion was successful.
func toStructOrStructPointer(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return reflect.ValueOf(nil)
	}
	if v.Kind() == reflect.Struct {
		return v
	}
	if v.Kind() == reflect.Ptr && v.Elem().Kind() == reflect.Struct {
		return v.Elem()
	}
	return reflect.ValueOf(nil)
}
