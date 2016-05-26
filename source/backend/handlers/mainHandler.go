package handlers

import (
	"github.com/gorilla/websocket"
	"encoding/json"
)

type Message struct {
	Action string `json:"action"`
	Data   interface{} `json:"data"`
}

type websocketConnection interface {
	Log(message string)
	Write(data interface{})
}

func Write(conn websocketConnection, action string, data interface{}) {
	message := Message{
		Action:action,
		Data:data,
	}
	conn.Write(message)
}

func HandleMessage(conn websocketConnection, messageType int, message []byte) {
	if messageType == websocket.TextMessage {
		var m Message
		json.Unmarshal(message, &m)
		switch m.Action {
		case "find-additional-folders": {
			findAdditionalFolders(conn, m.Data)
		}
		case "save-modpacks": {
			saveModpacks(conn, m.Data)
		}
		case "load-modpacks": {
			loadModpacks(conn)
		}
		case "gather-information": {
			gatherInformation(conn, m.Data)
		}
		case "build": {
			build(conn, m.Data)
		}
		}
	}
}

