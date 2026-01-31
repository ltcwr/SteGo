package main

import (
	"fmt"
	"os"
)

func clear() {
	fmt.Print("\033[H\033[2J")
}

func mainMenu() string {
	fmt.Println("1) Encode")
	fmt.Println("2) Decode")
	fmt.Println("3) Exit")
	return prompt("\nChoice: ")
}

func main() {
	for {
		clear()
		banner()

		choice := mainMenu()

		switch choice {
		case "1":
			clear()
			banner()
			fmt.Println(Green + "Encode Mode\n" + Reset)

			in := prompt("Input Image: ")
			out := prompt("Output Image: ")
			msg := prompt("Message : ")
			key := prompt("Secret Key: ")

			err := encode(in, out, msg, key)
			if err != nil {
				fmt.Println(Red+"Error:", err, Reset)
			} else {
				fmt.Println(Green + "\n Input encoded with success" + Reset)
			}
			pause()

		case "2":
			clear()
			banner()
			fmt.Println(Green + "Decode Mode\n" + Reset)

			in := prompt("Input image: ")
			key := prompt("Secret Key: ")

			msg, err := decode(in, key)
			if err != nil {
				fmt.Println(Red+"Error:", err, Reset)
			} else {
				fmt.Println(Green+"\nFound message :\n"+Reset, msg)
			}
			pause()

		case "3":
			fmt.Println(Green + "Bye !" + Reset)
			os.Exit(0)

		default:
			fmt.Println(Red + "Invalid Choice" + Reset)
			pause()
		}
	}
}
