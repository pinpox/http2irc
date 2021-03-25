package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/thoj/go-ircevent"
	"log"
	"net/http"
	"os"
)

// Configuration variables
var (
	botToken     string
	nick         string
	channel      string
	server       string
	saslUser     string
	saslPass     string
	hookListen   string
	ircDebug     string
	notice       string
	templatePath string
)

var ircBot *irc.Connection

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Token") != botToken {
		http.Error(w, "Not authorized", 401)
		return
	}

	var dataPost interface{}

	err := json.NewDecoder(r.Body).Decode(&dataPost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := mustache.RenderFile(templatePath, dataPost)
	// fmt.Println(data)

	// ircBot.Privmsg(channel, "msg") // sends a message to either a certain nick or a channel
	ircBot.Notice(channel, data) //send notices

	// _, err := io.Copy(os.Stdout, r.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
}

func main() {

	nick = os.Getenv("IRC_NICK")
	botToken = os.Getenv("IRC_BOT_TOKEN")
	channel = os.Getenv("IRC_CHANNEL")
	server = os.Getenv("IRC_SERVER")
	saslUser = os.Getenv("IRC_SASL_USER")
	saslPass = os.Getenv("IRC_SASL_PASS")
	hookListen = os.Getenv("IRC_LISTEN")
	ircDebug = os.Getenv("IRC_DEBUG")
	notice = os.Getenv("IRC_NOTICE")
	templatePath = os.Getenv("IRC_TEMPLATE")

	ircBot = irc.IRC(nick, "test") //nick, user
	ircBot.UseTLS = true
	ircBot.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// For debugging info
	// ircBot.VerboseCallbackHandler = true
	// ircBot.Debug = true

	ircBot.AddCallback("001", func(e *irc.Event) {
		ircBot.Join(channel)
		// ircBot.Privmsg(channel, "Hi, I'm a bot.")
	})
	ircBot.UseSASL = true
	ircBot.SASLLogin = saslUser
	ircBot.SASLPassword = saslPass

	err := ircBot.Connect(server)
	if err != nil {
		fmt.Printf("Err %s", err)
		return
	}
	go ircBot.Loop()

	http.HandleFunc("/webhook", handleWebhook)

	log.Println("server started")
	log.Fatal(http.ListenAndServe(hookListen, nil))
}
