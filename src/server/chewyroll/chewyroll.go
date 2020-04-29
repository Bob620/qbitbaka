package chewyroll

import (
	"fmt"
	"log"
	"net/url"

	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/rpc"
	"github.com/gorilla/websocket"
)

type Chewyroll struct {
	client *rpc.BakaRpc
}

func MakeChewyroll() *Chewyroll {
	u := url.URL{Scheme: "ws", Host: "localhost:7000", Path: "/"}
	fmt.Println("Chewyroll:", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	rpcClient := rpc.CreateBakaRpc(rpc.MakeSocketReaderChan(conn), rpc.MakeSocketWriterChan(conn))
	cr := Chewyroll{rpcClient}
	cr.authenticate("")

	return &cr
}

func (cr *Chewyroll) authenticate(code string) {
	_, err := cr.client.CallMethod(nil, "authenticate", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: code},
	}))
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (cr *Chewyroll) SeriesLookup(url string) []byte {
	fmt.Println("Looking up", url)
	res, err := cr.client.CallMethod(nil, "series.lookup", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: url},
	}))
	if err != nil {
		fmt.Println(err)
		return []byte("{}")
	}
	resJson, _ := res.MarshalJSON()

	fmt.Println(resJson)
	return resJson
}
