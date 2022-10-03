package coinbasepro

import (
	"context"
	"fmt"
	"github.com/quantstop/quantstopterminal/pkg/exchange/qsx"
)

// Withdrawal represents a movement of CurrencyName out of accounts to both external and internal destinations.
// Withdrawals are implemented as Transfers, but I assume this was confusing/naive, as the documentation
// rebrands Transfers as Withdrawals. I have followed the hint and done the same.
type Withdrawal struct {
	// Details provides more fine-grained information describing the Withdrawal
	Details WithdrawalDetails `json:"details"`
	// CanceledAt is the time of cancellation, if the Withdrawal was canceled
	CanceledAt *Time `json:"canceled_at"`
	// CompletedAt is the time of completion, if the Withdrawal was completed
	CompletedAt *Time `json:"completed_at"`
	// CreatedAt is the time of creation
	CreatedAt Time `json:"created_at"`
	// CreatedAt is the time the transfer was processed
	ProcessedAt *Time `json:"processed_at"`
	// AccountID identifies the Account to which the Withdrawal applies
	AccountID string `json:"account_id"`
	// Amount is the amount of the Withdrawal
	// TODO: in what currency?
	Amount   float64      `json:"amount,string"`
	Currency CurrencyName `json:"currency"`
	ID       string       `json:"id"`
	// Type identifies the type of the Withdrawal (`withdraw` or `internal_withdraw`)
	Type WithdrawalType `json:"type"`
	// UserID that initiated the Withdrawal
	UserID    string `json:"user_id"`
	UserNonce string `json:"user_nonce"`
}

// WithdrawalDetails is not well documented; until proven or requested otherwise, I will simply treat the details
// as free form annotations or labels.
type WithdrawalDetails map[string]interface{}

type WithdrawalType string

const (
	WithdrawalTypeWithdraw WithdrawalType = "withdraw"
	WithdrawalTypeInternal WithdrawalType = "internal_withdraw"
)

type WithdrawalFilter struct {
	// ProfileID limits the list of Withdrawals to the ProfileID. By default, Withdrawals retrieves Withdrawals for the default profile.
	ProfileID string `json:"profile_id"`
	// Type identifies the type of the Withdrawal (`withdraw` or `internal_withdraw`)
	Type WithdrawalType `json:"type"`
}

func (d WithdrawalFilter) Params() []string {
	var params []string
	if d.ProfileID != "" {
		params = append(params, fmt.Sprintf("profile_id=%s", d.ProfileID))
	}
	if d.Type != "" {
		params = append(params, fmt.Sprintf("type=%s", d.Type))
	}
	return params
}

type Withdrawals struct {
	Withdrawals []*Withdrawal `json:"withdrawals"`
	Page        *Pagination   `json:"page"`
}

type WithdrawalCancelCode int

