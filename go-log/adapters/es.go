package adapters

import (
	goes "github.com/gif-gif/go.io/go-db/go-es"
	"github.com/gif-gif/go.io/go-log"
	"sync"
)

type EsAdapter struct {
	index string
	id    string
	opt   goes.Config
	mu    sync.Mutex
}

func NewEsLog(index, id string, opt goes.Config) *golog.Logger {
	return golog.New(NewEsAdapter(index, id, opt))
}

func NewEsAdapter(index, id string, opt goes.Config) *EsAdapter {
	err := goes.Init(opt)
	if err != nil {
		return nil
	}
	fa := &EsAdapter{
		index: index,
		id:    id,
		opt:   opt,
	}
	return fa
}

func (fa *EsAdapter) Write(msg *golog.Message) {
	client := goes.GetClient(fa.opt.Name)
	if client == nil {
		return
	}
	_, err := goes.GetClient(fa.opt.Name).DocCreate(fa.index, fa.id, msg.MAP())
	if err != nil {
		return
	}
}

func (fa *EsAdapter) closeEs() {
	client := goes.GetClient(fa.opt.Name)
	if client == nil {
		return
	}
}
