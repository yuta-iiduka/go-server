package main

import (
	// 標準
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	// 自製
	"example.com/http/src/service/sample"
)

type Config struct {
	Host        string `json:"Host"`
	Port        int    `json:"Port"`
	CertFile    string `json:"CertFile"`
	KeyFile     string `json:"KeyFile"`
	TemplateDir string `json:"Template"`
}

type Message struct {
}

type Constant struct {
}

func loadConfig(filename string) (Config, error) {
	var c Config

	data, err := os.ReadFile(filename)
	if err != nil {
		return c, err
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return c, err
	}

	if c.Host == "" {
		c.Host = "127.0.0.1"
	}
	if c.Port == 0 {
		c.Port = 8443
	}
	if c.Port < 1 || c.Port > 65535 {
		return c, os.ErrInvalid
	}
	if c.CertFile == "" || c.KeyFile == "" ||
		c.TemplateDir == "" {
		return c, os.ErrInvalid
	}

	return c, nil
}

func main() {

	sample.Hello()

	// カレントディレクトリ
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 設定ファイルを読み込む
	config, err := loadConfig(filepath.Join(workDir, "config", "settings.json"))

	// サーバオブジェクト生成
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, HTTPS!\n")
		fmt.Fprintf(w, "Method: %s\n", r.Method)
		fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
	})

	server := &http.Server{
		Addr:              net.JoinHostPort(config.Host, strconv.Itoa(config.Port)),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	log.Println("HTTPS server listening on :8443")

	// 証明書ファイル
	if err := server.ListenAndServeTLS(config.CertFile, config.KeyFile); err != nil {
		log.Fatal(err)
	}
}
