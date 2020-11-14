package flag_unmarshaler

import (
	"fmt"
	into_struct "github.com/wojnosystems/go-into-struct"
	optional_parse_registry "github.com/wojnosystems/go-optional-parse-registry"
	"regexp"
	"strconv"
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

func (f *flagsInternal) SetValue(structFullPath into_struct.Path) (handled bool, err error) {
	field := structFullPath.Top()
	if field == nil {
		return
	}
	flagPaths := f.getFlagPaths(structFullPath)
	if !f.parseRegistry.IsSupported(field.Value().Addr().Interface()) {
		return
	}
	// handled is true because we support deserializing this structure type.
	handled = true
	for _, flagPath := range flagPaths {
		value, ok := f.flags.Get(flagPath)
		if ok {
			// Some flag value was set, use it
			valueWasSet := false
			valueWasSet, err = f.parseRegistry.SetValue(field.Value().Addr().Interface(), value)
			if err != nil {
				err = newParseError(structFullPath.String(), flagPath, err)
				return
			}
			if valueWasSet {
				f.emitter.ReceiveSet(structFullPath, flagPath, value)
			}
			return
		}
	}
	return
}

var flagIndexRegexp = regexp.MustCompile(`^(\d+)`)

func (f *flagsInternal) SliceLen(structFullPath into_struct.Path) (length int, err error) {
	flagPaths := f.getFlagPaths(structFullPath)
	maxIndex := int64(-1)
	for _, path := range flagPaths {
		for _, key := range f.flags.Keys(path + "[") {
			possibleNumber := flagIndexRegexp.FindString(key[len(path+"["):])
			if "" != possibleNumber {
				var index int64
				index, err = strconv.ParseInt(possibleNumber, 10, 0)
				if err != nil {
					err = newParseError(structFullPath.String(), key, err)
					return
				}
				if index > maxIndex {
					maxIndex = index
				}
			} else {
				err = newParseError(structFullPath.String(), key, fmt.Errorf("index was not a number"))
			}
		}
	}
	length = int(maxIndex + 1)
	return
}

// getFlagPaths calculates all possible flag names for a path in the structure, converting field names (the default names) into flag and short-flag variants.
// You can mix-and-match short, long and default names. All of them will be tried by the SetValue method.
func (f *flagsInternal) getFlagPaths(structFullPath into_struct.Path) (flagPaths []string) {
	for pathIndex, pathPart := range structFullPath.Parts() {
		newPaths := make([]string, 0, 2)
		index := ""
		if slicePart, ok := pathPart.(into_struct.PathSliceParter); ok {
			index = fmt.Sprintf("[%d]", slicePart.Index())
		}
		longFlagName := pathPart.StructField().Tag.Get("flag")
		shortFlagName := pathPart.StructField().Tag.Get("flag-short")
		if longFlagName == "" && shortFlagName == "" {
			longFlagName = pathPart.Name()
		}
		if pathIndex == 0 {
			newPaths = append(newPaths, "--"+longFlagName+index)
		} else {
			for _, path := range flagPaths {
				newPaths = append(newPaths, path+index+"."+longFlagName)
			}
		}
		if shortFlagName != "" {
			if pathIndex == 0 {
				newPaths = append(newPaths, "-"+shortFlagName+index)
			} else {
				for _, path := range flagPaths {
					newPaths = append(newPaths, path+index+"."+shortFlagName)
				}
			}
		}
		// next cycle
		flagPaths = newPaths
	}
	return
}
