package golog

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	Level   Level
	Message []interface{}
	Trace   []string
	Time    time.Time
	Entry   *Entry
}

func (msg *Message) JSON() []byte {
	buf, _ := json.Marshal(msg.JSON())
	return buf
}

func (msg *Message) MAP() *map[string]interface{} {
	data := map[string]interface{}{}

	if l := len(msg.Entry.Data); l > 0 {
		for _, i := range msg.Entry.Data {
			msg.Message = append(msg.Message, fmt.Sprintf("%s=%s", i.Field, fmt.Sprint(i.Value)))
		}
	}

	{
		data["log_level"] = LevelText[msg.Level]
		data["log_datetime"] = msg.Time.Format("2006-01-02 15:04:05")

		if l := len(msg.Entry.Tags); l > 0 {
			data["log_tags"] = msg.Entry.Tags
		}

		if l := len(msg.Message); l > 0 {
			var arr []string
			for _, i := range msg.Message {
				arr = append(arr, fmt.Sprint(i))
			}
			data["log_message"] = arr
		}

		if l := len(msg.Trace); l > 0 {
			data["log_trace"] = msg.Trace
		}
	}

	return &data
}
