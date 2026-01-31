package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func banner() {
	fmt.Println(Blue + `
███████╗████████╗███████╗ ██████╗  ██████╗ 
██╔════╝╚══██╔══╝██╔════╝██╔════╝ ██╔═══██╗
███████╗   ██║   █████╗  ██║  ███╗██║   ██║
╚════██║   ██║   ██╔══╝  ██║   ██║██║   ██║
███████║   ██║   ███████╗╚██████╔╝╚██████╔╝
╚══════╝   ╚═╝   ╚══════╝ ╚═════╝  ╚═════╝ ` + Reset)
}

func prompt(label string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(Yellow + label + Reset)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func pause() {
	fmt.Print("\nPress enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
