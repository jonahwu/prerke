package main

import (
	"bytes"
	"fmt"
	//"glog"
	gexpect "github.com/ThomasRooney/gexpect"
	//	"github.com/google/goexpect"
	//	"github.com/google/goterm/term"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	//"os"
	"regexp"
	"strings"
	"time"
)

var sshuser string
var sshpassword string
var sshport string
var sshAddress string
var deployUser string
var piRE = regexp.MustCompile(`3.14[0-9]*`)
var promptRE = regexp.MustCompile(`password`)

const (
	//command = `bc -l`
	timeout = 10 * time.Minute
)

type Config struct {
	Address string
	Info    []string
}

type Configs struct {
	Cfgs []Config `nodes`
}

//sshuser="root"
//sshpassword="promise"

func RemoteSSHRun(addr string, port string, cmd string) string {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	addrPort := fmt.Sprintf("%s:%s", addr, port)
	//client, err := ssh.Dial("tcp", "172.16.155.137:22", &ssh.ClientConfig{
	fmt.Println(addrPort)
	client, err := ssh.Dial("tcp", addrPort, &ssh.ClientConfig{
		User:            sshuser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshpassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		//Auth: []ssh.AuthMethod{ssh.Password("^Two^Ten=1024$")},
	})
	ce(err, "dial")

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	//ce(err, "new session")
	defer session.Close()
	var b bytes.Buffer
	session.Stdout = &b
	//if err := session.Run("/usr/bin/whoami"); err != nil {
	//cmd := "ls -al > scrremote"
	if err := session.Run(cmd); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())
	return b.String()

}

