package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var timeFormat = time.RFC3339

type ILogger interface {
	Date() time.Time
	Timezone() *time.Location
	Error(ctx context.Context, variables ...interface{})
	Log(ctx context.Context, variables ...interface{})
	Fatal(variables ...interface{})
	Critical(ctx context.Context, variables ...interface{})
}

type logType string

const (
	iError    logType = "ERROR"
	iCritical logType = "CRITICAL"
	iLog      logType = "LOG"
	iFatal    logType = "FATAL"
)

type discord struct {
	RB     *RequestBody `json:"rb"`
	Status logType      `json:"status"`
	Time   time.Time    `json:"time"`
	Info   string       `json:"info"`
}

type Logger struct {
	TZ      *time.Location
	Client  http.Client
	Webhook string
}

func ProdLogger(timezone, webhook string) (ILogger, error) {
	tz, err := Timezone(timezone)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return &Logger{
		TZ:      tz,
		Client:  http.Client{Timeout: 2 * time.Second},
		Webhook: webhook,
	}, nil
}

func (l *Logger) Timezone() *time.Location {
	return l.TZ
}

func (l *Logger) Date() time.Time {
	dt, err := time.Parse(timeFormat, time.Now().In(l.TZ).Format(timeFormat))
	if err != nil {
		log.Printf(err.Error())
		return time.Time{}
	}
	return dt
}

func (l *Logger) post(d *discord) {
	var title strings.Builder
	title.WriteString("📄 New Log Entry")
	if d.Status == iCritical || d.Status == iError {
		title.WriteString(" @everyone")
	}

	payload := map[string]interface{}{
		"embeds": []map[string]interface{}{
			{
				"title":       title.String(),
				"description": fmt.Sprintf("Status %s", d.Status),
				"color":       5814783, // color
				"fields": []map[string]string{
					{"name": "Request Id", "value": d.RB.Id, "inline": "true"},
					{"name": "IP", "value": d.RB.Ip, "inline": "true"},
					{"name": "Method", "value": d.RB.Method, "inline": "true"},
					{"name": "Path", "value": d.RB.Path, "inline": "false"},
					{"name": "Time", "value": d.Time.Format(time.RFC822), "inline": "false"},
					{"name": "Info", "value": d.Info, "inline": "false"},
				},
			},
		},
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		log.Printf("%s %s", iCritical, err.Error())
		return
	}

	if _, err := http.Post(l.Webhook, "application/json", buf); err != nil {
		log.Printf("%s %s", iCritical, err.Error())
	}
}

func (l *Logger) Error(ctx context.Context, variables ...interface{}) {
	d, str, err := logformat(ctx, iError, l.Date(), variables)
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(str)
	l.post(d)
}

func (l *Logger) Critical(ctx context.Context, variables ...interface{}) {
	d, str, err := logformat(ctx, iCritical, l.Date(), variables)
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(str)
	l.post(d)
}

func (l *Logger) Log(ctx context.Context, variables ...interface{}) {
	d, str, err := logformat(ctx, iLog, l.Date(), variables)
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(str)
	l.post(d)
}

func (l *Logger) Fatal(variables ...interface{}) {
	_, str, err := logformat(context.Background(), iFatal, l.Date(), variables)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(str)
}

func logformat(ctx context.Context, status logType, d time.Time, variables ...interface{}) (*discord, string, error) {
	obj, ok := ctx.Value(RequestKey).(*RequestBody)
	if !ok || obj == nil {
		return nil, "", errors.New(string(RequestKey) + "not found")
	}

	var sb strings.Builder
	for _, v := range variables {
		sb.WriteString(fmt.Sprintf("%v", v))
	}

	o := discord{
		RB:     obj,
		Status: status,
		Time:   d,
		Info:   sb.String(),
	}

	js, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return nil, "", err
	}

	return &o, string(js), nil
}
