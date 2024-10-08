package goserver

import (
	"encoding/json"
	"errors"
	golog "github.com/gif-gif/go.io/go-log"
	goutils "github.com/gif-gif/go.io/go-utils"
	"time"
)

type Token struct {
	AppId     string `json:"app_id"`
	OpenId    int64  `json:"data"`
	NonceStr  string `json:"nonce"`
	Timestamp int64  `json:"ts"`
}

func (t *Token) Bytes() []byte {
	buf, _ := json.Marshal(t)
	return buf
}

func (t *Token) String() string {
	return string(t.Bytes())
}

func CreateToken(appId string, openid int64) (tokenStr string, err error) {
	token := &Token{
		AppId:     appId,
		OpenId:    openid,
		NonceStr:  goutils.NonceStr(),
		Timestamp: time.Now().Unix(),
	}

	var (
		key    = goutils.MD5([]byte(appId))
		iv     = key[8:24]
		encBuf []byte
	)

	encBuf, err = goutils.AESCBCEncrypt(token.Bytes(), []byte(key), []byte(iv))
	if err != nil {
		golog.Error(err.Error())
		return
	}

	tokenStr = goutils.Base64Encode(encBuf)
	return
}

func ParseToken(tokenStr, appId string) (token *Token, err error) {
	var (
		tokenBuf = goutils.Base64Decode(tokenStr)
		key      = goutils.MD5([]byte(appId))
		iv       = key[8:24]
		b        []byte
	)

	b, err = goutils.AESCBCDecrypt(tokenBuf, []byte(key), []byte(iv))
	if err != nil {
		golog.Error(err.Error())
		return
	}

	token = new(Token)
	if err = json.Unmarshal(b, token); err != nil {
		golog.Error(err.Error())
		return
	}
	if token.AppId != appId {
		err = errors.New("appid invalid")
	}
	return
}
