package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type QBitUpdate struct {
	Categories  map[string]QBitCategory     `json:"categories"`
	ServerState QBitServerState             `json:"server_state"`
	Torrents    map[string]QBitTorrentState `json:"torrents"`
}

type QBitCategory struct {
	Name     string `json:"name"`
	SavePath string `json:"savePath"`
}

type QBitServerState struct {
	AllTimeDownload      int    `json:"alltime_dl"`
	AllTimeUpload        int    `json:"alltime_ul"`
	AverageTimeQueue     int    `json:"average_time_queue"`
	ConnectionStatus     string `json:"connection_status"`
	DhtNodes             int    `json:"dht_nodes"`
	DownloadInfoData     int    `json:"dl_info_data"`
	DownloadInfoSpeed    int    `json:"dl_info_speed"`
	DownloadRateLimit    int    `json:"dl_rate_limit"`
	FreeDiskSpace        int    `json:"free_space_on_disk"`
	GlobalRatio          string `json:"global_ratio"`
	QueuedIoJobs         int    `json:"queued_io_jobs"`
	Queueing             bool   `json:"queueing"`
	ReadCacheHits        string `json:"read_cache_hits"`
	ReadCacheOverload    string `json:"read_cache_overload"`
	RefreshInterval      int    `json:"refresh_interval"`
	TotalBufferSize      int    `json:"total_buffers_size"`
	TotalPeerConnections int    `json:"total_peer_connections"`
	TotalQueuedSize      int    `json:"total_queued_size"`
	TotalWastedSession   int    `json:"total_wasted_session"`
	UpInfoData           int    `json:"up_info_data"`
	UpInfoSpeed          int    `json:"up_info_speed"`
	UpRateLimit          int    `json:"up_rate_limit"`
	UsingAltSpeedLimits  bool   `json:"use_alt_speed_limits"`
	WriteCacheOverload   string `json:"write_cache_overload"`
}

type QBitTorrentState struct {
	AddedOn                int     `json:"added_on"`
	AmountLeft             int     `json:"amount_left"`
	AutoTMM                bool    `json:"auto_tmm"`
	Category               string  `json:"category"`
	Completed              int     `json:"completed"`
	CompletedOn            int     `json:"completed_on"`
	DownloadLimit          int     `json:"dl_limit"`
	DownloadSpeed          int     `json:"dlspeed"`
	Downloaded             int     `json:"downloaded"`
	DownloadedSession      int     `json:"downloaded_session"`
	ETA                    int     `json:"eta"`
	FirstLastPiecePriority bool    `json:"f_l_piece_prio"`
	ForceStart             bool    `json:"force_start"`
	LastActivity           int     `json:"last_activity"`
	MagnetURL              string  `json:"magnet_url"`
	MaxRatio               int     `json:"max_ratio"`
	MaxSeedingTime         int     `json:"max_seeding_time"`
	Name                   string  `json:"name"`
	NumberComplete         int     `json:"num_complete"`
	NumberIncomplete       int     `json:"num_incomplete"`
	NumberLeeches          int     `json:"num_leechs"`
	NumberSeeders          int     `json:"num_seeds"`
	Priority               int     `json:"priority"`
	Progress               float64 `json:"progress"`
	Ratio                  float64 `json:"ratio"`
	RatioLimit             int     `json:"ratio_limit"`
	SavePath               string  `json:"save_path"`
	SeedingTimeLimit       int     `json:"seeding_time_limit"`
	LastSeenComplete       int     `json:"seen_complete"`
	SequentialDownload     bool    `json:"seq_download"`
	Size                   int     `json:"size"`
	State                  string  `json:"state"`
	SuperSeeding           bool    `json:"super_seeding"`
	Tags                   string  `json:"tags"`
	TimeActive             int     `json:"time_active"`
	TotalSize              int     `json:"total_size"`
	Tracker                string  `json:"tracker"`
	UploadLimit            int     `json:"up_limit"`
	Uploaded               int     `json:"uploaded"`
	UploadSession          int     `json:"upload_session"`
	UploadSpeed            int     `json:"upspeed"`
}

func NewQBit() *QBit {
	updateDuration := time.Second * 2

	var qBit QBit
	qBit.timer = time.NewTimer(updateDuration)

	go func() {
		<-qBit.timer.C
		qBit.updateData()
		qBit.timer.Reset(updateDuration)
	}()
	qBit.updateSid()
	qBit.updateData()

	return &qBit
}

type QBit struct {
	sid        string
	timer      *time.Timer
	remoteData *QBitUpdate
}

func (qBit *QBit) GetData() QBitUpdate {
	return *qBit.remoteData
}

func (qBit *QBit) GetCategory(name string) QBitCategory {
	return qBit.remoteData.Categories[name]
}

func (qBit *QBit) GetTorrent(hash string) QBitTorrentState {
	return qBit.remoteData.Torrents[hash]
}

func (qBit *QBit) StopUpdate() {
	qBit.timer.Stop()
}

func (qBit *QBit) updateSid() bool {
	res, err := http.Get("http://localhost:8080/api/v2/auth/login")
	if err != nil {
		fmt.Println(err)
		return false
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil && fmt.Sprintf("%s", data) != "Ok." {
		fmt.Println(err)
		return false
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "SID" {
			qBit.sid = cookie.Value
		}
	}

	return true
}

func (qBit *QBit) updateData() {
	if len(qBit.sid) == 0 {
		qBit.updateSid()
	}

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v2/sync/maindata", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	req.AddCookie(&http.Cookie{Name: "SID", Value: qBit.sid})
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Host", "localhost")
	req.Header.Add("Accept-Encoding", "gzip")

	client := &http.Client{}
	res, err := client.Do(req)

	if res != nil {
		for _, cookie := range res.Cookies() {
			if cookie.Name == "SID" {
				qBit.sid = cookie.Value
			}
		}
	}

	if err != nil || res.ContentLength <= 0 {
		fmt.Println(err)
		return
	}

	reader, err := gzip.NewReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := ioutil.ReadAll(reader)
	_ = reader.Close()
	_ = res.Body.Close()

	var qBitData QBitUpdate
	err = json.Unmarshal(data, &qBitData)
	if err != nil {
		fmt.Println(err)
		return
	}

	qBit.remoteData = &qBitData
}
