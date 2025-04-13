package main

import (
	"context"
	"log"
	"time"
)

type mockLogger struct {
	location *time.Location
}

func DevLogger(timezone string) ILogger {
	loc, err := Timezone(timezone)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return &mockLogger{location: loc}
}

func (m *mockLogger) Timezone() *time.Location {
	return m.location
}

func (m *mockLogger) Date() time.Time {
	dt, err := time.Parse(timeFormat, time.Now().In(m.location).Format(timeFormat))
	if err != nil {
		log.Print(err.Error())
	}
	return dt
}

func (m *mockLogger) Error(ctx context.Context, variables ...interface{}) {
	_, str, err := logformat(ctx, iError, m.Date(), variables)
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(str)
}

func (m *mockLogger) Log(ctx context.Context, variables ...interface{}) {
	_, str, err := logformat(ctx, iLog, m.Date(), variables)
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(str)
}

func (m *mockLogger) Fatal(variables ...interface{}) {
	_, str, err := logformat(context.Background(), iFatal, m.Date(), variables)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(str)
}

func (m *mockLogger) Critical(ctx context.Context, variables ...interface{}) {
	_, str, err := logformat(ctx, iCritical, m.Date(), variables)
	if err != nil {
		log.Print(err.Error())
		return
	}
	log.Print(str)
}
