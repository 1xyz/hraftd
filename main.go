package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/1xyz/hraftd/http"
	"github.com/1xyz/hraftd/store"
)

// Command line defaults
const (
	DefaultHTTPAddr = ":11000"
	DefaultRaftAddr = "127.0.0.1:12000"
)

// Command line parameters
var (
	inmem           bool
	httpAddr        string
	raftAddr        string
	joinAddr        string
	nodeID          string
	bootstrapNodeID string
	dataDir         string
)

func init() {
	flag.BoolVar(&inmem, "inmem", false, "Use in-memory storage for Raft")
	flag.StringVar(&httpAddr, "haddr", DefaultHTTPAddr, "Set the HTTP bind address")
	flag.StringVar(&raftAddr, "raddr", DefaultRaftAddr, "Set Raft bind address")
	flag.StringVar(&joinAddr, "join", "", "Set join address, if any")
	flag.StringVar(&nodeID, "id", "", "Node ID")
	flag.StringVar(&bootstrapNodeID, "bootstrap-id", "", "Bootstrap Node ID")
	flag.StringVar(&dataDir, "data-dir", "", "Root Data Directory")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if nodeID == "" {
		fmt.Fprintf(os.Stderr, "No nodeID is specified")
		os.Exit(1)
	}
	if bootstrapNodeID == "" {
		fmt.Fprintf(os.Stderr, "no bootstrap nodeID is specified")
		os.Exit(1)
	}
	if dataDir == "" {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}

	// Ensure Raft storage exists.
	raftDir := filepath.Join(dataDir, nodeID)
	if raftDir == "" {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}
	if err := os.MkdirAll(raftDir, 0700); err != nil {
		fmt.Fprintf(os.Stderr, "MakedirAll failed with %v\n", err)
		os.Exit(1)
	}

	s := store.New(inmem)
	s.RaftDir = raftDir
	s.RaftBind = raftAddr
	enableSingle := nodeID == bootstrapNodeID
	log.Printf("EnableSingle %v\n", enableSingle)
	if err := s.Open(enableSingle, nodeID); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
	}

	h := httpd.New(httpAddr, s)
	if err := h.Start(); err != nil {
		log.Fatalf("failed to start HTTP service: %s", err.Error())
	}

	// If join was specified, make the join request.
	if !enableSingle {
		if joinAddr == "" {
			log.Fatalf("joinAddr is not specified")
		}
		if err := join(joinAddr, raftAddr, nodeID); err != nil {
			log.Fatalf("failed to join node at %s: %s", joinAddr, err.Error())
		}
	} else {
		log.Println("Switching to bootstrap mode. Note bootstrapping is done only once!!")
	}

	log.Println("hraftd started successfully")

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	<-terminate
	log.Println("hraftd exiting...")
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	for {
		time.Sleep(interval())
		log.Println("Request to join...", joinAddr)
		if resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr),
			"application-type/json", bytes.NewReader(b)); err != nil {
			log.Printf("http.Post error = %v. Retrying... \n", err)
			resp.Body.Close()
		} else {
			resp.Body.Close()
			return nil
		}
	}
}

func interval() time.Duration  {
	max := 10
	min := 3
	v := rand.Intn(max - min) + min
	return time.Duration(v) * time.Second
}