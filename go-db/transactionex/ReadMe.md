# 事务性操作 主要在go-zero 框架下执行
```go
初始化 ： TransactionModel: transactionex.NewModel(conn, cacheConf),

userInserter := func(session sqlx.Session) error {
    return l.svcCtx.UserModel.InsertWithTx(l.ctx, session, newUserInfo)
}
deviceInserter := func(session sqlx.Session) error {
    return l.svcCtx.DeviceModel.InsertWithTx(l.ctx, session, deviceInfo)
}
inserters := []transactionex.TableTransactionFunc{
    userInserter,
    deviceInserter,
}
err := l.svcCtx.TransactionModel.Transactions(inserters)
if err != nil {
    return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Register InsertWithTx error:%v in:%+v", err, in)
}

```