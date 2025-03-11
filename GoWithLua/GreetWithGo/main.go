package main

import (
	"fmt"
	"os/exec"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func main() {
	var name string
	fmt.Println("what is your name champion? ")
	fmt.Scanln(&name)
	fmt.Println("its nice to meet you", name)

	now := time.Now()
	hour := now.Hour()

	Lua := lua.NewState()
	defer Lua.Close()

	Lua.SetGlobal("name", lua.LString(name))
	Lua.SetGlobal("hour", lua.LNumber(hour))

	if err := Lua.DoFile("greet.lua"); err != nil {
		fmt.Println("Sorry, I can't find or run the Lua file: ", err)
		return
	}

	fn := Lua.GetGlobal("greet")
	if fn.Type() != lua.LTFunction {
		fmt.Println("Lua function 'greet' not found")
		return
	}
	Lua.Push(fn)
	if err := Lua.PCall(0, 1, nil); err != nil {
		fmt.Println("Error calling Lua function:", err)
		return
	}

	result := Lua.Get(-1)
	if result.Type() != lua.LTString {
		fmt.Println("Lua function returned nothing or wrong type.")
		return
	}

	fmt.Println(result.String())
	if err := exec.Command("say", "-v", "Zarvox", result.String()).Run(); err != nil {
		fmt.Println("sorry i can't say this at this moment 🤐")
	}
}
