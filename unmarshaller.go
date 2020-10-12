package go_flag_unmarshaller

import (
	"fmt"
	"github.com/wojnosystems/go-optional-parse-registry"
	"github.com/wojnosystems/go-parse-register"
	"reflect"
	"regexp"
	"strconv"
)

// unmarshaller creates an environment parser given the provided registry
type unmarshaller struct {
	// flags is the source of environment variables.
	// If you leave it blank, it will default to using the operating system environment variables with no prefixes.
	flags flagReader
	// ParseRegistry maps go-default and custom types to members of the provided structure. If left blank, defaults to just Go's primitives being mapped
	ParseRegistry *parse_register.Registry
}

func New(flags flagReader) Unmarshaller {
	return NewWithTypeParsers(
		flags,
		defaultParseRegister)
}

func NewWithTypeParsers(flags flagReader, parseRegistry *parse_register.Registry) Unmarshaller {
	return &unmarshaller{
		flags:         flags,
		ParseRegistry: parseRegistry,
	}
}

var (
	defaultNoOpSetReceiver = &SetReceiverNoOp{}
)

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func (e *unmarshaller) Unmarshall(into interface{}) (err error) {
	return e.UnmarshallWithEmitter(into, defaultNoOpSetReceiver)
}

// Unmarshall reads the environment variables and writes them to into.
// into should be a reference to a struct
// This method will do some basic checks on the into value, but to help developers pass in the correct values
func (e *unmarshaller) UnmarshallWithEmitter(into interface{}, emitter SetReceiver) (err error) {
	rootV := reflect.ValueOf(into)
	err = e.validateDestination(rootV, rootV.Type())
	if err != nil {
		return
	}
	structV := rootV.Elem()
	structT := structV.Type()
	err = e.unmarshallStruct("", structV, structT, emitter)
	return
}

// validateDestination does some basic checks to help users of this class avoid common pitfalls with more helpful messages
func (e *unmarshaller) validateDestination(rootV reflect.Value, rootT reflect.Type) (err error) {
	if rootV.IsNil() {
		return NewErrProgramming("'into' argument must be not be nil")
	}
	if rootT.Kind() != reflect.Ptr {
		return NewErrProgramming("'into' argument must be a reference")
	}
	if rootV.Elem().Kind() != reflect.Struct {
		return NewErrProgramming("'into' argument must be a struct")
	}
	err = validateStruct(rootT.Elem())
	return
}

// unmarshallStruct is the internal method, which can be called recursively. This performs the heavy-lifting
func (e *unmarshaller) unmarshallStruct(structParentPath string, structRefV reflect.Value, structRefT reflect.Type, emitter SetReceiver) (err error) {
	for i := 0; i < structRefV.NumField(); i++ {
		fieldV := structRefV.Field(i)
		fieldT := structRefT.Field(i)
		err = e.unmarshallField(structParentPath, fieldV, fieldT, emitter)
		if err != nil {
			return
		}
	}
	return
}

// unmarshallField unmarshalls a value into a single field in a struct. Could be the root struct or a nested struct
func (e *unmarshaller) unmarshallField(structParentPath string, fieldV reflect.Value, fieldT reflect.StructField, emitter SetReceiver) (err error) {
	if fieldV.CanSet() {
		if fieldT.Type.Kind() == reflect.Slice {
			flagNames := flagNamesOrDefault(fieldT, structParentPath)
			for _, name := range flagNames {
				fieldPath := appendStructPath(structParentPath, name)
				err = e.unmarshallSlice(fieldPath, fieldV, emitter)
				if err != nil {
					return
				}
			}
		} else {
			err = e.unmarshallValue(structParentPath, fieldV, fieldT, emitter)
			if err != nil {
				return
			}
		}
	}
	return
}

// unmarshallValue extracts a single value and sets it to a value in a struct
func (e *unmarshaller) unmarshallValue(structFullPath string, fieldV reflect.Value, fieldT reflect.StructField, emitter SetReceiver) (err error) {
	flagNames := flagNamesOrDefault(fieldT, structFullPath)
	for _, name := range flagNames {
		value, ok := e.flags.Get(name)
		if ok {
			// Some environment value was set, use it
			var valueSet bool
			valueSet, err = e.parseRegistry().SetValue(fieldV.Addr().Interface(), value)
			if err != nil {
				err = &ParseError{
					Path: StructFlagPath{
						StructPath: structFullPath,
						FlagPath:   structFullPath,
					},
					originalErr: err,
				}
				return
			}
			if valueSet {
				emitter.ReceiveSet(structFullPath, structFullPath, value)
			}
			return
		}
	}
	// fall back: no environment value found or was not set due to lack of type support
	if fieldT.Type.Kind() == reflect.Struct {
		err = e.unmarshallStruct(structFullPath, fieldV, fieldT.Type, emitter)
	}
	return
}

