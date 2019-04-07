//code from www.361way.com
package main

import (
	"fmt"
	gexpect "github.com/ThomasRooney/gexpect"
)

func main() {
	child, err := gexpect.Spawn("ssh-copy-id 172.16.155.137")
	if err != nil {
		panic(err)
	}
	fmt.Println("into expect")
	child.Expect("password")
	fmt.Println("end expect")
	child.SendLine("promise")
	child.Interact()
	child.Close()
}
