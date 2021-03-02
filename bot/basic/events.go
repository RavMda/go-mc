package basic

import (
	"github.com/google/uuid"

	"github.com/RavMda/go-mc/bot"
	"github.com/RavMda/go-mc/chat"
	"github.com/RavMda/go-mc/data/packetid"
	pk "github.com/RavMda/go-mc/net/packet"
)

type EventsListener struct {
	GameStart    func(c *bot.Client) error
	ChatMsg      func(cl *bot.Client, c chat.Message, pos byte, uuid uuid.UUID) error
	Disconnect   func(c *bot.Client, reason chat.Message) error
	HealthChange func(c *bot.Client, health float32) error
	Death        func(c *bot.Client) error
}

func (e EventsListener) Attach(c *bot.Client) {
	c.Events.AddListener(
		bot.PacketHandler{Priority: 64, ID: packetid.Login, F: e.onJoinGame},
		bot.PacketHandler{Priority: 64, ID: packetid.ChatClientbound, F: e.onChatMsg},
		bot.PacketHandler{Priority: 64, ID: packetid.KickDisconnect, F: e.onDisconnect},
		bot.PacketHandler{Priority: 64, ID: packetid.UpdateHealth, F: e.onUpdateHealth},
	)
}

func (e *EventsListener) onJoinGame(c *bot.Client, _ pk.Packet) error {
	if e.GameStart != nil {
		return e.GameStart(c)
	}
	return nil
}

func (e *EventsListener) onDisconnect(c *bot.Client, p pk.Packet) error {
	if e.Disconnect != nil {
		var reason chat.Message
		if err := p.Scan(&reason); err != nil {
			return Error{err}
		}
		return e.Disconnect(c, reason)
	}
	return nil
}

func (e *EventsListener) onChatMsg(c *bot.Client, p pk.Packet) error {
	if e.ChatMsg != nil {
		var msg chat.Message
		var pos pk.Byte
		var sender pk.UUID

		if err := p.Scan(&msg, &pos, &sender); err != nil {
			return Error{err}
		}

		return e.ChatMsg(c, msg, byte(pos), uuid.UUID(sender))
	}
	return nil
}

func (e *EventsListener) onUpdateHealth(c *bot.Client, p pk.Packet) error {
	if e.ChatMsg != nil {
		var health pk.Float
		var food pk.VarInt
		var foodSaturation pk.Float

		if err := p.Scan(&health, &food, &foodSaturation); err != nil {
			return Error{err}
		}
		if e.HealthChange != nil {
			if err := e.HealthChange(c, float32(health)); err != nil {
				return err
			}
		}
		if e.Death != nil && health <= 0 {
			if err := e.Death(c); err != nil {
				return err
			}
		}
	}
	return nil
}
