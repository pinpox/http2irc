package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/hoisie/mustache"
	"github.com/thoj/go-ircevent"
	"log"
	"net/http"
	"os"
	"strconv"
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
	ircDebug     bool
	notice       bool
	templatePath string
)

var ircBot *irc.Connection

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	// Check the content of the token and if it is set (lenght > 0 )
	if r.Header.Get("Token") != botToken || len(r.Header.Get("Token")) == 0 {
		log.Printf("Invalid token provided %s", r.Header.Get("Token"))
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

	if notice {
		ircBot.Notice(channel, data)
	} else {
		ircBot.Privmsg(channel, data)
	}
}

func main() {

	nick = os.Getenv("IRC_NICK")
	botToken = os.Getenv("IRC_BOT_TOKEN")
	channel = os.Getenv("IRC_CHANNEL")
	server = os.Getenv("IRC_SERVER")
	saslUser = os.Getenv("IRC_SASL_USER")
	saslPass = os.Getenv("IRC_SASL_PASS")
	hookListen = os.Getenv("IRC_LISTEN")
	templatePath = os.Getenv("IRC_TEMPLATE")
	ircDebug, _ = strconv.ParseBool(os.Getenv("IRC_DEBUG"))
	notice, _ = strconv.ParseBool(os.Getenv("IRC_NOTICE"))

	// Setup the bot
	ircBot = irc.IRC(nick, "none") // Don't care about user here, we use SASL
	ircBot.UseTLS = true
	ircBot.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	ircBot.AddCallback("001", func(e *irc.Event) { ircBot.Join(channel) })
	ircBot.UseSASL = true
	ircBot.SASLLogin = saslUser
	ircBot.SASLPassword = saslPass

	if ircDebug {
		ircBot.VerboseCallbackHandler = true
		ircBot.Debug = true
	}

	err := ircBot.Connect(server)
	if err != nil {
		log.Printf("Err %s", err)
		return
	}
	go ircBot.Loop()

	// Setup the handler
	http.HandleFunc("/webhook", handleWebhook)
	log.Println("server started")
	log.Fatal(http.ListenAndServe(hookListen, nil))
}
