package chewyroll

import (
	"encoding/json"
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
		fmt.Println("CR Auth Error:", err)
	}
}

func (cr *Chewyroll) SeriesLookupById(id string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "series.lookup.id", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: id},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) SeriesLookupByUuid(uuid string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "series.lookup.uuid", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: uuid},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) SeriesLookupByUrl(url string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "series.lookup.url", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: url},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) SeriesSearchById(id string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "series.search.id", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: id},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) SeriesSearchByName(name string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "series.search.name", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: name},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) SeriesUpdate(uuid string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "series.update", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: uuid},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) SeriesDownload(uuid string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "series.queueDownload", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: uuid},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) EpisodesLookup(uuid string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "episodes.lookup", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: uuid},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) EpisodesUpdate(uuid string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "episodes.update", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: uuid},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) EpisodesDownload(uuid string) *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "episodes.download", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: uuid},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}
