package transactionx

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type TransactionHandler interface {
	InsertTransactions(ctx context.Context, session sqlx.Session, info string) error
	DeleteTransactions(ctx context.Context, session sqlx.Session, info string) error
	UpdateTransactions(ctx context.Context, session sqlx.Session, info string) error
}

var (
	transactionManager = TransactionManager{
		TransactionHandlers: map[string]TransactionHandler{},
	}
)

type TransactionManager struct {
	TransactionHandlers map[string]TransactionHandler
}

func RegisterTransactionHandler(key string, handleFunc TransactionHandler) {
	transactionManager.TransactionHandlers[key] = handleFunc
}

func GetTransactionHandler(key string) TransactionHandler {
	key = fmt.Sprintf("`%s`", key)
	if handler, ok := transactionManager.TransactionHandlers[key]; ok {
		return handler
	}
	return nil
}
