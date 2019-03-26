package dns

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Config struct {
	server    string
	port      int
	transport string
	timeout   time.Duration
	retries   int
	keyname   string
	keyalgo   string
	keysecret string
}

type DNSClient struct {
	c         *dns.Client
	srv_addr  string
	transport string
	retries   int
	keyname   string
	keysecret string
	keyalgo   string
}

// Configures and returns a fully initialized DNSClient
func (c *Config) Client() (interface{}, error) {
	log.Println("[INFO] Building DNSClient config structure")

	var client DNSClient
	client.srv_addr = net.JoinHostPort(c.server, strconv.Itoa(c.port))
	authCfgOk := false
	if (c.keyname == "" && c.keysecret == "" && c.keyalgo == "") ||
		(c.keyname != "" && c.keysecret != "" && c.keyalgo != "") {
		authCfgOk = true
	}
	if !authCfgOk {
		return nil, fmt.Errorf("Error configuring provider: when using authentication, \"key_name\", \"key_secret\" and \"key_algorithm\" should be non empty")
	}
	client.c = new(dns.Client)
	client.transport = c.transport
	client.c.Timeout = c.timeout
	client.retries = c.retries
	if c.keyname != "" {
		if !dns.IsFqdn(c.keyname) {
			return nil, fmt.Errorf("Error configuring provider: \"key_name\" should be fully-qualified")
		}
		keyname := strings.ToLower(c.keyname)
		client.keyname = keyname
		client.keysecret = c.keysecret
		keyalgo, err := convertHMACAlgorithm(c.keyalgo)
		if err != nil {
			return nil, fmt.Errorf("Error configuring provider: %s", err)
		}
		client.keyalgo = keyalgo
		client.c.TsigSecret = map[string]string{keyname: c.keysecret}
	}
	client.transport = c.transport
	client.c.Timeout = c.timeout
	client.retries = c.retries
	return &client, nil
}

// Validates and converts HMAC algorithm
func convertHMACAlgorithm(name string) (string, error) {
	switch name {
	case "hmac-md5":
		return dns.HmacMD5, nil
	case "hmac-sha1":
		return dns.HmacSHA1, nil
	case "hmac-sha256":
		return dns.HmacSHA256, nil
	case "hmac-sha512":
		return dns.HmacSHA512, nil
	default:
		return "", fmt.Errorf("Unknown HMAC algorithm: %s", name)
	}
}
