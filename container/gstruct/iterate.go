package gstruct

import (
	"github.com/cryptowilliam/goutil/basic/gerrors"
	"reflect"
)

type (
	IterateFn = func(v reflect.Value) (newVal reflect.Value, modified bool, err error)
)

// Iterate over the basic type members of a structure with a specific flag.
// Reference:
// https://github.com/IQ-tech/go-crypto-layer/blob/master/datacrypto/aesecb.go#L42
// EncryptStruct crawls all annotated struct properties and encrypts them in place
func Iterate(toModify interface{}, tagKey, tagVal string, iterFn IterateFn) error {
	if toModify == nil || reflect.TypeOf(toModify) == nil {
		return nil
	}
	cipherVType := reflect.TypeOf(toModify)
	if cipherVType.Kind() != reflect.Ptr {
		return gerrors.New("must receive a pointer, but received " + cipherVType.Kind().String())
	}

	cipherVType = reflect.TypeOf(toModify).Elem()
	if cipherVType.Kind() != reflect.Struct {
		return gerrors.New("must receive a pointer to a struct, but received " + cipherVType.Kind().String())
	}

	instanceValue := reflect.ValueOf(toModify).Elem()

	for i := 0; i < cipherVType.NumField(); i++ {
		currentFieldTag := cipherVType.Field(i).Tag
		hasTagVal, hasTagKey := currentFieldTag.Lookup(tagKey)
		field := instanceValue.Field(i)

		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.Ptr:
				if !field.IsNil() && field.Elem().IsValid() {
					switch field.Elem().Kind() {
					case reflect.Struct:
						err := Iterate(field.Interface(), tagKey, tagVal, iterFn)
						if err != nil {
							return err
						}
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
						reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
						if hasTagKey && hasTagVal == tagVal {
							modifiedV, modified, err := iterFn(field)
							if err != nil {
								return err
							}
							if modified {
								field.Set(modifiedV)
							}
						}
					default:
						if hasTagKey && hasTagVal == tagVal {
							return gerrors.New("Field must be a basic types or a pointer to a them")
						}
					}
				}
			case reflect.Struct:
				err := Iterate(field.Addr().Interface(), tagKey, tagVal, iterFn)
				if err != nil {
					return err
				}
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
				if hasTagKey && hasTagVal == tagVal {
					modifiedV, modified, err := iterFn(field)
					if err != nil {
						return err
					}
					if modified {
						field.Set(modifiedV)
					}
				}
			default:
				if hasTagKey && hasTagVal == tagVal {
					return gerrors.New("Field must be a basic types or a pointer to a them")
				}
			}
		}
	}

	return nil
}
