package main

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

func verify(files []fs.DirEntry, rule string, path string) bool {
	var verified bool = false
	if len(rule) == 0 {
		rule = "default"
	}
	if len(files) > 0 { // Insert actual rule check here
		verified = true
	}

	if !verified {
		var req map[string]string = map[string]string{}
		req["bot_token"] = ":"

		req["chat_id"] = ""

		req["text"] = "No new files: " + path // Replace with actual message from rule

		res, err := http.Get("https://api.telegram.org/bot" + req["bot_token"] + "/sendMessage?chat_id=" + req["chat_id"] + "&text=" + req["text"])
		if err != nil {
			log.Println("Could not send Telegram notification!")
			return false
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
	}

	return verified
}
