package flag_unmarshaler

import (
	optional_parse_registry "github.com/wojnosystems/go-optional-parse-registry"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// unmarshaler creates an environment parser given the provided registry
type flagsInternal struct {
	flagsConfig
	root interface{}
}

func newFlagsInternal(config flagsConfig, root interface{}) *flagsInternal {
	return &flagsInternal{
		flagsConfig: config,
		root:        root,
	}
}

var (
	defaultNoOpSetReceiver = &SetReceiverNoOp{}
	defaultParseRegister   = optional_parse_registry.NewWithGoPrimitives()
)

func (f *flagsInternal) SetValue(structFullPath string, fieldV reflect.Value, fieldStruct reflect.StructField) (handled bool, err error) {
	flagPaths := f.getFlagPaths(structFullPath, fieldStruct)
	if !f.parseRegistry.IsSupported(fieldV.Addr().Interface()) {
		return
	}
	// handled is true because we support deserializing this structure type.
	handled = true
	for _, name := range flagPaths {
		value, ok := f.flags.Get(name)
		if ok {
			// Some environment value was set, use it
			valueWasSet := false
			valueWasSet, err = f.parseRegistry.SetValue(fieldV.Addr().Interface(), value)
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
			if valueWasSet {
				f.emitter.ReceiveSet(structFullPath, structFullPath, value)
			}
			return
		}
	}
	return
}

var flagIndexRegexp = regexp.MustCompile(`^(\d+)`)

func (f *flagsInternal) SliceLen(structFullPath string, fieldV reflect.Value, fieldStruct reflect.StructField) (length int, err error) {
	flagPaths := f.getFlagPaths(structFullPath, fieldStruct)
	maxIndex := int64(-1)
	for _, path := range flagPaths {
		for _, key := range f.flags.Keys(path + "[") {
			possibleNumber := flagIndexRegexp.FindString(key[len(path+"["):])
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
	}
	length = int(maxIndex + 1)
	return
}

// getFlagPaths calculates all possible flag names for a path in the structure, converting field names (the default names) into flag and short-flag variants.
// You can mix-and-match short, long and default names. All of them will be tried by the SetValue method.
func (f *flagsInternal) getFlagPaths(structFullPath string, fieldStruct reflect.StructField) (flagPaths []string) {
	structPathComponents := strings.Split(structFullPath, ".")
	currentT := reflect.ValueOf(f.root).Elem().Type()
	for componentIndex, componentName := range structPathComponents {
		newPaths := make([]string, 0, 2)
		fieldName := ""
		index := ""
		if isComponentSlice(componentName) {
			parts := strings.Split(componentName, "[")
			fieldName = parts[0]
			index = "[" + parts[1]
		} else {
			fieldName = componentName
		}
		component, _ := currentT.FieldByName(fieldName)
		longFlagName := component.Tag.Get("flag")
		shortFlagName := component.Tag.Get("flag-short")
		if longFlagName == "" && shortFlagName == "" {
			longFlagName = fieldName
		}
		if componentIndex == 0 {
			newPaths = append(newPaths, "--"+longFlagName+index)
		} else {
			for _, path := range flagPaths {
				newPaths = append(newPaths, path+index+"."+longFlagName)
			}
		}
		if shortFlagName != "" {
			if componentIndex == 0 {
				newPaths = append(newPaths, "-"+shortFlagName+index)
			} else {
				for _, path := range flagPaths {
					newPaths = append(newPaths, path+index+"."+shortFlagName)
				}
			}
		}
		// next cycle
		flagPaths = newPaths
		currentT = component.Type
		if index != "" {
			currentT = currentT.Elem()
		}
	}
	return
}

func isComponentSlice(componentName string) bool {
	return strings.HasSuffix(componentName, "]")
}
