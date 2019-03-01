package command

import (
	"fmt"
	"os"
	"runtime"
)

func InstallSystem() {
	installScript()
	Exit("installed done")
}

func installScript() {
	fmt.Println(os.Args[0])
	fmt.Println(os.Getwd())
	fmt.Println(runtime.GOOS)
	script := "asdfasdfsadfasdf"
	fmt.Printf("%x", script)
}
