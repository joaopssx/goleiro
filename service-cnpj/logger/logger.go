package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Level string

const (
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "erro"
)

type LogEntry struct {
	Timestamp      string `json:"timestamp"`
	Nivel          Level  `json:"nivel"`
	Servico        string `json:"servico"`
	Mensagem       string `json:"mensagem"`
	Metodo         string `json:"metodo,omitempty"`
	Path           string `json:"path,omitempty"`
	IP             string `json:"ip,omitempty"`
	StatusCode     int    `json:"status_code,omitempty"`
	DuracaoMs      int64  `json:"duracao_ms,omitempty"`
	Arquivo        string `json:"arquivo,omitempty"`
	Linha          int    `json:"linha,omitempty"`
	Erro           string `json:"erro,omitempty"`
	EstadoAnterior string `json:"estado_anterior,omitempty"`
	EstadoNovo     string `json:"estado_novo,omitempty"`
}

var currentLevel = LevelInfo
var serviceName = "unknown"

func Init(name string) {
	serviceName = name
	env := os.Getenv("LOG_LEVEL")
	if env == "warn" {
		currentLevel = LevelWarn
	} else if env == "erro" {
		currentLevel = LevelError
	}
}

func shouldLog(lvl Level) bool {
	if currentLevel == LevelInfo {
		return true
	}
	if currentLevel == LevelWarn && (lvl == LevelWarn || lvl == LevelError) {
		return true
	}
	if currentLevel == LevelError && lvl == LevelError {
		return true
	}
	return false
}

func logEvent(entry LogEntry) {
	if !shouldLog(entry.Nivel) {
		return
	}
	entry.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	entry.Servico = serviceName
	
	b, err := json.Marshal(entry)
	if err == nil {
		fmt.Println(string(b))
	}
}

func Info(msg string) {
	logEvent(LogEntry{Nivel: LevelInfo, Mensagem: msg})
}

func Warn(msg string) {
	logEvent(LogEntry{Nivel: LevelWarn, Mensagem: msg})
}

func Error(msg string, arquivo string, linha int, err error) {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	logEvent(LogEntry{
		Nivel:    LevelError,
		Mensagem: msg,
		Arquivo:  arquivo,
		Linha:    linha,
		Erro:     errStr,
	})
}

func Request(msg, metodo, path, ip string, status int, duracao int64) {
	logEvent(LogEntry{
		Nivel:      LevelInfo,
		Mensagem:   msg,
		Metodo:     metodo,
		Path:       path,
		IP:         ip,
		StatusCode: status,
		DuracaoMs:  duracao,
	})
}

func StateChange(msg, estadoAnterior, estadoNovo string) {
	logEvent(LogEntry{
		Nivel:          LevelInfo,
		Mensagem:       msg,
		EstadoAnterior: estadoAnterior,
		EstadoNovo:     estadoNovo,
	})
}
