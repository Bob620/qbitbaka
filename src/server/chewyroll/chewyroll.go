package chewyroll

import (
	"encoding/json"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"net/url"
	"time"

	"github.com/bob620/baka-rpc-go/parameters"
	"github.com/bob620/baka-rpc-go/rpc"
	"github.com/gorilla/websocket"
)

type Chewyroll struct {
	url    url.URL
	conn   *websocket.Conn
	client *rpc.BakaRpc
}

func MakeChewyroll(host string) *Chewyroll {
	cr := Chewyroll{url.URL{Scheme: "ws", Host: host, Path: "/"}, nil, rpc.CreateBakaRpc(nil, nil)}
	cr.init()

	return &cr
}

func (cr *Chewyroll) init() {
	if cr.conn != nil {
		fmt.Println("Closing hanging Chewyroll connection...")
		_ = cr.conn.Close()
	}

	fmt.Println("Connecting to Chewyroll:", cr.url.String())

	conn, _, err := websocket.DefaultDialer.Dial(cr.url.String(), nil)
	if err != nil {
		fmt.Println("Unable to connect to Chewyroll, retrying in 5 seconds...")
		time.AfterFunc(5*time.Second, cr.init)
		return
	}

	cr.conn = conn
	cr.client.RemoveChannels(nil)
	cr.client.AddChannels(rpc.MakeSocketReaderChan(conn), rpc.MakeSocketWriterChan(conn))

	cr.client.HandleDisconnect(func(uuid *uuid.UUID) {
		fmt.Println("Lost connection to Chewyroll, restarting websocket in 5 seconds...")
		time.AfterFunc(5*time.Second, cr.init)
	})
	fmt.Println("Connected to Chewyroll, authenticating...")
	cr.authenticate("")
}

func (cr *Chewyroll) authenticate(code string) {
	_, err := cr.client.CallMethod(nil, "authenticate", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{Default: code},
	}))
	if err != nil {
		fmt.Println("CR Auth Error:", err)
	}
	fmt.Println("Chewyroll authenticated")
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

func (cr *Chewyroll) CRListingUpdate() *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "cr.listing.update", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) CRSeasonUpdate() *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "cr.season.update", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) CRSeasonGet() *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "cr.season.get", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) MalSeasonUpdate() *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "mal.season.update", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (cr *Chewyroll) MalSeasonGet() *json.RawMessage {
	res, err := cr.client.CallMethod(nil, "mal.season.get", parameters.NewParametersByPosition([]parameters.Param{
		&parameters.StringParam{},
	}))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}
