package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"golang.org/x/net/context"
)

type ipInfo struct {
	IP       string `json:"ip"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
}

func main() {
	cfAPIKey := flag.String("CF_API_KEY", os.Getenv("CF_API_KEY"), "CF API Key")
	cfAPIEmail := flag.String("CF_API_EMAIL", os.Getenv("CF_API_EMAIL"), "CF API Email")
	zoneName := flag.String("CF_ZONE_NAME", os.Getenv("CF_ZONE_NAME"), "CF Zone Name")
	dnsName := flag.String("CF_DNS_NAME", os.Getenv("CF_DNS_NAME"), "CF DNS Name")
	defaultTTL, err := strconv.Atoi(os.Getenv("CF_DNS_TTL"))
	if err != nil {
		defaultTTL = 0
	}
	ttl := flag.Int("CF_DNS_TTL", defaultTTL, "CF NDS TTL")
	defaultIPCheckDuration, err := time.ParseDuration(os.Getenv("IP_CHECK_DURATION"))
	if err != nil {
		defaultIPCheckDuration = time.Minute
	}
	ipCheckDuration := flag.Duration("IP_CHECK_DURATION", defaultIPCheckDuration, "IP Check Duration")
	defaultTimeout, err := time.ParseDuration(os.Getenv("CF_TIMEOUT"))
	if err != nil {
		defaultTimeout = 30 * time.Second
	}
	timeout := flag.Duration("CF_TIMEOUT", defaultTimeout, "IP Check Duration")

	flag.Parse()
	for {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), *timeout)
			defer cancel()
			if err = DDNS(ctx, *cfAPIKey, *cfAPIEmail, *zoneName, *dnsName, *ttl); err != nil {
				log.Printf("[FAIL] %s\n", err.Error())
			}
		}()
		time.Sleep(*ipCheckDuration)
	}
}

func DDNS(ctx context.Context, cfAPIKey, cfAPIEmail, zoneName, dnsName string, ttl int) (err error) {
	ip, err := getIP(ctx)
	if err != nil {
		return fmt.Errorf("getIP fail, err: %w", err)
	}
	name := fmt.Sprintf("%s.%s", dnsName, zoneName)
	api, err := cloudflare.New(cfAPIKey, cfAPIEmail)
	if err != nil {
		return fmt.Errorf("new Cloudflare API fail, err: %w", err)
	}

	// get zone
	id, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return fmt.Errorf("get zoneID fail, err: %w", err)
	}
	zoneID := cloudflare.ZoneIdentifier(id)

	// get dns records
	dnsRecords, _, err := api.ListDNSRecords(ctx, zoneID, cloudflare.ListDNSRecordsParams{Name: name})
	if err != nil {
		return fmt.Errorf("list DNS records fail, err: %w", err)
	}

	proxy := false
	var dnsRecord cloudflare.DNSRecord
	if len(dnsRecords) == 0 {
		dnsRecord, err = api.CreateDNSRecord(ctx, zoneID, cloudflare.CreateDNSRecordParams{
			Type:    "A",
			Name:    name,
			Content: ip,
			TTL:     ttl,
			Proxied: &proxy,
		})
		if err != nil {
			return fmt.Errorf("create DNS record fail, err: %w", err)
		}
		log.Printf("DNS Record Created: %+v\n", dnsRecord)
	} else {
		dnsRecord = dnsRecords[0]
	}

	// same dns record, do nothing
	if dnsRecord.Content == ip {
		log.Printf("DNS Record Not Change (%s), Skip\n", ip)
		return nil
	}

	// update dns record
	_, err = api.UpdateDNSRecord(ctx, zoneID, cloudflare.UpdateDNSRecordParams{
		ID:      dnsRecord.ID,
		Type:    "A",
		Name:    name,
		Content: ip,
		TTL:     ttl,
		Proxied: &proxy,
	})
	if err != nil {
		return fmt.Errorf("update DNS record fail, err: %w", err)
	}
	log.Printf("update DNS record success, %s : %s -> %s\n", name, dnsRecord.Content, ip)
	return nil
}

func getIP(ctx context.Context) (string, error) {
	url := "https://ipinfo.io/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	info := ipInfo{}
	if err = json.Unmarshal(body, &info); err != nil {
		return "", err
	}
	if ip := net.ParseIP(info.IP); ip == nil {
		return "", fmt.Errorf("invalid ip, IPInfo: %+v", info)
	}
	log.Printf("get IPInfo success: %+v\n", info)
	return info.IP, nil
}
