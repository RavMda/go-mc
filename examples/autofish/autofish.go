package main

import (
	"log"
	"time"

	"github.com/RavMda/go-mc/bot/basic"
	"github.com/RavMda/go-mc/data/packetid"
	pk "github.com/RavMda/go-mc/net/packet"

	"github.com/google/uuid"

	"github.com/RavMda/go-mc/bot"
	"github.com/RavMda/go-mc/chat"
	_ "github.com/RavMda/go-mc/data/lang/en-us"
	"github.com/mattn/go-colorable"
)

const timeout = 45

var (
	c *bot.Client
	p *basic.Player

	watch chan time.Time
)

func main() {
	log.SetOutput(colorable.NewColorableStdout())
	c = bot.NewClient()
	p = basic.NewPlayer(c, basic.DefaultSettings)

	//Register event handlers
	basic.EventsListener{
		GameStart:  onGameStart,
		ChatMsg:    onChatMsg,
		Disconnect: onDisconnect,
		Death:      onDeath,
	}.Attach(c)
	c.Events.AddListener(soundListener)

	//Login
	err := c.JoinServer("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Login success")

	//JoinGame
	err = c.HandleGame()
	if err != nil {
		log.Fatal(err)
	}
}

func onDeath(c *bot.Client) error {
	log.Println("Died and Respawned")
	// If we exclude Respawn(...) then the player won't press the "Respawn" button upon death
	return p.Respawn()
}

func onGameStart(c *bot.Client) error {
	log.Println("Game start")

	watch = make(chan time.Time)
	go watchDog()

	return UseItem(0)
}

var soundListener = bot.PacketHandler{
	ID:       packetid.NamedSoundEffect,
	Priority: 0,
	F: func(c *bot.Client, p pk.Packet) error {
		var (
			SoundName     pk.Identifier
			SoundCategory pk.VarInt
			X, Y, Z       pk.Int
			Volume, Pitch pk.Float
		)
		if err := p.Scan(&SoundName, &SoundCategory, &X, &Y, &Z, &Volume, &Pitch); err != nil {
			return err
		}
		return onSound(string(SoundName), int(SoundCategory), float64(X)/8, float64(Y)/8, float64(Z)/8, float32(Volume), float32(Pitch))
	},
}

func UseItem(hand int32) error {
	return c.Conn.WritePacket(pk.Marshal(
		packetid.UseItem,
		pk.VarInt(hand),
	))
}

//goland:noinspection SpellCheckingInspection
func onSound(name string, category int, x, y, z float64, volume, pitch float32) error {
	if name == "entity.fishing_bobber.splash" {
		if err := UseItem(0); err != nil { //retrieve
			return err
		}
		log.Println("gra~")
		time.Sleep(time.Millisecond * 300)
		if err := UseItem(0); err != nil { //throw
			return err
		}
		watch <- time.Now()
	}
	return nil
}

func onChatMsg(cl *bot.Client, c chat.Message, pos byte, uuid uuid.UUID) error {
	log.Println("Chat:", c)
	return nil
}

func onDisconnect(cl *bot.Client, c chat.Message) error {
	log.Println("Disconnect:", c)
	return nil
}

func watchDog() {
	to := time.NewTimer(time.Second * timeout)
	for {
		select {
		case <-watch:
		case <-to.C:
			log.Println("rethrow")
			if err := UseItem(0); err != nil {
				panic(err)
			}
		}
		to.Reset(time.Second * timeout)
	}
}
