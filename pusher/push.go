package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mysql_byroad/model"
	"net/http"
	"net/url"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

type SendClient struct {
	http.Client
}

var sendClient *SendClient

func NewSendClient() *SendClient {
	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1000,
		},
	}
	sendClient := &SendClient{
		Client: httpClient,
	}
	return sendClient
}

func init() {
	sendClient = NewSendClient()
}

/*
发送消息
*/
func (sc *SendClient) SendMessage(evt *model.NotifyEvent) (string, error) {
	task := taskManager.GetTask(evt.TaskID)
	if task == nil {
		return "success", nil
	}
	evt.LastSendTime = time.Now()
	msg, _ := json.Marshal(evt)
	timeout := time.Millisecond * time.Duration(task.Timeout)
	if task.PackProtocal == model.PackProtocalEventCenter {
		idStr := strconv.FormatInt(task.ID, 10)
		retryCountStr := strconv.Itoa(evt.RetryCount)
		pushurl := task.Apiurl + "?" + url.Values{"jobid": {idStr}, "retry_times": {retryCountStr}}.Encode()
		body := url.Values{"message": {string(msg)}}
		resp, err := sendClient.PostForm(pushurl, body)
		if err != nil {
			return "fail", err
		}
		defer resp.Body.Close()
		retStat, err := ioutil.ReadAll(resp.Body)
		return string(retStat), err
	} else {
		body := bytes.NewBuffer(msg)
		sendClient.Timeout = timeout
		resp, err := sendClient.Post(task.Apiurl, "application/json", body)
		if err != nil {
			return "fail", err
		}
		defer resp.Body.Close()
		retStat, err := ioutil.ReadAll(resp.Body)
		return string(retStat), err
	}
	return "success", nil
}

func isSuccessSend(msg string) bool {
	if msg == "success" {
		return true
	} else {
		type SendResp struct {
			Status int `json:"status"`
		}
		var sendResp SendResp
		if json.Unmarshal([]byte(msg), &sendResp) == nil {
			if sendResp.Status == 1 {
				return true
			}
		}
		return false
	}
}

func (sc *SendClient) ResendMessage(evt *model.NotifyEvent) {
	task := taskManager.GetTask(evt.TaskID)
	ticker := time.NewTicker(time.Duration(task.ReSendTime) * time.Millisecond)
	var err error
	var ret string
	for i := 0; i < task.RetryCount; i++ {
		<-ticker.C
		evt.RetryCount++
		ret, err = sc.SendMessage(evt)
		log.Debugf("resend message ret: %s, err: %v", ret, err)
		if isSuccessSend(ret) {
			return
		}
	}
	if err != nil {
		sc.LogSendError(evt, err.Error())
	} else {
		sc.LogSendError(evt, ret)
	}
}

func (sc *SendClient) LogSendError(evt *model.NotifyEvent, reason string) {
	log.Debugf("log send error: %+v, reason: %s", evt, reason)
}
