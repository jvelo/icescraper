// Copyright 2021 Jérôme Velociter
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:generate go run github.com/prisma/prisma-client-go generate
//go:generate go run github.com/99designs/gqlgen generate

package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/jvelo/icecast-monitor/config"
	"github.com/jvelo/icecast-monitor/model"
	"github.com/jvelo/icecast-monitor/updater"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/joho/godotenv/autoload"
	"github.com/jvelo/icecast-monitor/prisma/db"
)

type options struct {
	Config string `short:"c" long:"config" description:"Path to configuration file" default:"./config.yml"`
}

type Source struct {
	AudioInfo     string `json:"audio_info"`
	Bitrate       int    `json:"bitrate"`
	Genre         string `json:"genre"`
	ListenersPeak int    `json:"listeners_peak"`
	Listeners     int    `json:"listeners"`
	ListenURL     string `json:"listenurl"`
	Description   string `json:"server_description"`
	Name          string `json:"server_name"`
	Type          string `json:"server_type"`
	Url           string `json:"server_url"`
	StreamStart   string `json:"stream_start_iso8601"`
	Title         string `json:"title"`
}

type Icestats struct {
	Admin       string `json:"admin"`
	Host        string `json:"host"`
	Location    string `json:"location"`
	Id          string `json:"server_id"`
	ServerStart string `json:"server_start_iso8601"`
	Source      Source `json:"source"`
}

type Response struct {
	Stats Icestats `json:"icestats"`
}

var (
	opts         options
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	parser := flags.NewParser(&opts, flags.HelpFlag|flags.PrintErrors)
	_, err := parser.Parse()
	if err != nil {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	client := db.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		return err
	}

	defer func() {
		if err := client.Prisma.Disconnect(); err != nil {
			panic(err)
		}
	}()

	conf, err := config.LoadFile(opts.Config)
	if err != nil {
		return err
	}
	log.Infof("conf: %v", conf)

	ticker := time.NewTicker(conf.ScrapeInterval)
	defer ticker.Stop()

	insecureTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	secureTransport := &http.Transport{}
	c := http.Client{Transport: secureTransport}

	stream := make(chan *model.Record)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go updater.Update(ctx, stream)

	go func() {
		for {
			select {
			case t := <-ticker.C:
				log.Infof("tick %v", t)
				go func() {
					for _, target := range conf.Servers {
						log.Infof("Polling target: %v", target.Url)

						if target.SkipCertCheck {
							c.Transport = insecureTransport
						} else {
							c.Transport = secureTransport
						}

						body, err := doRequest(target, c)
						if err != nil {
							log.Errorf("polling target: %v", err)
							continue
						}
						var response Response
						err = json.Unmarshal(body, &response)
						if err != nil {
							log.Errorf("unmarshalling target: %v", err)
							continue
						}

						cast := model.NewCast(
							response.Stats.Source.Name,
							response.Stats.Source.Description,
							target.Url,
						)
						track := model.NewTrack(
							response.Stats.Source.Title,
							response.Stats.Source.Listeners,
						)

						go func() {
							stream <- &model.Record{
								Cast:  cast,
								Track: track,
							}
						}()
					}
				}()
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe("127.0.0.1:2112", nil)
}

func doRequest(target config.IcecastServer, c http.Client) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%v/status-json.xsl", target.Url), nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating request")
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("didn't get a OK status: %v", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading body")
	}
	return body, nil
}
