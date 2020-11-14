package flag_unmarshaler

import (
	"os"
	"strings"
)

const (
	keyValueSeparator = "="
	argListEnd        = "--"
)

func SplitArgs() (out []Group) {
	return Split(os.Args[1:])
}

// Split converts os.Args[1:] into a slice of Groups. Each group may contain flags as key-value pairs.
// In clis, each command/sub-command divides a domain of flags which may be interpreted by the parser
// For example:
// `mycli --global command1 --for-command1 subcommand --for-subcommand --another-for-sub-command required positional args`
// Split will take the above command and break it into sub-commands. The first group is always the global and will always exist.
// args should not contain the first item, the executable's path. e.g. we expect os.Args[1:]
func Split(args []string) (out []Group) {
	// Always create the first group. Even it there's no arguments, this is considered the global group
	out = append(out, Group{
		CommandName: "",
	})
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if argListEnd == arg {
			// indicates that no more values should be parsed
			return
		}
		if isFlag(arg) {
			// Create the initial group without a name, as we have a flag with no group-name
			key := arg
			value := "true"
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