var defaultParseRegister = optional_parse_registry.Register(parse_register.RegisterGoPrimitives(&parse_register.Registry{}))

// parseRegistry obtains a copy of the current registry, or uses the default go primitives, for convenience
func (e *unmarshaller) parseRegistry() *parse_register.Registry {
	if e.ParseRegistry == nil {
		e.ParseRegistry = defaultParseRegister
	}
	return e.ParseRegistry
}

// flagNamesOrDefault attempts to read the tags to obtain an alternate name, if no tag found, defaults back to
// using the name provided to the field when the member was defined in Go
func flagNamesOrDefault(fieldT reflect.StructField, structPrefix string) (fieldNames []string) {
	fieldName := fieldT.Tag.Get("flag")
	if "" != fieldName {
		fieldNames = append(fieldNames, "--"+appendStructPath(structPrefix, fieldName))
	}
	fieldName = fieldT.Tag.Get("flag-short")
	if "" != fieldName {
		fieldNames = append(fieldNames, "-"+appendStructPath(structPrefix, fieldName))
	}
	// default back to long name if and only if nothing was set
	if len(fieldNames) == 0 {
		fieldNames = append(fieldNames, "--"+appendStructPath(structPrefix, fieldT.Name))
	}
	return
}

// appendStructPath concatenates the parent path name with the current field's name
func appendStructPath(parent string, name string) string {
	if parent != "" {
		return parent + "." + name
	}
	return name
}

// unmarshallSlice operates on a slice of objects. It will initialize the slice, then populate all of its members
// from the environment variables
func (e *unmarshaller) unmarshallSlice(sliceFieldPath string, sliceValue reflect.Value, emitter SetReceiver) (err error) {
	var length int
	length, err = elementsInSliceWithAddressPrefix(e.flags, sliceFieldPath+"[")
	if err != nil {
		return
	}
	if length > 0 {
		newSlice := reflect.MakeSlice(sliceValue.Type(), length, length)
		sliceValue.Set(newSlice)
		for i := 0; i < length; i++ {
			sliceElement := newSlice.Index(i)
			err = e.unmarshallStruct(sliceFieldPath+"["+strconv.FormatInt(int64(i), 10)+"]", sliceElement, sliceElement.Type(), emitter)
			if err != nil {
				return
			}
		}
	}
	return
}

var flagIndexRegexp = regexp.MustCompile(`^(\d+)`)

// elementsInSliceWithAddressPrefix returns how big a slice should be to hold all of the variables defined in the environment
func elementsInSliceWithAddressPrefix(flags flagReader, pathPrefix string) (length int, err error) {
	maxIndex := int64(-1)
	for _, key := range flags.Keys(pathPrefix) {
		possibleNumber := flagIndexRegexp.FindString(key[len(pathPrefix):])
		if "" != possibleNumber {
			var index int64
			index, err = strconv.ParseInt(possibleNumber, 10, 0)
			if err != nil {
				return
			}
			if index > maxIndex {
				maxIndex = index
			}
		}
	}
	length = int(maxIndex + 1)
	return
}

func validateStruct(structT reflect.Type) (err error) {
	for i := 0; i < structT.NumField(); i++ {
		fieldT := structT.Field(i)
		err = validateShortFlag(fieldT, structT)
		if err != nil {
			return
		}
		if fieldT.Type.Kind() == reflect.Slice {
			err = validateStruct(fieldT.Type.Elem())
			if err != nil {
				return
			}
		}
	}
	return
}

func validateShortFlag(fieldT reflect.StructField, structT reflect.Type) error {
	shortFlagDef := fieldT.Tag.Get("flag-short")
	if len(shortFlagDef) != 0 {
		if len(shortFlagDef) != 1 {
			return fmt.Errorf(`flag-short "%s" on field named: "%s.%s.%s" must be exactly 1 rune`, shortFlagDef, fieldT.PkgPath, structT.Name(), fieldT.Name)
		}
	}
	return nil
}
