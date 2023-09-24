package main

import (
	"fmt"
	"github.com/amarnathcjd/gogram/telegram"
	"github.com/joho/godotenv"
	"github.com/k0kubun/pp"
	"mangaParser/config"
	"os"
	"regexp"
	"strings"
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

	var lastId int = 9290 // TODO read from file
	ReadLastId()
	for {
		fmt.Println(lastId)
		posts, postErr := client.GetMessages(conf.Channel, &telegram.SearchOption{IDs: lastId})
		fmt.Println("-----")
		if postErr != nil {
			pp.Println(postErr) // TODO log in file
			continue
		}
		post := posts[0].Message
		if post.Message == "" {
			fmt.Println("FINISH")
			WriteNewLastId()
			os.Exit(0)
		} else {
			lastId++
		}
		//pp.Println(post.Message)
		if Contains(post.Message, titles) {
			r, _ := regexp.Compile("https://teletype.in/@lrmanga/[^\x1b]+")
			t := r.FindString(pp.Sprint(post.ReplyMarkup))
			pp.Println(t)
		}
	}

	//p, _ := client.GetSendablePeer("@prbezposhady")
	//pp.Println(client.MessagesGetHistory(&telegram.MessagesGetHistoryParams{Peer: p, Limit: 10, OffsetID: -1}))
	//.GetSendablePeer(PeerID)
}

func ReadLastId()     {}
func WriteNewLastId() {}
func Contains(v string, arr []string) bool {
	for _, el := range arr {
		if strings.Contains(v, el) {
			return true
		}
	}
	return false
}
