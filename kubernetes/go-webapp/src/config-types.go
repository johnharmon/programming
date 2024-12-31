package main

import "time"

//#region BASIC STRUCTS

type ServerConfig struct {
	Name    string   `yaml:"name"`
	Port    []int    `yaml:"port"`
	Address []string `yaml:"address"`
}

type TokenCookieConfig struct {
	ExpirationUnit string `yaml:"expiration_unit"`
	ExpirationMult int    `yaml:"expiration_mult"`
	ExpirationTime time.Time
}

func (c TokenCookieConfig) GetExpirationTime() (exp time.Time) {
	return time.Now().Add(Times[c.ExpirationUnit] * time.Duration(c.ExpirationMult))
}

func (c *TokenCookieConfig) SetExpirationTime() (exp time.Time) {
	c.ExpirationTime = time.Now().Add(Times[c.ExpirationUnit] * time.Duration(c.ExpirationMult))
	return c.ExpirationTime
}

type TokenExpirationConfig struct {
	ExpirationUnit string `yaml:"expiration_unit"`
	ExpirationMult int    `yaml:"expiration_mult"`
	ExpirationTime time.Time
}

func (e TokenExpirationConfig) GetExpirationTime() (exp time.Time) {
	return time.Now().Add(Times[e.ExpirationUnit] * time.Duration(e.ExpirationMult))
}

func (e *TokenExpirationConfig) SetExpirationTime() (exp time.Time) {
	e.ExpirationTime = time.Now().Add(Times[e.ExpirationUnit] * time.Duration(e.ExpirationMult))
	return e.ExpirationTime
}

// #endregion

// region COMPOSITE STRUCTS

type TokenConfig struct {
	Issuer     string                `yaml:"issuer"`
	Audience   string                `yaml:"audience"`
	Subject    string                `yaml:"subject"`
	Expiration TokenExpirationConfig `yaml:"expiration"`
}

//#endregion

// #region NESTED COMPOSITE STRUCTS

type JwtConfig struct {
	Token  TokenConfig       `yaml:"token"`
	Cookie TokenCookieConfig `yaml:"cookie"`
}

type Config struct {
	Server    ServerConfig `yaml:"server"`
	JWTConfig JwtConfig    `yaml:"jwt_config"`
}

//#endregion

type TimeMap map[string]time.Duration

var Times = TimeMap{
	"hour":   time.Hour,
	"second": time.Second,
	"minute": time.Minute,
}

var DefaultConfig = Config{
	Server: ServerConfig{
		Name:    "Webapp",
		Port:    []int{8080, 8443},
		Address: []string{"127.0.0.1"},
	},
	JWTConfig: JwtConfig{
		Token: TokenConfig{
			Issuer:   "webapp",
			Audience: "dev",
			Subject:  "testing",
			Expiration: TokenExpirationConfig{
				ExpirationUnit: "hour",
				ExpirationMult: 1,
			},
		},
		Cookie: TokenCookieConfig{
			ExpirationUnit: "hour",
			ExpirationMult: 1,
		},
	},
}
