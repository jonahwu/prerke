package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"

	"golang.org/x/crypto/ssh"

	//"google.golang.org/grpc/codes"

	"github.com/google/goexpect"
	"github.com/google/goterm/term"
)

const (
	//timeout = 10 * time.Minute
	timeout = 10 * time.Second
)

var (
	addr  = flag.String("address", "172.16.155.137:22", "address of telnet server")
	user  = flag.String("user", "root", "username to use")
	pass1 = flag.String("pass1", "promise", "password to use")
	pass2 = flag.String("pass2", "proimse", "alternate password to use")
)

func main() {
	flag.Parse()
	fmt.Println(term.Bluef("SSH Example"))

	sshClt, err := ssh.Dial("tcp", *addr, &ssh.ClientConfig{
		User:            *user,
		Auth:            []ssh.AuthMethod{ssh.Password(*pass1)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalf("ssh.Dial(%q) failed: %v", *addr, err)
	}
	defer sshClt.Close()

	e, _, err := expect.SpawnSSH(sshClt, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer e.Close()
	fmt.Println("into expect")

	e.ExpectBatch([]expect.Batcher{
		&expect.BCas{[]expect.Caser{
			//&expect.Case{R: regexp.MustCompile(`router#`), T: expect.OK()},
			//&expect.Case{R: regexp.MustCompile(`Login: `), S: *user,
			//		T: expect.Continue(expect.NewStatus(codes.PermissionDenied, "wrong username")), Rt: 3},
			&expect.Case{R: regexp.MustCompile("password:"), S: *pass1, T: expect.Next(), Rt: 1},
			//&expect.Case{R: regexp.MustCompile(`Password: `), S: *pass2,
			//	T: expect.Continue(expect.NewStatus(codes.PermissionDenied, "wrong password")), Rt: 1},
		}},
	}, timeout)

	fmt.Println(term.Greenf("All done"))
}
