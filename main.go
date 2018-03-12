package main

import (
	"fmt"
	"time"
	"log"
	"strconv"
	"strings"
	"regexp"
	"encoding/json"
)

var (
	Mode   int
	Tokens []string
	Valid  []string
	Config Configuration
)

func main() {
	if !exists("config.json") {
		log.Fatal("Error, no config.json file found in base directory.")
	}

	config, _ := read("config.json")
	err := json.Unmarshal([]byte(config), &Config)

	if err != nil {
		log.Fatal("Error, config.json was unable to be parsed.")
	}

	if !exists("tokens.txt") {
		log.Fatal("Error, no tokens.txt file found in base directory.")
	}

	fmt.Println("Discord Token Checker v1.0\n	1) Line by line\n	2) Regex matching (slower)")
	for {
		fmt.Print("\nSelect a mode for parsing: ")
		_, err := fmt.Scanln(&Mode)

		if err != nil {
			fmt.Println("You didn't chose a valid option.")
			time.Sleep(time.Second * 2)
		} else {
			break
		}
	}

	switch Mode {
	case 1:

		Tokens, err := readln("tokens.txt")
		if err != nil {
			log.Fatal("Error reading tokens.txt, " + err.Error())
		}

		Tokens = strip(Tokens)

		for i, t := range Tokens {
			valid, str := validate(t)
			if valid {
				fmt.Println("("+strconv.Itoa(i)+"/"+strconv.Itoa(len(Tokens))+")", t, "is valid.")
				str += fmt.Sprint("\n("+strconv.Itoa(i)+"/"+strconv.Itoa(len(Tokens))+")", t, "is valid.")
				Valid = append(Valid, str)
			}
		}
	case 2:

		tokens, err := read("tokens.txt")
		if err != nil {
			log.Fatal("Error reading tokens.txt, " + err.Error())
		}

		r, _ := regexp.Compile(`([N|M][a-zA-Z\d-_]{23}[.][a-zA-Z\d-_]{6}[.][a-zA-Z\d-_]{27})`)
		Tokens := strip(r.FindAllString(tokens, -1))

		for i, t := range Tokens {
			valid, str := validate(t)
			if valid {
				fmt.Println("("+strconv.Itoa(i)+"/"+strconv.Itoa(len(Tokens))+")", t, "is valid.")
				str += fmt.Sprint("\n("+strconv.Itoa(i)+"/"+strconv.Itoa(len(Tokens))+")", t, "is valid.")
				Valid = append(Valid, str)
			}
		}
	default:
		log.Fatal("You chose a value outside of the option index.")
	}

	if err := write("valid.txt", strings.Join(Valid, "\n")); err == nil {
		fmt.Println("Wrote", len(Valid), "tokens to valid.txt")
	} else {
		fmt.Println("Failed to write", len(Valid), "tokens to valid.txt")
		fmt.Println("Printing them to chat to not lose data:")

		time.Sleep(time.Second * 2)

		fmt.Println(Valid)
	}

}

type Configuration struct {
	MinimumMembers int  `json:"minimum_members"`
	IncludeBots    bool `json:"include_bots"`
	MinimumGuilds  int  `json:"minimum_guilds"`
	Nukeable       bool `json:"nukeable"`
}
