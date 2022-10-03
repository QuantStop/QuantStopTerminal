package coinbasepro

// Channel is a feed of specific types messages for as specific set of products.
type Channel struct {
	Name       ChannelName `json:"name"`
	ProductIDs []ProductID `json:"product_ids"`
}

type ChannelName string

const (
	// ChannelNameHeartbeat messages for specific products once a second subscribe to the heartbeat channel. Heartbeats
	// also include sequence numbers and last trade ids that can be used to verify no messages were missed.
	ChannelNameHeartbeat ChannelName = "heartbeat"
	// ChannelNameStatus will send all products and currencies on a preset interval.
	ChannelNameStatus ChannelName = "status"
	// ChannelNameTicker provides real-time price updates every time a match happens. It batches updates in case of
	// cascading matches, greatly reducing bandwidth requirements.
	ChannelNameTicker ChannelName = "ticker"
	// ChannelNameLevel2 provides a snapshot of the order book. It guarantees delivery of all updates, which has reduced
	// overhead compared to consuming the full channel.
	ChannelNameLevel2 ChannelName = "level2"
	// ChannelNameFull provides real-time updates on orders and trades. These updates can be applied on to a level 3
	// order book snapshot to maintain an accurate and up-to-date copy of the exchange order book.
	ChannelNameFull ChannelName = "full"
	// ChannelNameUser is a subset of the ChannelName_Full channel that only contains messages that reference the authenticated user.
	ChannelNameUser ChannelName = "user"
	// ChannelNameMatches only includes matches. Note that messages can be dropped from this channel.
	ChannelNameMatches ChannelName = "matches"
)

type MessageType string

const (
	MessageTypeError    MessageType = "error"
	MessageTypeOpen     MessageType = "open"
	MessageTypeDone     MessageType = "done"
	MessageTypeChange   MessageType = "change"
	MessageTypeActivate MessageType = "activate"
	MessageTypeAuction  MessageType = "auction"
	// MessageTypeSubscribe messages are sent to the server and indicate which channels and products to receive.
	MessageTypeSubscribe     MessageType = "subscribe"
	MessageTypeUnsubscribe   MessageType = "unsubscribe"
	MessageTypeSubscriptions MessageType = "subscriptions"
	MessageTypeL2Update      MessageType = "l2update"
	MessageTypeTicker        MessageType = "ticker"
	MessageTypeSnapshot      MessageType = "snapshot"
	MessageTypeMatch         MessageType = "match"
	MessageTypeHeartbeat     MessageType = "heartbeat"
	MessageTypeStatus        MessageType = "status"
	MessageTypeReceived      MessageType = "received"
)

// A SubscriptionRequest describes the products and channels to be provided by the feed.
type SubscriptionRequest struct {
	Type       MessageType   `json:"type"` // must be 'subscribe'
	ProductIDs []ProductID   `json:"product_ids"`
	Channels   []interface{} `json:"channels"`
}

type SubscriptionResponse struct {
	Type     MessageType `json:"type"`
	Channels []Channel   `json:"channels"`
}

type UnsubscriptionRequest struct {
	Type       MessageType   `json:"type"` // must be 'subscribe'
	ProductIDs []ProductID   `json:"product_ids"`
	Channels   []interface{} `json:"channels"`
}

// NewSubscriptionRequest creates the initial message to the server indicating which channels and products to receive.
// This message is mandatory — you will be disconnected if no subscribe has been received within 5 seconds.
//
// There are two ways to specify products to listen for within each channel:
// - By specifying a list of products for all subscribed channels
// - By specifying the product ids for an individual channel
//
// Both mechanisms can be used together in a single subscription request.
func NewSubscriptionRequest(productIDs []ProductID, channelNames []ChannelName, productChannels []Channel) SubscriptionRequest {
	channels := make([]interface{}, 0, len(channelNames)+len(productChannels))
	for _, channelName := range channelNames {
		channels = append(channels, channelName)
	}
	for _, productChannel := range productChannels {
		channels = append(channels, productChannel)
	}
	return SubscriptionRequest{
		Type:       MessageTypeSubscribe,
		ProductIDs: productIDs,
		Channels:   channels,
	}
}

func NewUnsubscriptionRequest(productIDs []ProductID, channelNames []ChannelName, productChannels []Channel) UnsubscriptionRequest {
	channels := make([]interface{}, 0, len(channelNames)+len(productChannels))
	for _, channelName := range channelNames {
		channels = append(channels, channelName)
	}
	for _, productChannel := range productChannels {
		channels = append(channels, productChannel)
	}
	return UnsubscriptionRequest{
		Type:       MessageTypeUnsubscribe,
		ProductIDs: productIDs,
		Channels:   channels,
	}
}
