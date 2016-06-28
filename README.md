# colorgful ![report](https://goreportcard.com/badge/github.com/reconquest/colorgful)

Colorizer extension for [lorg](https://github.com/kovetskiy/lorg).

Basically, turns lorg output into this:

![dark](https://raw.githubusercontent.com/reconquest/colorgful/master/dark.png)

Or this:

![light](https://raw.githubusercontent.com/reconquest/colorgful/master/light.png)

colorgful comes with two embedded themes: light and dark.

# Example

```go
package main

import (
	"github.com/kovetskiy/lorg"
	"github.com/reconquest/colorgful"
)

func main() {
	log := lorg.NewLog()

	log.SetFormat(colorgful.MustApplyDefaultTheme(
		`* ${time} ${level} %s`,
		colorgful.Light,
	))

	log.SetLevel(lorg.LevelTrace)

	log.Trace("tracing dead status")

	log.Debug("debuggin dead status")

	log.Info("soon you will be dead")

	log.Warning("you are not prepared to be dead")

	log.Error("you're dead!")

	log.Fatal("stopping")
}
```