func sshCopy(deployUser string, sshpassword string) {
	//command := ""
	// cannot su user here
	//cmd := fmt.Sprintf("/usr/bin/ssh-copy-id %s", sshAddress)
	//cmd := fmt.Sprintf("sudo -iu %s ssh-copy-id %s@%s", deployUser, deployUser, sshAddress)
	cmd := fmt.Sprintf("sudo -iu %s ssh-copy-id -o StrictHostKeyChecking=no %s@%s", deployUser, deployUser, sshAddress)
	child, err := gexpect.Spawn(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println("into expect")
	child.Expect("password")
	fmt.Println("end expect")
	child.SendLine(sshpassword)
	child.Interact()
	child.Close()
}
func remoteTaskPipes(addr string, port string, cmds string) {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	addrPort := fmt.Sprintf("%s:%s", addr, port)
	fmt.Println(addrPort)
	client, err := ssh.Dial("tcp", addrPort, &ssh.ClientConfig{
		User:            sshuser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshpassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	ce(err, "dial")

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	stdinBuf, _ := session.StdinPipe()

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal(err)
	}
	err = session.Shell()
	cmds = fmt.Sprintf("%s;%s", cmds, "exit")
	//cmds := "ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \".ssh/id_rsa\" -q; exit"
	cmdlist := strings.Split(cmds, ";")
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
		fmt.Println(c)

	}

	session.Wait()
	fmt.Println((outbt.String() + errbt.String()))
	return
	//if err := session.Run("/usr/bin/whoami"); err != nil {
	//cmd := "ls -al > scrremote"
}

func createSshKey(addr string, port string) {
	ce := func(err error, msg string) {
		if err != nil {
			log.Fatalf("%s error: %v", msg, err)
		}
	}
	addrPort := fmt.Sprintf("%s:%s", addr, port)
	fmt.Println(addrPort)
	client, err := ssh.Dial("tcp", addrPort, &ssh.ClientConfig{
		User:            sshuser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshpassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	ce(err, "dial")

	session, err := client.NewSession()
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	stdinBuf, _ := session.StdinPipe()

	var outbt, errbt bytes.Buffer
	session.Stdout = &outbt

	session.Stderr = &errbt
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		log.Fatal(err)
	}
	err = session.Shell()
	cmds := "ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \".ssh/id_rsa\" -q; exit"
	cmdlist := strings.Split(cmds, ";")
	for _, c := range cmdlist {
		c = c + "\n"
		stdinBuf.Write([]byte(c))
		fmt.Println(c)

	}

	session.Wait()
	fmt.Println((outbt.String() + errbt.String()))
	return
	//if err := session.Run("/usr/bin/whoami"); err != nil {
	//cmd := "ls -al > scrremote"
}

func getAddress() Configs {

	var config Configs

	//filename := os.Args[1]
	filename := "precluster.yml"
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	//  source := []byte(data)

	if err := yaml.Unmarshal(source, &config); err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- config:\n%v\n\n", config)
	//fmt.Println("len of cfg", len(config.Cfgs))
	//fmt.Println("len of value", len(config.Cfgs[0].Info))
	//fmt.Println("first info value", config.Cfgs[0].Info[0])
	//fmt.Println(config.Cfgs[0].Address)
	return config

}

func main() {
	sshuser = "root"
	sshpassword = "promise"
	sshport = "22"
	sshAddress = "172.16.155.170"
	deployUser = "pentium"
	var cmds string
	config := getAddress()
	//fmt.Println("show first address:", config.Cfgs[0].Address)
	//fmt.Println("lens of address:", len(config.Cfgs))
	fmt.Println("----- check for 5 secs--------")
	time.Sleep(5 * time.Second)
	for ia := 0; ia < len(config.Cfgs); ia++ {
		sshAddress = config.Cfgs[ia].Address
		fmt.Println("installed of address:", sshAddress)

		cmds = "systemctl stop firewalld;systemctl disable firewalld"
		remoteTaskPipes(sshAddress, sshport, cmds)

		//os.Exit(3)

		/*
			cmd := "curl https://releases.rancher.com/install-docker/17.03.sh | sh"
			ret := RemoteSSHRun(sshAddress, sshport, cmd)
			fmt.Println(ret)
			userTask(sshAddress, sshport)
		*/
		/* install docker */
		cmds = "curl https://releases.rancher.com/install-docker/17.03.sh | sh"
		remoteTaskPipes(sshAddress, sshport, cmds)

		/* create user */
		// for ubuntu
		//cmds = fmt.Sprintf("adduser pentium --gecos \"First Last,RoomNumber,WorkPhone,HomePhone\" --disabled-password;echo \"pentium:%s\" | sudo chpasswd;gpasswd -a pentium docker", sshpassword)
		//cmds = fmt.Sprintf("adduser pentium ;echo \"pentium:%s\" | sudo chpasswd;usermod -aG wheel pentium;gpasswd -a pentium docker", sshpassword)
		cmds = fmt.Sprintf("adduser %s ;echo \"%s:%s\" | sudo chpasswd;usermod -aG wheel %s;gpasswd -a %s docker", deployUser, deployUser, sshpassword, deployUser, deployUser)
		remoteTaskPipes(sshAddress, sshport, cmds)

		/* generate ssh key */
		//for ubuntu
		//cmds = "ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \".ssh/id_rsa\" -q"
		//for centos
		cmds = "ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \"/root/.ssh/id_rsa\" -q"
		remoteTaskPipes(sshAddress, sshport, cmds)
		//remoteExpect()

		/* ssh-keygen in user pentium */
		//for ubuntu
		//cmds = fmt.Sprintf("sudo -iu %s ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \".ssh/id_rsa\" -q", deployUser)
		//for centos
		fmt.Println("create remote pentum sshkey")
		cmds = fmt.Sprintf("sudo -iu %s ssh-keygen -t rsa -C \"comment\" -P \"examplePassphrase\" -f \"/home/%s/.ssh/id_rsa\" -q", deployUser, deployUser)
		remoteTaskPipes(sshAddress, sshport, cmds)

		/* ssh-copy-id something */
		fmt.Println("into sshcopy")
		sshCopy(deployUser, sshpassword)
	}

	fmt.Println("trying launch the following command for testing")
	fmt.Println("sudo -u pentium ssh 'pentium@172.16.155.170' docker ps")
	fmt.Println("sudo -u pentium ssh 'pentium@172.16.155.170' systemctl status firewalld")

}
