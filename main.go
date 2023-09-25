package main

import (
	"fmt"
	"github.com/amarnathcjd/gogram/telegram"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"io"
	"mangaParser/config"
	"os"
	"regexp"
	"strings"
	"time"
)

var titles = []string{"Жрец порчи"}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	client, _ := telegram.NewClient(telegram.ClientConfig{
		AppID:    conf.Bot.AppId,
		AppHash:  conf.Bot.AppHash,
		LogLevel: telegram.LogInfo,
		Session:  "./client.session",
	})
	if err := client.Connect(); err != nil {
		panic(err)
	}
	if err := client.LoginBot(conf.Bot.Token); err != nil {
		panic(err)
	}

	lastId, err := ReadLastId("post_id.txt")
	if err != nil {
		panic(err)
	}
	for {
		ch := make(chan string)
		go func() {
			defer close(ch)
			posts, postErr := client.GetMessages(conf.Channel, &telegram.SearchOption{IDs: lastId})
			msg := ""

			if postErr != nil {
				pp.Println(postErr) // TODO log in file
				ch <- "error"
				return
			}
			post := posts[0].Message

			if post.Message != "" && Contains(post.Message, titles) {
				r, _ := regexp.Compile("https://teletype.in/@lrmanga/[^\x1b]+")
				t := r.FindString(pp.Sprint(post.ReplyMarkup))
				msg = t
			}

			ch <- msg
		}()

		select {
		case msg := <-ch:
			fmt.Println("Received", lastId, " - ", msg)
			if msg == "" {
				fmt.Println("Finishing at", lastId)
				if err := WriteNewLastId("post_id.txt", lastId); err != nil {
					panic(err)
				}
			}
			//} else if msg != "error" {
			//	fmt.Println(msg)
			//}

		case <-time.After(5 * time.Second):
			fmt.Println("Timed out", lastId)
			continue
		}
		lastId++

	}

}

func ReadLastId(filename string) (int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	var number int
	if _, err := fmt.Fscanf(file, "%d\n", &number); err != nil && err != io.EOF {
		return 0, err
	}
	return number, nil
}
func WriteNewLastId(filename string, id int) error {
	file, err := os.OpenFile(filename, os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("%d", id))
	return err
}
func Contains(v string, arr []string) bool {
	for _, el := range arr {
		if strings.Contains(v, el) {
			return true
		}
	}
	return false
}
