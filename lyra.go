package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	influx "github.com/influxdata/influxdb-client-go"
)

type result struct {
	DNSLookup        float64
	TCPConnection    float64
	TLSHandshake     float64
	ServerProcessing float64
	ContentTransfer  float64
	Total            float64
	Availability     int
	cdn              string
	workflow         string
	contentType      string
	Headers          string
	Error            error
}

type Config struct {
	TagPrefix                  string `json:"tag_prefix"`
	NumUsers                   int    `json:"num_users"`
	Debug                      bool   `json:"debug"`
	UseInfluxDB                bool   `json:"use_influx_db"`
	InfluxDBUrl                string `json:"influx_db_uri"`
	InfluxDBApiKey             string `json:"influx_db_api_key"`
	InfluxDBOrg                string `json:"influx_db_org"`
	InfluxDBBucket             string `json:"influx_db_bucket"`
	InfluxPointMeasurementName string `json:"influx_point_measurement_name"`
	Resources                  []struct {
		URL      string `json:"url"`
		CDN      string `json:"cdn"`
		Workflow string `json:"workflow"`
	} `json:"resources"`
}

func writePoint(result result, influxUrl string, org string, bucket string, key string, pointname string, tagprefix string) error {
	client := influx.NewClient(influxUrl, key)
	writeAPI := client.WriteAPI(org, bucket)
	defer client.Close()

	// DNS Lookup   TCP Connection   TLS Handshake   Server Processing   Content Transfer Total
	point := influx.NewPointWithMeasurement(pointname).
		AddTag(tagprefix+"cdn", result.cdn).
		AddTag(tagprefix+"workflow", result.workflow).
		AddTag(tagprefix+"contenttype", result.contentType).
		AddField(tagprefix+"dnslookup", result.DNSLookup).
		AddField(tagprefix+"tcpconnection", result.TCPConnection).
		AddField(tagprefix+"tlshandshake", result.TLSHandshake).
		AddField(tagprefix+"serverprocessing", result.ServerProcessing).
		AddField(tagprefix+"contenttransfer", result.ContentTransfer).
		AddField(tagprefix+"total", result.Total).
		AddField(tagprefix+"availability", result.Availability).
		AddField(tagprefix+"headers", result.Headers).
		AddField(tagprefix+"error", result.Error).
		SetTime(time.Now().UTC())

	writeAPI.WritePoint(point)

	return nil
}

func parseConfig(filePath string) (Config, error) {
	var config Config

	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func getPort(url string) string {
	if strings.Contains(url, "https://") {
		return ":443"
	}
	return ":80"
}
func checkResource(url string, cdn string, workflow string) result {

	start := time.Now()

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   0,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result{Error: err}
	}
	req.Header.Add("User-Agent", "Akamai Lyra/1.2;  Performance Metrics Agent")

	dnsStart := time.Now()
	_, err = net.LookupHost(req.URL.Hostname())
	if err != nil {
		return result{Error: err}
	}
	dnsLookup := time.Since(dnsStart)

	tcpStart := time.Now()
	conn, err := net.Dial("tcp", req.URL.Host)
	if err != nil {
		return result{Error: err}
	}
	defer conn.Close()
	tcpConnection := time.Since(tcpStart)

	tlsStart := time.Now()
	tlsConn := tls.Client(conn, &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         req.URL.Hostname(),
	})

	err = tlsConn.Handshake()
	if err != nil {
		return result{Error: err}
	}
	defer tlsConn.Close()
	tlsHandshake := time.Since(tlsStart)

	serverStart := time.Now()
	err = req.Write(tlsConn)
	if err != nil {
		return result{Error: err}
	}
	_, err = http.ReadResponse(bufio.NewReader(tlsConn), req)
	if err != nil {
		return result{Error: err}
	}
	serverProcessing := time.Since(serverStart)

	contentStart := time.Now()
	res, err := client.Do(req)

	if err != nil {
		return result{Error: err}
	}

	availability := res.StatusCode
	contenttype := res.Header.Get("content-type")

	headerBytes, err := json.Marshal(res.Header)
	if err != nil {
		return result{Error: err}
	}
	headers := string(headerBytes)

	//_, err = io.Copy(io.Discard, res.Body) // avoud ERR_CLIENT_ABORT
	_, err = io.ReadAll(res.Body)
	if err != nil {
		return result{Error: err}
	}
	defer res.Body.Close()

	contentTransfer := time.Since(contentStart)

	return result{
		DNSLookup:        float64(dnsLookup.Milliseconds()),
		TCPConnection:    float64(tcpConnection.Milliseconds()),
		TLSHandshake:     float64(tlsHandshake.Milliseconds()),
		ServerProcessing: float64(serverProcessing.Milliseconds()),
		ContentTransfer:  float64(contentTransfer.Milliseconds()),
		Total:            float64(time.Since(start).Milliseconds()),
		Availability:     availability,
		cdn:              cdn,
		workflow:         workflow,
		contentType:      contenttype,
		Headers:          headers,
	}
}

func main() {
	config, err := parseConfig("config.json")
	// Accept user input for the number of virtual users
	numUsers := config.NumUsers
	debug := config.Debug
	resources := config.Resources
	useInfluxDB := config.UseInfluxDB
	influxDBUrl := config.InfluxDBUrl
	influxDBApiKey := config.InfluxDBApiKey
	influxDBOrg := config.InfluxDBOrg
	InfluxDBBucket := config.InfluxDBBucket
	InfluxPointMeasurementName := config.InfluxPointMeasurementName
	InfluxTagPrevix := config.TagPrefix

	if err != nil {
		log.Fatal(err)
	}

	results := make(chan result)

	for i := 0; i < numUsers; i++ {
		start := time.Now()
		for j := 0; j < len(resources); j++ {
			go func(url string, cdn string, workflow string) result {
				results <- checkResource(url, cdn, workflow)
				return result{}
			}(resources[j].URL, resources[j].CDN, resources[j].Workflow)
		}

		for j := 0; j < len(resources); j++ {
			r := <-results

			// DNS Lookup   TCP Connection   TLS Handshake   Server Processing   Content Transfer Total
			if debug {
				resultJSON, err := json.Marshal(r)
				if err != nil {
					fmt.Println("In Debug: Error marshaling resultJSON:", err)
					return
				}
				fmt.Println(string(resultJSON))
			}

			r.Total = float64(time.Since(start).Milliseconds())

			if useInfluxDB {
				err := writePoint(r, influxDBUrl, influxDBOrg, InfluxDBBucket, influxDBApiKey, InfluxPointMeasurementName, InfluxTagPrevix)

				if err != nil {
					fmt.Printf("Error writing point: %v\n", err)
				}
			}

		}
	}
}
