package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mysql_byroad/model"
	"net/url"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
)

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
		resp, err := sendClient.Get(timeout).PostForm(pushurl, body)
		if err != nil {
			return "fail", err
		}
		defer resp.Body.Close()
		retStat, err := ioutil.ReadAll(resp.Body)
		return string(retStat), err
	} else {
		body := bytes.NewBuffer(msg)
		resp, err := sendClient.Get(timeout).Post(task.Apiurl, "application/json", body)
		if err != nil {
			return "fail", err
		}
		defer resp.Body.Close()
		retStat, err := ioutil.ReadAll(resp.Body)
		return string(retStat), err
	}
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
	log.Errorf("log send error: %+v, reason: %s", evt, reason)
	msg, _ := json.Marshal(evt)
	tl := &model.TaskLog{
		TaskId:     evt.TaskID,
		Message:    string(msg),
		Reason:     reason,
		CreateTime: time.Now(),
	}
	_, err := tl.Insert()
	if err != nil {
		log.Errorf("senderror insert error: %s", err.Error())
	}
}
