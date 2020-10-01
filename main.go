package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	cfAPIKey := flag.String("CF_API_KEY", os.Getenv("CF_API_KEY"), "CF API Key")
	cfAPIEmail := flag.String("CF_API_EMAIL", os.Getenv("CF_API_EMAIL"), "CF API Email")
	zoneName := flag.String("CF_ZONE_NAME", os.Getenv("CF_ZONE_NAME"), "CF Zone Name")
	dnsName := flag.String("CF_DNS_NAME", os.Getenv("CF_DNS_NAME"), "CF DNS Name")
	sTTL := os.Getenv("CF_DNS_TTL")
	iTTL := 0
	err := errors.New("OK")
	if sTTL != "" {
		iTTL, err = strconv.Atoi(sTTL)
		if err != nil {
			panic(err)
		}
	}
	ttl := flag.Int("CF_DNS_TTL", iTTL, "CF NDS TTL")
	flag.Parse()

	for true {
		err := DDNS(*cfAPIKey, *cfAPIEmail, *zoneName, *dnsName, *ttl)
		if err != nil {
			log.Println(fmt.Sprintf("[FAIL] %s", err.Error()))
		}
		time.Sleep(1 * time.Minute)
	}
}

func DDNS(cfAPIKey, cfAPIEmail, zoneName, dnsName string, ttl int) (err error) {
	// Construct a new API object
	ip, err := getIP()
	if err != nil {
		return err
	}
	api, err := cloudflare.New(cfAPIKey, cfAPIEmail)
	if err != nil {
		return err
	}

	// Fetch the zone ID
	id, err := api.ZoneIDByName(zoneName) // Assuming example.com exists in your Cloudflare account already
	if err != nil {
		return err
	}

	// Fetch zone details
	zone, err := api.ZoneDetails(id)
	if err != nil {
		return err
	}

	domain := strings.Join([]string{dnsName, zoneName}, ".")
	dnsRecords, err := api.DNSRecords(zone.ID, cloudflare.DNSRecord{
		Name: domain,
	})
	if err != nil {
		return err
	}

	if len(dnsRecords) == 0 {
		dnsRecordResponse, err := api.CreateDNSRecord(zone.ID, cloudflare.DNSRecord{
			Type:    "A",
			Name:    domain,
			Content: ip,
			TTL:     ttl,
		})
		if err != nil {
			return err
		}
		if !dnsRecordResponse.Success {
			return errors.New(fmt.Sprintf("Create DNS Record(%s) Fail, rsp(%v)", domain, dnsRecordResponse))
		}
		log.Println(fmt.Sprintf("DNS Record Created, %s : %s", domain, ip))
		return nil
	}
	if dnsRecords[0].Content == ip {
		log.Println(fmt.Sprintf("DNS Record Not Change (%s), Skip", ip))
		return nil
	}
	err = api.UpdateDNSRecord(zone.ID, dnsRecords[0].ID, cloudflare.DNSRecord{
		Type:    "A",
		Name:    domain,
		Content: ip,
		TTL:     ttl,
	})
	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("DNS Record Updated, %s : %s -> %s", domain, dnsRecords[0].Content, ip))
	return nil
}

func getIP() (ip string, err error) {
	url := "http://myip.ipip.net/"
	method := "GET"

	client := &http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	reg, err := regexp.Compile(`^当前 IP：(\d{1,3}(\.\d{1,3}){3})`)
	if err != nil {
		return "", err
	}
	find := reg.FindSubmatch(body)
	if len(find) < 2 {
		return "", errors.New(fmt.Sprintf("Fail To Get IP, Rsp(%s)", string(body)))
	}
	log.Println(fmt.Sprintf("Get IP Success, IP(%s)", string(find[1])))
	return string(find[1]), nil
}
