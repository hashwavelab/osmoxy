/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/hashwavelab/osmoxy/pb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:9094"
)

var m map[string]bool = make(map[string]bool)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewOsmoxyClient(conn)
	Get1(c)
	//GetWallet(c)
	//SubAllowance(c)
	go Sub(c)
	Sub1(c)
}

func Get(c pb.OsmoxyClient) {
	start := time.Now()
	log.Println("getting")
	r, err := c.GetPoolsSnapshot(context.Background(), &pb.EmptyRequest{}, grpc.MaxCallRecvMsgSize(200*1024*1024))
	if err != nil {
		log.Fatal(err)
	}
	for _, pool := range r.Pools {
		log.Println(pool)
	}
	log.Println(len(r.Pools))
	for _, p := range r.Pools {
		log.Println(p)
	}
	log.Println("get taken", time.Since(start))
}

func Sub(c pb.OsmoxyClient) {
	stream, err := c.SubscribePoolsUpdate(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatal(err)
	}
	ts := time.Now()
	for {
		updates, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		log.Println("SubscribePoolsUpdate", updates.FromBlockNumber, updates.ToBlockNumber, len(updates.Updates), "time since last event block", time.Since(ts))
		// for _, update := range updates.Updates {
		// 	log.Println(update)
		// }
		ts = time.Now()
	}
}

func Get1(c pb.OsmoxyClient) {
	start := time.Now()
	log.Println("getting")
	r, err := c.GetUniV2PairsSnapshot(context.Background(), &pb.EmptyRequest{}, grpc.MaxCallRecvMsgSize(200*1024*1024))
	if err != nil {
		log.Fatal(err)
	}
	for _, pool := range r.Pairs {
		log.Println(pool)
	}
	log.Println(len(r.Pairs))
	for _, p := range r.Pairs {
		log.Println(p)
	}
	log.Println("get taken", time.Since(start))
}

func Sub1(c pb.OsmoxyClient) {
	stream, err := c.SubscribeUniV2PairsUpdate(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		log.Fatal(err)
	}
	ts := time.Now()
	for {
		updates, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
		}
		log.Println("SubscribeUniV2PairsUpdate", updates.FromBlockNumber, updates.ToBlockNumber, len(updates.Univ2Updates), "time since last event block", time.Since(ts))
		// for _, update := range updates.Univ2Updates {
		// 	log.Println(update)
		// }
		ts = time.Now()
	}
}
