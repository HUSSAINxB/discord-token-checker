package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"io/ioutil"
	"strings"
	"os"
	"path/filepath"
	"time"
)

func validate(token string) (bool, string) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return false, ""
	}

	if s.Open() != nil {
		return false, ""
	}

	time.Sleep(time.Second * 2)
	defer s.Close()

	if len(s.State.Guilds) == 0 {
		return false, ""
	}

	members := members(s)
	if members < Config.MinimumMembers {
		return false, ""
	}

	fmt.Println("\nUsername:", s.State.User.Username, "\nTotal Members:", members, "\nID:", s.State.User.ID)
	return true, fmt.Sprint("\nUsername:", s.State.User.Username, "\nTotal Members:", members, "\nID:", s.State.User.ID)

}

func members(s *discordgo.Session) int {
	members := 0
	for _, e := range s.State.Guilds {
		g, err := s.Guild(e.ID)
		if err != nil {
			continue
		}

		for _, m := range g.Members {
			if !m.User.Bot {
				members += 1
			}
		}
	}

	return members
}

func strip(elements []string) []string {
	encountered := map[string]bool{}
	var result []string

	for v := range elements {
		if encountered[elements[v]] == false {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	return result
}

func readln(path string) ([]string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(bytes), "\n"), nil
}

func read(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func write(path, content string) error {
	os.MkdirAll(filepath.Dir(path), 0777)
	os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0777)

	return ioutil.WriteFile(path, []byte(content), 0777)
}

func exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
