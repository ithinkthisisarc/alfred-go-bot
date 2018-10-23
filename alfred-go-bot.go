package main

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    gocfg "gocfg-dev"
    "log"
)

type bot struct {
    Config *AppConf
    ConfigManager *gocfg.ConfigManager
    Session *discordgo.Session
}

var err error
var Bot *bot

func init() {
    Bot = &bot{}
}

func (b *bot) Run() {

    LoadConfig(b)
    b.Session, err = discordgo.New(fmt.Sprintf("Bot %s", b.Config.Auth.Token))
    FatalCheck("Could't create Discord session", err)

    // Handlers
    b.Session.AddHandler(func (s *discordgo.Session, nm *discordgo.MessageCreate) {
        if nm.Content == "!help" {
            s.ChannelMessageSend(nm.ChannelID, "Some help for you, my friend! <3")
        }
    })

    // Client Ready event
    if err = b.Session.Open(); err != nil {
        FatalReport("Could't open Discord websocket connection", err)
    }
    defer b.Session.Close()

    fmt.Println("|> Handlers attached")

    if _, err = b.Session.UserUpdateStatus("invisible"); err != nil {
        log.Printf("Could't change status to \"invisible\": \n\t%s",
            err.Error())
    }

    fmt.Printf("|> Bot started and logged as:\n\t[ID: %s | Name: %s]\n",
        b.Session.State.User.ID, b.Session.State.User.Username)

    <-make(chan struct{})
}