// 1xxx cancel codes are for fiat transfers and 2xxx cancel codes are for crypto transfer. Please make changes accordingly before retrying the withdrawal.
const (
	WithdrawalCancelCodeDefaultError                             WithdrawalCancelCode = 0
	WithdrawalCancelCodeMaxExceeded                              WithdrawalCancelCode = 1000
	WithdrawalCancelCodeZeroAmount                               WithdrawalCancelCode = 1001
	WithdrawalCancelCodeAccountNotAllowed                        WithdrawalCancelCode = 1002
	WithdrawalCancelCodePaymentMethodNotAllowed                  WithdrawalCancelCode = 1003
	WithdrawalCancelCodeCurrencyAndPaymentMethodNotAllowed       WithdrawalCancelCode = 1004
	WithdrawalCancelCodeAmountExceedsAccountFunds                WithdrawalCancelCode = 1005
	WithdrawalCancelCodeAmountMustBeAtLeastOne                   WithdrawalCancelCode = 1006
	WithdrawalCancelCodeAmountTooSmall                           WithdrawalCancelCode = 1007
	WithdrawalCancelCodeNoRecurringTransfersWithPaymentMethod    WithdrawalCancelCode = 1008
	WithdrawalCancelCodeCurrencyDoesNotMatchAccountCurrency      WithdrawalCancelCode = 1009
	WithdrawalCancelCodePaymentMethodUnsupported                 WithdrawalCancelCode = 1010
	WithdrawalCancelCodeWithdrawalRateLimitExceeded              WithdrawalCancelCode = 1011
	WithdrawalCancelCodeAmountExceedsMaximumAccountBalance       WithdrawalCancelCode = 1012
	WithdrawalCancelCodeNegativeAmount                           WithdrawalCancelCode = 1013
	WithdrawalCancelCodeNoTagNameProvided                        WithdrawalCancelCode = 2000
	WithdrawalCancelCodeAmountExceedsSendLimits                  WithdrawalCancelCode = 2004
	WithdrawalCancelCodeMaxSendsExceeded                         WithdrawalCancelCode = 2005
	WithdrawalCancelCodeSendAmountTooSmallForOnBlockchain        WithdrawalCancelCode = 2007
	WithdrawalCancelCodeTwoStepVerificationCodeRequired          WithdrawalCancelCode = 2008
	WithdrawalCancelCodeCurrencyRequiresTagName                  WithdrawalCancelCode = 2009
	WithdrawalCancelCodeInvalidAmount                            WithdrawalCancelCode = 2010
	WithdrawalCancelCodeCurrencyTemporarilyDisabled              WithdrawalCancelCode = 2011
	WithdrawalCancelCodeAmountExceedsCurrencyMaxWithdrawalAmount WithdrawalCancelCode = 2012
	WithdrawalCancelCodeAmountExceedsCurrencyMaxSendAmount       WithdrawalCancelCode = 2013
	WithdrawalCancelCodeSendFromFiatAccountsTemporarilyDisabled  WithdrawalCancelCode = 2014
	WithdrawalCancelCodePaymentRequestExpired                    WithdrawalCancelCode = 2015
	WithdrawalCancelCodeSendFromAccountNotAllowed                WithdrawalCancelCode = 2016
	WithdrawalCancelCodeUnableToSendToAddress                    WithdrawalCancelCode = 2017
	WithdrawalCancelCodeRecipientAddressNotWhitelisted           WithdrawalCancelCode = 2018
	WithdrawalCancelCodeRecipientAddressWhitelistPending         WithdrawalCancelCode = 2020
	WithdrawalCancelCodeUnableToSendToUser                       WithdrawalCancelCode = 2021
	WithdrawalCancelCodeUnableToSendToSelf                       WithdrawalCancelCode = 2022
	WithdrawalCancelCodeSendRateExceeded                         WithdrawalCancelCode = 2023
	WithdrawalCancelCodeInvalidEmailOrNetworkAddress             WithdrawalCancelCode = 2024
	WithdrawalCancelCodeAccountDoesNotSupportCurrency            WithdrawalCancelCode = 2025
)

func (w WithdrawalCancelCode) String() string {
	switch w {
	case WithdrawalCancelCodeDefaultError:
		return "default error"
	case WithdrawalCancelCodeMaxExceeded:
		return "transaction exceeds transaction limit"
	case WithdrawalCancelCodeZeroAmount:
		return "amount must be greater than 0"
	case WithdrawalCancelCodeAccountNotAllowed:
		return "account does not support withdrawal"
	case WithdrawalCancelCodePaymentMethodNotAllowed:
		return "payment method does not support withdrawal"
	case WithdrawalCancelCodeCurrencyAndPaymentMethodNotAllowed:
		return "cannot withdraw this currency with this payment method"
	case WithdrawalCancelCodeAmountExceedsAccountFunds:
		return "withdrawal amount exceeds funds in account"
	case WithdrawalCancelCodeAmountMustBeAtLeastOne:
		return "withdrawal amount must be at least 1.00"
	case WithdrawalCancelCodeAmountTooSmall:
		return "withdrawal amount too small"
	case WithdrawalCancelCodeNoRecurringTransfersWithPaymentMethod:
		return "payment method cannot be used with recurring transfers"
	case WithdrawalCancelCodeCurrencyDoesNotMatchAccountCurrency:
		return "withdrawal currency does not match account currency"
	case WithdrawalCancelCodePaymentMethodUnsupported:
		return "payment method unsupported"
	case WithdrawalCancelCodeWithdrawalRateLimitExceeded:
		return "withdrawal rate limit exceeded: try again in a few hours"
	case WithdrawalCancelCodeAmountExceedsMaximumAccountBalance:
		return "amount would exceed maximum account balance"
	case WithdrawalCancelCodeNegativeAmount:
		return "amount must be positive"
	case WithdrawalCancelCodeNoTagNameProvided:
		return "warning: with no tag name, recipient may lose funds: confirm that recipient does not require tag name"
	case WithdrawalCancelCodeAmountExceedsSendLimits:
		return "amount would exceed send limits: try a smaller amount or try again later"
	case WithdrawalCancelCodeMaxSendsExceeded:
		return "maximum number of sends per hour exceeded: contact support if you require a higher limit or try again later"
	case WithdrawalCancelCodeSendAmountTooSmallForOnBlockchain:
		return "send amount is below the minimum amount required to send on-blockchain"
	case WithdrawalCancelCodeTwoStepVerificationCodeRequired:
		return "two-step verification code required to complete this request: resend request with CB-2FA-Token header"
	case WithdrawalCancelCodeCurrencyRequiresTagName:
		return "withdrawal currency requires tag name"
	case WithdrawalCancelCodeInvalidAmount:
		return "amount is invalid"
	case WithdrawalCancelCodeCurrencyTemporarilyDisabled:
		return "withdrawal of this currency is temporarily disabled"
	case WithdrawalCancelCodeAmountExceedsCurrencyMaxWithdrawalAmount:
		return "withdrawal amount exceeds maximum withdrawal amount for currency"
	case WithdrawalCancelCodeAmountExceedsCurrencyMaxSendAmount:
		return "withdrawal amount exceeds maximum send amount for currency"
	case WithdrawalCancelCodeSendFromFiatAccountsTemporarilyDisabled:
		return "send from fiat accounts is temporarily disabled: try again later"
	case WithdrawalCancelCodePaymentRequestExpired:
		return "payment request has expired"
	case WithdrawalCancelCodeSendFromAccountNotAllowed:
		return "send from this account not allowed"
	case WithdrawalCancelCodeUnableToSendToAddress:
		return "unable to send to this address"
	case WithdrawalCancelCodeRecipientAddressNotWhitelisted:
		return "recipient address is not whitelisted"
	case WithdrawalCancelCodeRecipientAddressWhitelistPending:
		return "recipient address whitelist pending: 48 hour hold: try again later"
	case WithdrawalCancelCodeUnableToSendToUser:
		return "unable to send to this user"
	case WithdrawalCancelCodeUnableToSendToSelf:
		return "cannot send from an account to itself"
	case WithdrawalCancelCodeSendRateExceeded:
		return "too many sends, too quickly: wait for some transactions to confirm before sending more"
	case WithdrawalCancelCodeInvalidEmailOrNetworkAddress:
		return "invalid email or network address"
	case WithdrawalCancelCodeAccountDoesNotSupportCurrency:
		return "account does not support this currency"
	default:
		return fmt.Sprintf("unknown withdrawal cancel code(%d)", w)
	}
}

