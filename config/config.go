package config

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/sha3"
	"golang.org/x/term"
)

const CONFIG_FILE_PATH = "/.config/personal-website/config.json"

type Config struct {
	Scheme string `json:"scheme"`
	Host   string `json:"host"`
	Port   int    `json:"port"`

	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`

	Users []User `json:"users"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ReadConfigFile() (*Config, error) {
	h, err := os.UserHomeDir()

	if err != nil {
		log.Panic("Problem getting home directory: " + err.Error())
	}

	b, err := ioutil.ReadFile(h + CONFIG_FILE_PATH)

	if errors.Is(err, os.ErrNotExist) {
		log.Panic("Config file does not exist at ~" + CONFIG_FILE_PATH + "\n" + err.Error())
	}
	if err != nil {
		log.Panic("Problem opening ~" + CONFIG_FILE_PATH + "\n" + err.Error())
	}

	if err != nil {
		log.Panic("Problem reading config file\n" + err.Error())
	}

	var c Config
	err = json.Unmarshal(b, &c)

	if err != nil {
		log.Panic("Problem unmarshaling data from config\n" + err.Error())
	}

	if c.Users == nil {
		fmt.Print("Enter a username for the API user: ")
		reader := bufio.NewReader(os.Stdin)
		username, err := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		if err != nil {
			return nil, err
		}

		fmt.Print("Enter a password for the API user: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return nil, err
		}
		fmt.Println()

		hash := sha3.New256()
		hash.Write(passwordBytes)
		password := hex.EncodeToString(hash.Sum(nil))

		c.Users = []User{
			{
				Username: username,
				Password: password,
			},
		}

		f, err := os.OpenFile(h + CONFIG_FILE_PATH, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}

		_, err = f.Write(b)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}
