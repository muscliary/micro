package screen

import (
	"fmt"
	"os"

	"github.com/zyedidia/micro/cmd/micro/config"
	"github.com/zyedidia/micro/cmd/micro/terminfo"
	"github.com/zyedidia/tcell"
)

var Screen tcell.Screen

// Init creates and initializes the tcell screen
func Init() {
	// Should we enable true color?
	truecolor := os.Getenv("MICRO_TRUECOLOR") == "1"

	tcelldb := os.Getenv("TCELLDB")
	os.Setenv("TCELLDB", config.ConfigDir+"/.tcelldb")

	// In order to enable true color, we have to set the TERM to `xterm-truecolor` when
	// initializing tcell, but after that, we can set the TERM back to whatever it was
	oldTerm := os.Getenv("TERM")
	if truecolor {
		os.Setenv("TERM", "xterm-truecolor")
	}

	// Initilize tcell
	var err error
	Screen, err = tcell.NewScreen()
	if err != nil {
		if err == tcell.ErrTermNotFound {
			err = terminfo.WriteDB(config.ConfigDir + "/.tcelldb")
			if err != nil {
				fmt.Println(err)
				fmt.Println("Fatal: Micro could not create terminal database file", config.ConfigDir+"/.tcelldb")
				os.Exit(1)
			}
			Screen, err = tcell.NewScreen()
			if err != nil {
				fmt.Println(err)
				fmt.Println("Fatal: Micro could not initialize a screen.")
				os.Exit(1)
			}
		} else {
			fmt.Println(err)
			fmt.Println("Fatal: Micro could not initialize a screen.")
			os.Exit(1)
		}
	}
	if err = Screen.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Now we can put the TERM back to what it was before
	if truecolor {
		os.Setenv("TERM", oldTerm)
	}

	if config.GetGlobalOption("mouse").(bool) {
		Screen.EnableMouse()
	}

	os.Setenv("TCELLDB", tcelldb)
}
