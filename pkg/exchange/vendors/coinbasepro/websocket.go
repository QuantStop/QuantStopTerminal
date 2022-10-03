package coinbasepro

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx/orderbook"
	"log"
	"strconv"
	"sync"
	"time"
)

func NewFeed() *Feed {
	return &Feed{
		Subscriptions: make(chan SubscriptionRequest, 1),
		Messages:      make(chan []byte),
		Heartbeat:     make(chan HeartbeatMessage),
		Status:        make(chan StatusMessage),
		Ticker:        make(chan TickerMessage),
		TickerBatch:   make(chan TickerMessage),
		Level2:        make(chan L2UpdateMessage),
		Level2Snap:    make(chan L2SnapshotMessage),
		//Level2Batch:   make(chan L2Message),
		//User:	       make(chan HeartbeatMessage),
		Matches: make(chan MatchMessage),
		//Full:	       make(chan HeartbeatMessage),
		//Auction:	   make(chan HeartbeatMessage),
	}
}

type Feed struct {
	Subscriptions chan SubscriptionRequest
	Messages      chan []byte

	Heartbeat   chan HeartbeatMessage
	Status      chan StatusMessage
	Ticker      chan TickerMessage
	TickerBatch chan TickerMessage
	Level2      chan L2UpdateMessage
	Level2Snap  chan L2SnapshotMessage
	//Level2Batch   chan L2Message
	//User          chan User
	Matches chan MatchMessage
	//Full          chan Full
	//Auction       chan Auction
}

type HeartbeatMessage struct {
	Type        string    `json:"type"`
	Sequence    int       `json:"sequence"`
	LastTradeId int       `json:"last_trade_id"`
	ProductId   string    `json:"product_id"`
	Time        time.Time `json:"time"`
}

type StatusMessage struct {
	Type     string `json:"type"`
	Products []struct {
		Id             string      `json:"id"`
		BaseCurrency   string      `json:"base_currency"`
		QuoteCurrency  string      `json:"quote_currency"`
		BaseMinSize    string      `json:"base_min_size"`
		BaseMaxSize    string      `json:"base_max_size"`
		BaseIncrement  string      `json:"base_increment"`
		QuoteIncrement string      `json:"quote_increment"`
		DisplayName    string      `json:"display_name"`
		Status         string      `json:"status"`
		StatusMessage  interface{} `json:"status_message"`
		MinMarketFunds string      `json:"min_market_funds"`
		MaxMarketFunds string      `json:"max_market_funds"`
		PostOnly       bool        `json:"post_only"`
		LimitOnly      bool        `json:"limit_only"`
		CancelOnly     bool        `json:"cancel_only"`
		FxStablecoin   bool        `json:"fx_stablecoin"`
	} `json:"products"`
	Currencies []struct {
		Id            string      `json:"id"`
		Name          string      `json:"name"`
		MinSize       string      `json:"min_size"`
		Status        string      `json:"status"`
		StatusMessage interface{} `json:"status_message"`
		MaxPrecision  string      `json:"max_precision"`
		ConvertibleTo []string    `json:"convertible_to"`
		Details       struct {
		} `json:"details"`
	} `json:"currencies"`
}

type TickerMessage struct {
	Type      string    `json:"type"`
	Sequence  int       `json:"sequence"`
	ProductId string    `json:"product_id"`
	Price     string    `json:"price"`
	Open24H   string    `json:"open_24h"`
	Volume24H string    `json:"volume_24h"`
	Low24H    string    `json:"low_24h"`
	High24H   string    `json:"high_24h"`
	Volume30D string    `json:"volume_30d"`
	BestBid   string    `json:"best_bid"`
	BestAsk   string    `json:"best_ask"`
	Side      string    `json:"side"`
	Time      time.Time `json:"time"`
	TradeId   int       `json:"trade_id"`
	LastSize  string    `json:"last_size"`
}

type L2SnapshotMessage struct {
	Type      string     `json:"type"`
	ProductId string     `json:"product_id"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
}

type L2UpdateMessage struct {
	Type      string     `json:"type"`
	ProductId string     `json:"product_id"`
	Time      time.Time  `json:"time"`
	Changes   [][]string `json:"changes"`
}

type MatchMessage struct {
	Type         string    `json:"type"`
	TradeId      int       `json:"trade_id"`
	Sequence     int       `json:"sequence"`
	MakerOrderId string    `json:"maker_order_id"`
	TakerOrderId string    `json:"taker_order_id"`
	Time         time.Time `json:"time"`
	ProductId    string    `json:"product_id"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	Side         string    `json:"side"`
}

