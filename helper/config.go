package helper

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	DebugLevel       bool
	Host             string
	Port             string
	ReadTimeoutMS    time.Duration
	WriteTimeoutMS   time.Duration
	CacheTTLSec      time.Duration
	YamlRedirectPath string
}
type EnvConfig struct {
	Config
}

func LoadConfigWithDefaults(cfg Config) *EnvConfig {
	c := Config{}
	dl, err := strconv.ParseBool(os.Getenv("DEBUG_LEVEL"))
	c.DebugLevel = dl
	if err != nil {
		c.DebugLevel = cfg.DebugLevel
	}
	h := os.Getenv("HOST")
	c.Host = h
	if h == "" {
		c.Host = cfg.Host
	}
	p := os.Getenv("PORT")
	c.Port = p
	if p == "" {
		c.Port = cfg.Port
	}
	rt, err := strconv.ParseInt(os.Getenv("READ_TIMEOUT_MS"), 10, 64)
	c.ReadTimeoutMS = time.Duration(rt) * time.Millisecond
	if err != nil {
		c.ReadTimeoutMS = cfg.ReadTimeoutMS
	}
	wt, err := strconv.ParseInt(os.Getenv("WRITE_TIMEOUT_MS"), 10, 64)
	c.WriteTimeoutMS = time.Duration(wt) * time.Millisecond
	if err != nil {
		c.WriteTimeoutMS = cfg.WriteTimeoutMS
	}
	ttl, err := strconv.ParseInt(os.Getenv("CACHE_TTL_SEC"), 10, 64)
	c.CacheTTLSec = time.Duration(ttl) * time.Second
	if err != nil {
		c.CacheTTLSec = cfg.CacheTTLSec
	}
	fp := os.Getenv("YAML_REDIRECT_PATH")
	if fp == "" {
		c.YamlRedirectPath = cfg.YamlRedirectPath
	}
	return &EnvConfig{
		c,
	}

}
