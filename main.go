package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorcon/rcon"
	tele "gopkg.in/telebot.v3"
)

var modInstalled = false

func checkMod(rconClient *rcon.Conn) bool {
	resp, _ := rconClient.Execute("/say")
	return !strings.Contains(resp, "Unknown")
}

func initBot(token string, group int64, rconClient *rcon.Conn) (*tele.Bot, error) {
	var commandFormat string

	if modInstalled {
		commandFormat = "/say %s: %s"
	} else {
		commandFormat = "%s: %s"
	}

	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	groups := b.Group()
	groups.Use(Whitelist(group))

	groups.Handle("/players", func(c tele.Context) error {
		resp, err := rconClient.Execute("/players")
		if err != nil {
			return err
		}

		return c.Reply(resp)
	})

	groups.Handle(tele.OnText, func(c tele.Context) error {

		// skip all commands exept predefined
		if c.Message().Text[0] == '/' {
			return nil
		}

		rconClient.Execute(fmt.Sprintf(commandFormat, c.Sender().FirstName, c.Text()))
		return nil
	})

	return b, nil
}

func connectRCON(addr string, pass string) (*rcon.Conn, error) {
	rconClient, err := rcon.Dial(addr, pass)
	if err != nil {
		return nil, err
	}

	return rconClient, nil
}

func fwdFromFactorio(bot *tele.Bot, chatid int64) {

	chat, _ := bot.ChatByID(chatid)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		//echo all logs
		fmt.Println(scanner.Text())

		s := strings.Index(scanner.Text(), "[CHAT]")
		if s == -1 {
			continue
		}
		if strings.Contains(scanner.Text(), "[CHAT] <server>:") {
			continue
		}

		bot.Send(chat, scanner.Text()[s+6:])
		if !modInstalled {
			//skip 1 line which would be our message
			scanner.Scan()
		}
	}

	// if err, die
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

var wg sync.WaitGroup

func main() {
	rconPort := os.Getenv("RCON_PORT")
	rconPass := os.Getenv("RCON_PASS")

	tgToken := os.Getenv("TELEGRAM_TOKEN")
	groupId := os.Getenv("TELEGRAM_GROUP")

	groupIdInt64, err := strconv.ParseInt(groupId, 10, 64)
	if err != nil {
		log.Fatal("Invalid telegram token")
	}

	var rconClient *rcon.Conn

	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		rconClient, err = connectRCON("localhost:"+rconPort, rconPass)
		if err != nil {
			continue
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	modInstalled = checkMod(rconClient)

	bot, err := initBot(tgToken, groupIdInt64, rconClient)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(2)

	go bot.Start()
	go fwdFromFactorio(bot, groupIdInt64)

	wg.Wait()
}