type ChangeMessage struct {
	Type      string    `json:"type"`
	Time      time.Time `json:"time"`
	Sequence  int       `json:"sequence"`
	OrderId   string    `json:"order_id"`
	ProductId string    `json:"product_id"`
	NewSize   string    `json:"new_size"`
	OldSize   string    `json:"old_size"`
	Price     string    `json:"price"`
	Side      string    `json:"side"`
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Watch provides a feed of real-time market data updates for orders and trades.
func (c *CoinbasePro) Watch(shutdown chan struct{}, wg *sync.WaitGroup, subscriptionRequest SubscriptionRequest, feed *Feed) (*orderbook.Orderbook, error) {
	var err error

	// try dialing the connection
	if c.Conn, err = c.Websocket.Dial(); err != nil {
		return nil, err
	}

	// subscription request must be sent within 5 seconds of open or socket will auto-close
	if err = c.Conn.WriteJSON(subscriptionRequest); err != nil {
		return nil, err
	}

	// create the orderbook
	book := orderbook.NewOrderbook()

	// run the feed and return the book
	c.watch(shutdown, wg, feed, book)
	return book, nil
}

func (c *CoinbasePro) watch(shutdown chan struct{}, wg *sync.WaitGroup, feed *Feed, orderbook *orderbook.Orderbook) {

	log.Println("coinbase websocket watch() routine activated, starting up feeds ...")
	wg.Add(1)

	// read messages from coinbase
	go func() {
		defer wg.Done()
		for {
			select {
			case <-shutdown:
				log.Println("coinbase websocket received context.Done() request, returning from watch() routine ...")
				return
			default:
				_, message, err := c.Conn.ReadMessage()
				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("coinbase websocket read error: %v", err)
						return
					}
				}
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				decodeMessage(message, feed, orderbook)

			}
		}
	}()

	log.Println("coinbase websocket feeds started.")

}

func decodeMessage(message []byte, feed *Feed, orderbook *orderbook.Orderbook) {
	var messageInterface map[string]interface{}
	if err := json.Unmarshal(message, &messageInterface); err != nil {
		log.Printf("coinbase websocket decodeMessage error unmarshalling incoming bytes: %v", err)
		return
	}

	msgType := MessageType(messageInterface["type"].(string))
	switch msgType {
	case MessageTypeError:

	case MessageTypeL2Update:
		l2up := L2UpdateMessage{}
		if err := json.Unmarshal(message, &l2up); err != nil {
			log.Printf("coinbase websocket decodeMessage error unmarshalling %s: %v", MessageTypeL2Update, err)
			return
		}
		for _, change := range l2up.Changes {
			price, err := strconv.ParseFloat(change[1], 0)
			if err != nil {
				log.Printf("coinbase error converting bid price to orderbook Order: %v", err)
			}
			size, err := strconv.ParseFloat(change[2], 0)
			if err != nil {
				log.Printf("coinbase error converting bid size to orderbook Order: %v", err)
			}
			if size != 0.0 {
				order := orderbook.NewOrder()
				if change[0] == "buy" {
					order.BidOrAsk = true //bid
				} else {
					order.BidOrAsk = false //ask
				}
				order.Volume = size
				orderbook.Add(price, order)
			} else {
				if change[0] == "buy" {
					orderbook.DeleteBidLimit(price)
				} else {
					orderbook.DeleteAskLimit(price)
				}
			}
		}
		feed.Level2 <- l2up
	case MessageTypeSnapshot:
		l2snap := L2SnapshotMessage{}
		if err := json.Unmarshal(message, &l2snap); err != nil {
			log.Printf("coinbase websocket decodeMessage error unmarshalling %s: %v", MessageTypeSnapshot, err)
			return
		}
		for _, bid := range l2snap.Bids {
			price, err := strconv.ParseFloat(bid[0], 0)
			if err != nil {
				log.Printf("coinbase error converting bid price to orderbook Order: %v", err)
			}
			size, err := strconv.ParseFloat(bid[1], 0)
			if err != nil {
				log.Printf("coinbase error converting bid size to orderbook Order: %v", err)
			}
			order := orderbook.NewOrder()
			order.BidOrAsk = true //bid
			order.Volume = size
			orderbook.Add(price, order)
		}
		for _, ask := range l2snap.Asks {
			price, err := strconv.ParseFloat(ask[0], 0)
			if err != nil {
				log.Printf("coinbase error converting ask price to orderbook Order: %v", err)
			}
			size, err := strconv.ParseFloat(ask[1], 0)
			if err != nil {
				log.Printf("coinbase error converting ask size to orderbook Order: %v", err)
			}
			order := orderbook.NewOrder()
			order.BidOrAsk = false //ask
			order.Volume = size
			orderbook.Add(price, order)
		}
		feed.Level2Snap <- l2snap
	case MessageTypeHeartbeat:
		hbMsg := HeartbeatMessage{}
		if err := json.Unmarshal(message, &hbMsg); err != nil {
			log.Printf("coinbase websocket decodeMessage error unmarshalling %s: %v", MessageTypeHeartbeat, err)
			return
		}
		feed.Heartbeat <- hbMsg
	case MessageTypeMatch:
		match := MatchMessage{}
		if err := json.Unmarshal(message, &match); err != nil {
			log.Printf("coinbase websocket decodeMessage error unmarshalling %s: %v", MessageTypeMatch, err)
			return
		}
		feed.Matches <- match
	case MessageTypeSubscriptions:

	case MessageTypeTicker:

	}
}
