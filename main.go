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

type PostStatus int

const (
	POST_OK PostStatus = iota
	POST_EMPTY
	POST_ERR
)

type Post struct {
	Status PostStatus
	Name   string
	Url    string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	client := CreateClient(*conf)

	lastId, err := ReadLastId("post_id.txt")
	if err != nil {
		panic(err)
	}
	for {
		ch := make(chan Post)
		go func() {
			defer close(ch)
			posts, postErr := client.GetMessages(conf.Channel, &telegram.SearchOption{IDs: lastId})
			post := Post{}

			if postErr != nil {
				pp.Println(postErr.Error())
				post.Status = POST_ERR
				ch <- post
				return
			}
			postData := posts[0].Message
			if postData.Message == "" {
				post.Status = POST_EMPTY
			} else {
				post.Status = POST_OK
				post.Name = postData.Message

				r, _ := regexp.Compile("https://teletype.in/@lrmanga/[^\x1b]+")
				t := r.FindString(pp.Sprint(postData.ReplyMarkup))
				post.Url = t
			}
			ch <- post
		}()

		select {
		case post := <-ch:
			fmt.Println("Received", lastId, " - ", post.Name)
			if post.Status == POST_ERR && !client.IsConnected() {
				client.Disconnect()
				client = CreateClient(*conf)
				continue
			} else if post.Status == POST_EMPTY {
				fmt.Println("Finishing at", lastId)
				if err := WriteNewLastId("post_id.txt", lastId); err != nil {
					panic(err)
				}
				os.Exit(0)
			} else if Contains(post.Name, titles) {

				//client.SendMessage(conf.UserId, post.Url)
			}
		case <-time.After(4 * time.Second):
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
func CreateClient(conf config.Config) *telegram.Client {
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
	return client
}
