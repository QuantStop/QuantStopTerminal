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
	GetSQL(dbName string) (*sql.DB, error)
	GetExchange(string) qsx.IExchange
	GetSupportedExchangesList() []string
}
