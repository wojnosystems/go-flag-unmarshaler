# Overview

Reads os.Args and turns them into Group objects which can be passed into the Unmarshall method. The Unmarshall method converts these flags into Go primitives or even custom structures.

# Usage

```go
package main

import (
  "github.com/wojnosystems/go-flag-unmarshaller"
  "log"
  "os"
)

type globalFlags struct {
  Enabled bool `flag:"enabled"`
}

type command1Flags struct {
  Host string `flag:"host" flag-short:"h"`
}

func main() {
  commands := flag_unmarshaller.Split(os.Args[1:])
  if len(commands) < 2 {
    log.Panic("at least 1 command required")
  }
  flagParser := flag_unmarshaller.New(&commands[0])
  var globals globalFlags
  if commands[0].CommandName == "" {
    _ = flagParser.Unmarshall(&globals)
  }
  flagParser = flag_unmarshaller.New(&commands[1])
  switch commands[1].CommandName {
    case "connect":
      var command1Options command1Flags 
      _ = flagParser.Unmarshall(&command1Options)
      if globals.Enabled {
        // do thing if enabled
      }
      Connect(command1Options.Host)
    default:
      log.Fatal("unrecognized command")
  }
}
```

# Flag values

Flags can have long names and optionally short names. Flags must have a long name, which, if not set, will default to the field name in the Go struct, case-sensitive.

Flags look like this:

```
--fieldName=bob
```

Values are always separated from the key name using an equal sign.

If you have a struct like this:

```go
type myStruct struct {
  FieldName string `flag:"fieldName" flag-usage:"Set this to some string to demonstrate"`
}
```

After calling Unmarshall, myStruct.FieldName will be set to "bob".

# FAQ

## POSIX says I can separate option values with optional spaces!

Yea. I don't care. This makes things needlessly complicated. Use equal signs. LMK if there's a special case I should consider for this.
