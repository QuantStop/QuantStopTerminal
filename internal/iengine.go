package internal

import (
	"database/sql"
	"github.com/quantstop/quantstopexchange/qsx"
)

type IEngine interface {
	GetUptime() string
	SetConfig(string, string) error
	GetSubsystemsStatus() map[string]bool
	SetSubsystem(subSystemName string, enable bool) error
	GetVersion() map[string]string
	GetCoreSQL() (*sql.DB, error)
	GetCoinbaseSQL() (*sql.DB, error)
	GetTDAmeritradeSQL() (*sql.DB, error)
	GetExchange(string) qsx.IExchange
	GetSupportedExchangesList() []string
}
