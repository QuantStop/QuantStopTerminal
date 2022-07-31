package engine

import (
	"github.com/quantstop/quantstopexchange"
	"github.com/quantstop/quantstopexchange/qsx"
	"github.com/quantstop/quantstopterminal/internal/database/models"
	"github.com/quantstop/quantstopterminal/internal/log"
	"sync"
)

type ExchangeManager struct {
	Subsystem
	Exchanges map[string]qsx.IExchange
}

func (e *ExchangeManager) init(bot *Engine, name string) error {
	if err := e.Subsystem.init(bot, name, true); err != nil {
		return err
	}
	e.Exchanges = make(map[string]qsx.IExchange)
	e.initialized = true
	log.Debugln(log.ExchangeManager, e.name+MsgSubsystemInitialized)
	return nil
}

// start sets up the exchange manager subsystem to maintain external connections to each exchange
func (e *ExchangeManager) start(wg *sync.WaitGroup) (err error) {
	if err = e.Subsystem.start(wg); err != nil {
		return err
	}

	for _, name := range qsx.SupportedExchanges {
		//log.Debugln(log.Global, name)
		exchange := models.Exchange{}
		err = exchange.GetExchangeByName(e.bot.DatabaseSubsystem.coreDatabase.SQL, name)
		if err != nil {
			log.Error(log.Global, err)
			return err
		}
		ex, err := quantstopexchange.NewExchange(name, &qsx.Config{
			Auth: &qsx.Auth{
				Key:        exchange.AuthKey,
				Passphrase: exchange.AuthPassphrase,
				Secret:     exchange.AuthSecret,
				Token:      nil,
			},
			Sandbox: false,
		})
		if err != nil {
			log.Error(log.Global, err)
			return err
		}

		e.Exchanges[exchange.Name] = ex
	}

	log.Debugln(log.ExchangeManager, e.name+MsgSubsystemStarted)
	e.started = true
	return nil
}

// stop attempts to shut down the subsystem
func (e *ExchangeManager) stop() error {
	if err := e.Subsystem.stop(); err != nil {
		return err
	}

	e.started = false

	close(e.shutdown)
	log.Debugln(log.ExchangeManager, e.name+MsgSubsystemShutdown)
	return nil
}
