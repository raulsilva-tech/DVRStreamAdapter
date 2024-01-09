package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type VideoHandler struct{}

type Error struct {
	Message string `json:"message"`
}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{}
}

func (vh *VideoHandler) Stream(w http.ResponseWriter, r *http.Request) {

	userId := chi.URLParam(r, "id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Error{
			Message: "id is required",
		})
		return
	}

	url := "http://dankia:dankia77@192.168.7.192:9090/cgi-bin/loadfile.cgi?action=startLoad&channel=1&startTime=2023-12-24%2021:44:29&endTime=2023-12-24%2021:45:0&Types=mp4"

	dvrBody, err := onlyRequest(url)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Error{
			Message: err.Error(),
		})
		return
	}

	reader := bytes.NewReader(dvrBody)

	w.Header().Set("Content-Length", fmt.Sprint(len(dvrBody)))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", "video/mp4")
	buffer := make([]byte, 64*1024) // 64KB buffer size
	io.CopyBuffer(w, reader, buffer)
	// io.Copy(w,file)

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
