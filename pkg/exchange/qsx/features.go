package qsx

// ExchangeFeatures holds information about what services/features the exchange provides
type ExchangeFeatures struct {

	// HasCrypto returns true if the exchange has cryptocurrency support
	HasCrypto bool

	// HasWebsocket returns true if the exchange has websocket streaming support
	HasWebsocket bool

	// HasOptions returns true if the exchange has support for options data
	HasOptions bool
}
