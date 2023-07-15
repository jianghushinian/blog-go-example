package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
)

type Message struct {
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
	MsgType string `json:"msg_type"`
}

type Result struct {
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
	Code          int    `json:"code"`
	Data          any    `json:"data"`
	Msg           string `json:"msg"`
}

func sendFeishu(content, webhook string) (*Result, error) {
	msg := Message{
		Content: struct {
			Text string `json:"text"`
		}{
			Text: content,
		},
		MsgType: "text",
	}

	body, _ := json.Marshal(msg)
	resp, err := http.Post(webhook, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	result := new(Result)
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, fmt.Errorf("code: %d, error: %s", result.Code, result.Msg)
	}

	return result, nil
}

func monitor(pid int) (*Result, error) {
	for {
		// 检查进程是否存在
		err := syscall.Kill(pid, 0)
		if err != nil {
			log.Printf("Process %d exited\n", pid)
			webhook := os.Getenv("WEBHOOK")
			return sendFeishu(fmt.Sprintf("Process %d exited", pid), webhook)
		}

		log.Printf("Process %d is running\n", pid)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Println("Usage: ./monitor <pid>")
		return
	}

	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Printf("Invalid pid: %s\n", os.Args[1])
		return
	}

	result, err := monitor(pid)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(result)
}