// PaymentMethodWithdrawalSpec withdraw funds using a PaymentMethod
type PaymentMethodWithdrawalSpec struct {
	Amount          float64      `json:"amount,string"`
	Currency        CurrencyName `json:"currency"`
	PaymentMethodID string       `json:"payment_method_id"`
}

// CoinbaseAccountWithdrawalSpec creates payout to a CoinbaseAccount
type CoinbaseAccountWithdrawalSpec struct {
	Amount            float64      `json:"amount,string"`
	Currency          CurrencyName `json:"currency"`
	CoinbaseAccountID string       `json:"coinbase_account_id"`
}

type CryptoAddressWithdrawalSpec struct {
	// AddNetworkFeeToTotal indicates that the network fee should be added to the amount
	// By default, network fees are deducted from the amount
	AddNetworkFeeToTotal bool `json:"add_network_fee_to_total"`
	// Amount to withdraw
	Amount float64 `json:"amount,string"`
	// Crypto address of the recipient
	CryptoAddress string `json:"crypto_address"`
	// Currency to withdraw
	Currency CurrencyName `json:"currency"`
	// DestinationTag for currencies that support destination tagging
	DestinationTag string `json:"destination_tag"`
	// NoDestinationTag opts out of using a destination tag: required when not providing a destination tag
	NoDestinationTag bool `json:"no_destination_tag"`
}

type CryptoWithdrawal struct {
	Amount   float64      `json:"amount"`
	Currency CurrencyName `json:"currency"`
	Fee      float64      `json:"fee"`
	ID       string       `json:"id"`
	Subtotal float64      `json:"subtotal"`
}

type CryptoAddress struct {
	Currency      CurrencyName `json:"currency"`
	CryptoAddress string       `json:"crypto_address"`
}

func (c CryptoAddress) Params() []string {
	var params []string
	if c.Currency != "" {
		params = append(params, fmt.Sprintf("currency=%s", c.Currency))
	}
	if c.CryptoAddress != "" {
		params = append(params, fmt.Sprintf("crypto_address=%s", c.CryptoAddress))
	}
	return params
}

type WithdrawalFeeEstimate struct {
	Fee float64 `json:"fee"`
}

