package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
)

func verify(files []fs.DirEntry, rule string, path string) bool {
	var verified bool = false
	if len(rule) == 0 {
		rule = "default"
	}

	verified = checkrules(rule, files, path)

	if !verified {
		notify("Rule: _" + rule + "_ check at " + path + " failed")
	}

	return verified
}

func checkrules(rule string, files []fs.DirEntry, path string) bool {
	switch rule {
	case "default":
	case "exists":
		if len(files) > 0 { // Insert actual rule check here
			return true
		}
	case "unlocked":
		for _, f := range files {
			finfo, _ := f.Info()
			t, err := os.OpenFile(path+finfo.Name(), os.O_RDWR, 0755)
			if err != nil {
				return false
			}
			defer t.Close()
		}
		return true
	}
	return false
}

func notify(text string) error {
	res, err := http.Get("https://api.telegram.org/bot" + config2["bot_token"] + "/sendMessage?chat_id=" + config2["chat_id"] + "&text=" + text)
	if err != nil {
		log.Println("Could not send Telegram notification!")
		return errors.New("could not send Telegram notification")
	}
	defer res.Body.Close()

	var buf []byte = make([]byte, res.ContentLength)
	_, err = res.Body.Read(buf)
	if err != nil {
		log.Println("Could not get a response from Telegram, notification status unknown!")
	} else {
		response_string := fmt.Sprintln("Response:", string(buf))
		log.Println(response_string)
	}
	return nil
}
