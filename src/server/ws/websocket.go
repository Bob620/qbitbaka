package ws

import (
	"encoding/json"
	"net/http"

	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/rpc"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var rpcClient = rpc.CreateBakaRpc(nil, nil)

type WS struct {
	client *rpc.BakaRpc
}

func CreateWs() WS {
	ws := WS{rpcClient}

	rpcClient.RegisterMethod(
		"idk",
		[]parameters.Param{
			&parameters.StringParam{Name: "test", IsRequired: true},
		}, func(params map[string]parameters.Param) (returnMessage json.RawMessage, err error) {
			test, _ := params["test"].(*parameters.StringParam).GetString()

			return json.Marshal(test)
		})

	return ws
}

func (ws *WS) Handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()
	ws.client.UseChannels(rpc.MakeSocketReaderChan(conn), rpc.MakeSocketWriterChan(conn))
}
