package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func InitUI() {
	clearScreen()
	fmt.Println("----------------------------- Socket Chat -----------------------------")
	fmt.Println("  <YOU>\t\t\t\t\t\t\t\t<Peer>")
}

func DisplayPeerMessage(msg string) {
	fmt.Printf("\n\t\t\t\t\t\t\t\t%s\n", msg)
}

func GetUserInput(scanner *bufio.Scanner) string {
	fmt.Print("\n ")
	scanner.Scan()
	return scanner.Text()
}

func ShowMessage(msg string) {
	fmt.Println(msg)
}

func clearScreen() {
	switch runtime.GOOS {
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
