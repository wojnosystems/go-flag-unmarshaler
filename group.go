package flag_unmarshaller

import "strings"

type Group struct {
	CommandName string
	Flags       []KeyValue
}

// Get the value of a single environment with the name envNamed
func (v *Group) Get(flagNamed string) (value string, ok bool) {
	for _, flag := range v.Flags {
		if flag.Key == flagNamed {
			value, ok = flag.Value, true
		}
	}
	return
}

// Keys get a list of keys that begin with the prefix. If "" is passed, matches all and returns all keys
func (v *Group) Keys(prefix string) (out []string) {
	uniqueKeys := make(map[string]bool)
	for _, flag := range v.Flags {
		if strings.HasPrefix(flag.Key, prefix) {
			uniqueKeys[flag.Key] = true
		}
	}
	for key := range uniqueKeys {
		out = append(out, key)
	}
	return
}