// GetWithdrawals retrieves a paginated list of Withdrawals for the current Profile, in descending order by CreatedAt time.
func (c *CoinbasePro) GetWithdrawals(ctx context.Context, filter WithdrawalFilter, pagination PaginationParams) (Withdrawals, error) {
	params := append(filter.Params(), pagination.Params()...)
	var withdrawals Withdrawals
	path := fmt.Sprintf("/%s/%s", coinbaseproTransfers, qsx.Query(params))
	err := c.API.Get(ctx, path, &withdrawals)
	if err != nil {
		return Withdrawals{}, err
	}
	// Withdrawals are a flavor of Transfer and the coinbasepro API cannot filter by multiple types
	// TODO: this potentially screws up pagination
	transferWithdrawals := make([]*Withdrawal, 0, len(withdrawals.Withdrawals))
	for _, transfer := range withdrawals.Withdrawals {
		if transfer.Type == WithdrawalTypeInternal || transfer.Type == WithdrawalTypeWithdraw {
			transferWithdrawals = append(transferWithdrawals, transfer)
		}
	}
	if len(transferWithdrawals) == 0 {
		withdrawals.Page = &Pagination{}
	}
	withdrawals.Withdrawals = transferWithdrawals
	return withdrawals, nil
}

// GetWithdrawal retrieves the details of a single Withdrawal. The Withdrawal must belong to the current Profile.
func (c *CoinbasePro) GetWithdrawal(ctx context.Context, withdrawalID string) (Withdrawal, error) {
	var withdrawal Withdrawal
	path := fmt.Sprintf("/%s/%s", coinbaseproTransfers, withdrawalID)
	return withdrawal, c.API.Get(ctx, path, &withdrawal)
}

// CreatePaymentMethodWithdrawal creates a Withdrawal of funds to an external PaymentMethod. Use ListPaymentMethods to
// retrieve details of available PaymentMethods.
func (c *CoinbasePro) CreatePaymentMethodWithdrawal(ctx context.Context, paymentMethodWithdrawal PaymentMethodWithdrawalSpec) (Withdrawal, error) {
	result := struct {
		ID string `json:"id"`
	}{}
	path := fmt.Sprintf("/%s/", coinbaseproWithdrawalPaymentMethod)
	err := c.API.Post(ctx, path, paymentMethodWithdrawal, &result)
	if err != nil {
		return Withdrawal{}, err
	}
	// POST coinbasepro response is partial; retrieve full representation
	return c.GetWithdrawal(ctx, result.ID)
}

// CreateCoinbaseAccountWithdrawal creates a Withdrawal of funds to a CoinbaseAccount. Funds can be moved between
// CoinbaseAccounts and Coinbase Pro trading Accounts within daily limits. Moving funds between Coinbase and Coinbase Pro
// is instant and free. Use ListCoinbaseAccounts to retrieve available Coinbase accounts.
func (c *CoinbasePro) CreateCoinbaseAccountWithdrawal(ctx context.Context, coinbaseAccountWithdrawal CoinbaseAccountWithdrawalSpec) (Withdrawal, error) {
	result := struct {
		ID string `json:"id"`
	}{}
	path := fmt.Sprintf("/%s/", coinbaseproWithdrawalCoinbaseAcct)
	err := c.API.Post(ctx, path, coinbaseAccountWithdrawal, &result)
	if err != nil {
		return Withdrawal{}, err
	}
	// POST coinbasepro response is partial; retrieve full representation
	return c.GetWithdrawal(ctx, result.ID)
}

// CreateCryptoAddressWithdrawal creates a Withdrawal of funds to a crypto address.
func (c *CoinbasePro) CreateCryptoAddressWithdrawal(ctx context.Context, cryptoAddressWithdrawal CryptoAddressWithdrawalSpec) (Withdrawal, error) {
	result := struct {
		ID string `json:"id"`
	}{}
	path := fmt.Sprintf("/%s/", coinbaseproWithdrawalCrypto)
	err := c.API.Post(ctx, path, cryptoAddressWithdrawal, &result)
	if err != nil {
		return Withdrawal{}, err
	}
	// POST coinbasepro response is partial; retrieve full representation
	return c.GetWithdrawal(ctx, result.ID)
}

// GetWithdrawalFeeEstimate retrieves the estimated network fees that would apply when sending to the given address.
func (c *CoinbasePro) GetWithdrawalFeeEstimate(ctx context.Context, cryptoAddress CryptoAddress) (WithdrawalFeeEstimate, error) {
	var withdrawalFeeEstimate WithdrawalFeeEstimate
	path := fmt.Sprintf("/%s/%s", coinbaseproWithdrawalFeeEstimate, qsx.Query(cryptoAddress.Params()))
	return withdrawalFeeEstimate, c.API.Get(ctx, path, &withdrawalFeeEstimate)
}
