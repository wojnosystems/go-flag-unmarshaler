package go_flag_unmarshaller

import "strings"

const (
	keyValueSeparator = "="
	argListEnd        = "--"
)

func Split(args []string) (out []Group) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if argListEnd == arg {
			// indicates that no more values should be parsed
			return
		}
		if isFlag(arg) {
			// Create the initial group without a name, as we have a flag with no group-name
			if len(out) == 0 {
				out = append(out, Group{
					CommandName: "",
				})
			}
			var value string
			key := arg
			if flagHasValue(arg) {
				splits := strings.SplitN(arg, keyValueSeparator, 2)
				key = splits[0]
				value = splits[1]
			}
			if isShortFlag(key) {
				for _, shortFlag := range flagsFromShortFlag(key) {
					out[len(out)-1].Flags = append(out[len(out)-1].Flags, KeyValue{
						Key:   shortFlag,
						Value: value,
					})
				}
			} else {
				out[len(out)-1].Flags = append(out[len(out)-1].Flags, KeyValue{
					Key:   key,
					Value: value,
				})
			}
		} else {
			// positional argument, this is called a Group
			out = append(out, Group{
				CommandName: arg,
			})
		}
	}
	return
}

func isFlag(arg string) bool {
	return strings.HasPrefix(arg, "-") && len(arg) != 1
}

func flagHasValue(flagWithValue string) bool {
	return strings.Contains(flagWithValue, keyValueSeparator)
}

// isShortFlagDef returns true if the flag definition is a short-named flag (single rune)
func isShortFlagDef(flagName string) bool {
	if len(flagName) == 2 {
		return flagName[0] == '-' && flagName[1] != '-'
	}
	return false
}

func isShortFlag(flagKey string) bool {
	if len(flagKey) >= 2 {
		return flagKey[0] == '-' && flagKey[1] != '-'
	}
	return false
}

func flagsFromShortFlag(shortFlag string) (out []string) {
	runes := []rune(shortFlag)
	out = make([]string, len(runes)-1)
	for i, r := range runes[1:] {
		out[i] = "-" + string(r)
	}
	return
}
