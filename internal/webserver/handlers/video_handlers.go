package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/raulsilva-tech/DVRStreamAdapter/configs"
)

type VideoHandler struct{}

type Error struct {
	Message string `json:"message"`
}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{}
}

func (vh *VideoHandler) Stream(w http.ResponseWriter, r *http.Request) {

	//host, port, channel, start_time,end_time
	host := chi.URLParam(r, "host")
	port := chi.URLParam(r, "port")
	channel := chi.URLParam(r, "channel")
	startTime := chi.URLParam(r, "start_time")
	endTime := chi.URLParam(r, "end_time")
	if host == "" || port == "" || channel == "" || startTime == "" || endTime == "" {
		w.WriteHeader(http.StatusBadRequest)
		if host == "" {
			json.NewEncoder(w).Encode(Error{Message: "host is required"})
			return
		}
		if port == "" {
			json.NewEncoder(w).Encode(Error{Message: "port is required"})
			return
		}
		if channel == "" {
			json.NewEncoder(w).Encode(Error{Message: "channel is required"})
			return
		}
		if startTime == "" {
			json.NewEncoder(w).Encode(Error{Message: "startTime is required"})
			return
		}
		if endTime == "" {
			json.NewEncoder(w).Encode(Error{Message: "endTime is required"})
			return
		}
	}

	cfg, _ := configs.LoadConfig(".")

	// url := "http://dankia:dankia77@192.168.7.192:9090/cgi-bin/loadfile.cgi?action=startLoad&channel=1&startTime=2023-12-24%2021:44:29&endTime=2023-12-24%2021:45:0&Types=mp4"
	url := "http://" + cfg.DVRUser + ":" + cfg.DVRPassword + "@" + host + ":" + port + "/cgi-bin/loadfile.cgi?action=startLoad&channel=" + channel + "&startTime=" + startTime + "&endTime=" + endTime + "&Types=mp4"
	url = strings.Replace(url, " ", "%20", -1)

	//fmt.Println(url)

	dvrBody, err := onlyRequest(url)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Error{
			Message: err.Error(),
		})
		return
	}

	w.Header().Set("Content-Length", fmt.Sprint(len(dvrBody)))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", "video/mp4")

	src := bytes.NewBuffer(dvrBody)
	io.Copy(w, src)

}

func onlyRequest(url string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	dvrBody, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	return dvrBody, err
}
