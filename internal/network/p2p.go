package main

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/SohelAhmedJoni/Awazz-Backend/internal/model"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

type blockHeader struct {
	BlockID string
	PScode  string
}

var (
	topicNameFlag = flag.String("topicName", "/block/1.0.0", "name of topic to join")
)

func p2p_sync() {
	cfg := parseFlags()
	ctx := context.Background()

	db, err := leveldb.OpenFile("BLOCK", nil)
	if err != nil {
		panic(err)
	}
	h, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}
	// h.SetStreamHandler("/block/1.0.0", func(s network.Stream) {
	// 	log.Printf("/block/1.0.0 stream created")
	// 	handleStream(s, db)
	// })

	// go discoverPeers(ctx, h)
	go connectPeers(ctx, h, cfg)

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	topic, err := ps.Join(*topicNameFlag)
	if err != nil {
		panic(err)
	}
	go streamConsoleTo(ctx, topic, db)

	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}
	printMessagesFrom(ctx, sub, h, db)
}
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func connectPeers(ctx context.Context, h host.Host, cfg *config) {

	fmt.Printf("\n[*] Your Multiaddress Is: /ip4/%s/tcp/%v/p2p/%s\n", cfg.listenHost, cfg.listenPort, h.ID().Pretty())

	peerChan := initMDNS(h, cfg.RendezvousString)
	for { // allows multiple peers to join
		peer := <-peerChan // will block untill we discover a peer
		fmt.Println("Found peer:", peer, ", connecting")

		if err := h.Connect(ctx, peer); err != nil {
			fmt.Println("Connection failed:", err)
			continue
		}

	}
	// fmt.Println("Peer discovery complete")
}

// publish
func streamConsoleTo(ctx context.Context, topic *pubsub.Topic, db *leveldb.DB) {
	lbyte, err := db.Get([]byte("latest"), nil)
	check(err)
	var latestBlock model.Post
	err = proto.Unmarshal(lbyte, &latestBlock)
	check(err)

	bit, err := json.Marshal(latestBlock)
	check(err)

	// reader := bufio.NewReader(os.Stdin)
	for {
		// s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if err := topic.Publish(ctx, []byte(bit)); err != nil {
			fmt.Println("### Publish error:", err)
		}
	}
}

func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription, h host.Host, db *leveldb.DB) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}

		if m.ReceivedFrom == h.ID() {
			continue
		}
		// wait for 5 seconds
		time.Sleep(5 * time.Second)

		fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
		ReceivedBlock := BlockHeader{}
		json.Unmarshal(m.Message.Data, &ReceivedBlock)

		latestBlock := LiteBlock{}
		lbyte, err := db.Get([]byte("latest"), nil)
		check(err)
		err = latestBlock.XXX_Unmarshal(lbyte)
		check(err)

		if ReceivedBlock.Pscode == latestBlock.Header.Pscode {
			if ReceivedBlock.BlockID < latestBlock.Header.BlockID {
				stream, err := h.NewStream(ctx, m.ReceivedFrom, protocol.ID("/block/1.0.0"))
				if err != nil {
					fmt.Println("Stream open failed", err)
				} else {
					rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
					go writeData(rw, ReceivedBlock.BlockID, latestBlock.Header.BlockID, db)
					fmt.Println("Connected to:", m.ReceivedFrom.Pretty())
				}
			} else if ReceivedBlock.BlockID > latestBlock.Header.BlockID {
				stream, err := h.NewStream(ctx, m.ReceivedFrom, protocol.ID("/block/1.0.0"))
				if err != nil {
					fmt.Println("Stream open failed", err)
				} else {
					rw := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream))
					go readData(rw)
					fmt.Println("Connected to:", m.ReceivedFrom.Pretty())
				}
			}

		}
	}
}

func sendMsg(rw *bufio.ReadWriter, id int64, content []byte) error {
	return nil
}

func readMsg(rw *bufio.ReadWriter) {
	for {
		// read bytes until new line
		msg, err := rw.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error reading from buffer")
			continue
		}

		// get the id
		id := int64(binary.LittleEndian.Uint64(msg[0:8]))

		// get the content, last index is len(msg)-1 to remove the new line char
		content := string(msg[8 : len(msg)-1])

		if content != "" {
			// we print [message ID] content
			fmt.Printf("[%d] %s", id, content)
		}

		if err := sendMsg(rw, id, []byte("response")); err != nil {
			fmt.Println("Err while sending response: ", err)
			continue
		}
	}
}
