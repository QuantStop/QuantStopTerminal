package websocket

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	qse "github.com/quantstop/quantstopterminal/pkg/exchange"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
	"github.com/quantstop/quantstopterminal/pkg/exchange/vendors/coinbasepro"

	"github.com/quantstop/quantstopterminal/internal/database"
	"github.com/quantstop/quantstopterminal/internal/database/models"
	"github.com/quantstop/quantstopterminal/internal/log"
	"github.com/quantstop/quantstopterminal/internal/webserver/write"
	"net/http"
	"sync"
)

const (
	publish     = "publish"
	subscribe   = "subscribe"
	unsubscribe = "unsubscribe"
	getClients  = "getClients"
	getSubs     = "getSubs"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {

	// db holds a pointer to the database service for access to all databases
	db *database.Database

	// Shutdown channel to stop the Hub
	Shutdown chan struct{}

	// Registered clients.
	clients map[*Client]bool

	// Client subscriptions
	Subscriptions []*Subscription

	// Register channel for requests from the clients.
	Register chan *Client

	// Unregister channel for requests from clients.
	Unregister chan *Client
}

// Subscription type holds a single client <-> exchange connection
type Subscription struct {
	ExchangeClient qsx.IExchange
	Client         *Client
	Shutdown       chan struct{}
}

// Message is the type for a valid message from a client
type Message struct {
	Action     string `json:"action"`
	ExchangeID string `json:"exchange_id"`
	Message    string `json:"message"`
}

type MessageResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func NewHub(db *database.Database, shutdown chan struct{}) (*Hub, error) {
	if db == nil {
		return nil, errors.New("cannot create new hub, database is nil")
	}
	if shutdown == nil {
		return nil, errors.New("cannot create new hub, shutdown chan is nil")
	}
	return &Hub{
		db:            db,
		clients:       make(map[*Client]bool),
		Subscriptions: make([]*Subscription, 0),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Shutdown:      shutdown,
	}, nil
}

func (h *Hub) Run(serviceWG *sync.WaitGroup) {
	serviceWG.Add(1)
	go func() {
		for {
			select {
			case <-h.Shutdown:
				h.UnsubscribeAll()
				serviceWG.Done()
				return
			case client := <-h.Register:
				h.clients[client] = true
			case client := <-h.Unregister:
				if _, ok := h.clients[client]; ok {
					h.Unsubscribe(client)
					delete(h.clients, client)
					close(client.send)
				}
			}
		}
	}()

}

func (h *Hub) Send(client *Client, message []byte) {
	client.send <- message
}

func (h *Hub) ProcessMessage(client *Client, messageType int, payload []byte) *Hub {
	m := Message{}
	if err := json.Unmarshal(payload, &m); err != nil {
		msgRes := MessageResponse{
			Type:    "error",
			Message: "error: invalid payload",
		}
		res, err := json.Marshal(msgRes)
		if err != nil {
			log.Error(log.Webserver, err)
		}
		h.Send(client, res)
		return h
	}
	log.Debugf(log.Webserver, "Coinbasepro Message: %v", m)

	switch m.Action {
	case publish:
		//h.Publish(m.ExchangeID, []byte(m.Message))
		break

	case subscribe:
		h.Subscribe(client, m.ExchangeID, m.Message)
		break

	case unsubscribe:
		h.Unsubscribe(client)
		break

	case getClients:

		break

	case getSubs:

		break

	default:
		msgRes := MessageResponse{
			Type:    "error",
			Message: "error: unrecognized action",
		}
		res, err := json.Marshal(msgRes)
		if err != nil {
			log.Error(log.Webserver, err)
		}
		h.Send(client, res)
		break
	}

	return h
}

func (h *Hub) Subscribe(client *Client, exchange string, product string) {

	log.Debugf(log.Webserver, "Websocket Hub | Coinbasepro %v requested a subscription to %v for %v", client.ID, exchange, product)

	// check for existing subscription
	for _, sub := range h.Subscriptions {
		// if found, unsubscribe and create new subscription with new request
		if sub.Client.ID == client.ID {
			h.Unsubscribe(client)
		}
	}

	// create new subscription & add client to that subscription
	newSubscription := &Subscription{
		Client:   client,
		Shutdown: make(chan struct{}),
	}
	h.Subscriptions = append(h.Subscriptions, newSubscription)

	// run the exchange service
	go h.RunSubscriptionService(newSubscription, exchange, product, client)

}

func (h *Hub) Unsubscribe(client *Client) {
	log.Debugf(log.Webserver, "Websocket Hub | Coinbasepro %v requested unsubscribe", client.ID)
	for subIndex, sub := range h.Subscriptions {
		if sub.Client.ID == client.ID {
			close(sub.Shutdown)
			h.Subscriptions = append(h.Subscriptions[:subIndex], h.Subscriptions[subIndex+1:]...)
		}
	}
}

func (h *Hub) UnsubscribeAll() {
	log.Debugf(log.Webserver, "Websocket Hub | unsubscribing all clients")
	for subIndex, sub := range h.Subscriptions {
		close(sub.Shutdown)
		h.Subscriptions = append(h.Subscriptions[:subIndex], h.Subscriptions[subIndex+1:]...)
	}
}

func (h *Hub) RunSubscriptionService(sub *Subscription, exchangeName, product string, client *Client) {

	log.Debugf(log.Webserver, "Websocket Hub | %v subscription service starting ...", client.ID)

	e := models.Exchange{}
	err := e.GetExchangeByName(h.db.CoreDB.SQL, qsx.Name(exchangeName))
	if err != nil {
		log.Error(log.Webserver, err)
		return
	}

	// Create a client instance
	sub.ExchangeClient, err = qse.NewExchange(
		qsx.Name(exchangeName),
		&qsx.Config{
			Auth: &qsx.Auth{
				Key:        e.AuthKey,
				Passphrase: e.AuthPassphrase,
				Secret:     e.AuthSecret,
				Token:      nil,
			},
			Sandbox: false,
		},
	)

	if err != nil {
		log.Error(log.Webserver, err)
		return
	}

	// create a new subscription request
	// todo: feed needs to be analogous to exchange (common feed)
	feed := coinbasepro.NewFeed()

	wg := sync.WaitGroup{}
	wg.Add(4)

	// Start api client feed
	_, err = sub.ExchangeClient.WatchFeed(sub.Shutdown, &wg, product, feed)
	if err != nil {
		return
	}

	// Loop on Heartbeat channel
	go func() {
		defer wg.Done()
		for message := range feed.Heartbeat {
			select {
			case <-sub.Shutdown:
				log.Debugf(log.Webserver, "HUB | Coinbasepro %v subscription feed.Heartbeat service shutdown requested ...", client.ID)
				return
			default:
				res, err := json.Marshal(message)
				if err != nil {
					log.Error(log.Webserver, err)
				}
				client.send <- res
				//log.Debugf(log.TraderLogger, "%s | %s", message.Type, message.Time.String())
				continue
			}
		}
	}()

	// Loop on L2Channel channel
	go func() {
		defer wg.Done()
		for message := range feed.Level2 {
			select {
			case <-sub.Shutdown:
				log.Debugf(log.Webserver, "HUB | Coinbasepro %v subscription feed.Level2 service shutdown requested ...", client.ID)
				return
			default:
				res, err := json.Marshal(message)
				if err != nil {
					log.Error(log.Webserver, err)
				}
				client.send <- res
				//log.Debugf(log.TraderLogger, "%s | %s", message.Type, message.Time.String())
				continue
			}
		}
	}()

	// Loop on L2ChannelSnapshot channel
	go func() {
		defer wg.Done()
		for message := range feed.Level2Snap {
			select {
			case <-sub.Shutdown:
				log.Debugf(log.Webserver, "HUB | Coinbasepro %v subscription feed.Level2Snap service shutdown requested ...", client.ID)
				return
			default:
				res, err := json.Marshal(message)
				if err != nil {
					log.Error(log.Webserver, err)
				}
				client.send <- res
				//log.Debugf(log.TraderLogger, "%s | %s", message.Type, message.ProductId)
				continue
			}
		}
	}()

	// Loop on Matches channel
	go func() {
		defer wg.Done()
		for message := range feed.Matches {
			select {
			case <-sub.Shutdown:
				log.Debugf(log.Webserver, "HUB | Coinbasepro %v subscription feed.Matches service shutdown requested ...", client.ID)
				return
			default:
				res, err := json.Marshal(message)
				if err != nil {
					log.Error(log.Webserver, err)
				}
				client.send <- res
				//log.Debugf(log.TraderLogger, "%s | %s", message.Type, message.Time.String())
				continue
			}
		}
	}()

	// Wait for all go-routines to finish
	log.Debugf(log.Webserver, "Websocket Hub | %v subscription service started.", client.ID)
	wg.Wait()
	<-sub.Shutdown
	log.Debugf(log.Webserver, "Websocket Hub | %v subscription service shutdown.", client.ID)
}

var upgrader = websocket.Upgrader{
	//https://github.com/kataras/neffos/issues/11#issuecomment-520689681
	// todo this allows all origins ...
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(log.Webserver, err)
		write.Error(err)
		return
	}
	client := &Client{
		ID:   uuid.Must(uuid.NewRandom()).String(),
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.Register <- client

	// greet the new client
	msgRes := MessageResponse{
		Type:    "welcome",
		Message: "Welcome to QuantstopTerminal Websocket Server: Your ID is " + client.ID,
	}
	res, err := json.Marshal(msgRes)
	if err != nil {
		log.Error(log.Webserver, err)
		write.Error(err)
		return
	}
	hub.Send(client, res)

	log.Debugf(log.Webserver, "HUB | %v clients connected.", len(hub.clients))

	// Allow collection of memory referenced by the caller by doing all work in new goroutines.
	go client.writePump()
	go client.readPump()
}
