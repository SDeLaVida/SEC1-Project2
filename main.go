package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"example.com/sec_2/MPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/*
Some of the go specific semantics are inspired by
  - Thor's (from the SWU discord) guide on TLS in GO (LastnameUnknown but is attending the course)
  - https://github.com/NaddiNadja/grpc101
  - https://github.com/NaddiNadja/peer-to-peer
*/

func main() {
	hospitalID := int32(5000) // The hospital is always at port 5000, this also reflected in the flow of the program.

	// Peer to peer communication with grpc boilerplate
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + hospitalID
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:         ownPort,
		clients:    make(map[int32]MPC.MPCClient),
		ctx:        ctx,
		HospitalID: hospitalID,
		chunks:     make(map[int32]int),
		channel:    make(chan Result, 1),
	}

	// Sets up the server
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", ownPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// TLS Certification. It is however self signed, which is fine for this assignment
	serverCert, err := credentials.NewServerTLSFromFile("certificate/server.crt", "certificate/priv.key")
	if err != nil {
		log.Fatalln("failed to create cert", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(serverCert))
	MPC.RegisterMPCServer(grpcServer, p)

	// Starts the server
	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()

	// Connects to the other patients
	for i := 0; i <= 3; i++ {
		port := hospitalID + int32(i)

		if port == ownPort {
			continue
		}

		//Set up client connections with TLS
		clientCert, err := credentials.NewClientTLSFromFile("certificate/server.crt", "")
		if err != nil {
			log.Fatalln("failed to create cert", err)
		}

		fmt.Printf("Trying to dial: %v\n", port)

		// Dial the server, and store the connection in the clients map
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", port), grpc.WithTransportCredentials(clientCert), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Patient (ID: %v):did not connect: %s", p.id, err)
		}
		defer conn.Close()
		c := MPC.NewMPCClient(conn)
		p.clients[port] = c
		fmt.Printf("%v", p.clients[port])
	}

	scanner := bufio.NewScanner(os.Stdin)
	if ownPort != hospitalID { // The following logic is at the patient side
		fmt.Print("Enter a number that you want to share with the other patients. This can overflow so be nice :)\nNumber: ") // This can overflow, so be nice
		for scanner.Scan() {
			secret, _ := strconv.ParseInt(scanner.Text(), 10, 32)
			p.ShareChunks(int(secret)) // Break the secret into chunks and send them to the other patients
			break
		}
		// Wait for the other patients to send their chunks
		current_chunks := -1 // Ignore this: just some logic for the print statement while we wait
		for current_chunks < len(p.clients)-1 {
			if current_chunks != len(p.chunks) {
				current_chunks = len(p.chunks)
				fmt.Printf("Patient (ID: %v): Waiting for data from other patients (%v out of %v)\n", p.id, len(p.chunks)+1, len(p.clients))
			}
			time.Sleep(2 * time.Second)
		}

		// Now we have all the chunks, so we can sum them and send them to the hospital
		sumChunks := p.SumChunks()
		fmt.Printf("Patient (ID: %v): I got all the data from the other patients!\nMy secret is: %v\n", p.id, sumChunks)
		fmt.Printf("Patient (ID: %v): Sending my secret to the hospital at id: %v\n", p.id, p.HospitalID)
		hospitalClient := p.clients[p.HospitalID]
		request := &MPC.Result{Id: p.id, Result: int32(sumChunks)}
		reply, err := hospitalClient.SendResult(p.ctx, request)

		if err != nil {
			log.Fatalf("Patient (ID: %v): Could not send chunk to id:%v : %v", p.id, p.HospitalID, err)
		}
		fmt.Printf("Patient (ID: %v): I got reply from id %v: %v\n", p.id, p.HospitalID, reply)
		return
	} else { // The following logic is aimed at the hospital side

		for (len(p.chunks)) < len(p.clients) {
			// Wait for all the patients to send their chunks. Using a channel to avoid race conditions
			fmt.Printf("Hospital (ID: %v): Waiting for data from other patients. Got %v out of %v\n", p.id, len(p.chunks), len(p.clients))
			time.Sleep(2 * time.Second)
			result := <-p.channel
			fmt.Printf("Hospital (ID: %v): I just got result %v from id %v\n", p.id, result.Result, result.Id)
			p.chunks[result.Id] = int(result.Result)
		}
		for id, chunk := range p.chunks {
			fmt.Printf("Hospital (ID: %v): Patient (Id: %v) sent the following result: %v\n", p.id, id, chunk)
		}
		fmt.Print("Press any button to quit\n")
		scanner.Scan()

	}

}

func (p *peer) SendChunk(ctx context.Context, req *MPC.Request) (*MPC.Reply, error) {
	// Receive chunks from other peers and store them in a the chunks array. Remember GRPC so the method name is backwards.
	id := req.Id
	p.chunks[id] = int(req.Chunk)
	rString := fmt.Sprintf("Patient (ID: %v): I got %v from Patient (ID: %v)", p.id, p.chunks[id], id)
	return &MPC.Reply{Reply: rString}, nil
}

func (p *peer) SendResult(ctx context.Context, req *MPC.Result) (*MPC.Reply, error) {
	// Receive result from other peers and store them in a the chunks array.
	// Should in a perfect world only be called for the hospital connection.
	p.channel <- Result{Id: req.Id, Result: req.Result}
	rString := fmt.Sprintf("Hopsital (ID: %v): I got %v from Patient (ID: %v)", p.id, req.Result, req.Id)
	return &MPC.Reply{Reply: rString}, nil
}

func splitChunks(secret int, n int) []int {
	// Split the secret into chunks, while avoiding floating point numbers but still keeping the sum of the chunks equal to the secret
	temp_secret := secret
	chunks := []int{}
	for i := 0; i < n-1; i++ { // n-1 because we want to keep the last chunk as the remainder (might leak information? Unsure how to avoid this)
		rand := rand.Intn(temp_secret)
		chunks = append(chunks, rand)
		temp_secret -= rand
	}
	chunks = append(chunks, temp_secret)
	return chunks

}

func (p *peer) SumChunks() int {
	// Sum the chunks
	sum := int(0)
	for _, chunk := range p.chunks {
		sum += chunk
	}
	return sum
}

func (p *peer) ShareChunks(secret int) {
	// Split the secret into chunks
	patient_amount := len(p.clients) - 1 // minus one as the hospital is not a patient but is still in the clients map
	split := splitChunks(secret, patient_amount)
	for id, client := range p.clients {
		if id == p.HospitalID {
			fmt.Printf("Patient (ID: %v): I skip sending a chunk to the hospital at id: %v\n", p.id, id)
			continue
		}
		split_id := int(id) % patient_amount // Our port number is 5000 + id, so we need to mod it with the number of chunks. A consequence of the implementation (splitChunks function)
		if id == p.id {
			fmt.Printf("Patient (ID: %v): I keep chunk %v\n", p.id, split[split_id])
			p.chunks[id] = split[split_id]
		}
		request := &MPC.Request{Chunk: int32(split[split_id]), Id: p.id}
		reply, err := client.SendChunk(p.ctx, request)
		if err != nil {
			log.Fatalf("Patient (ID: %v): Could not send chunk to patient (id:%v): %v", p.id, id, err)
		}
		fmt.Printf("Patient (ID: %v): I got reply from patient (id %v): %v\n", p.id, id, reply)
	}
}

type peer struct {
	MPC.UnimplementedMPCServer
	id         int32
	clients    map[int32]MPC.MPCClient
	HospitalID int32
	ctx        context.Context
	chunks     map[int32]int
	channel    chan Result
}
type Result struct {
	// Mirrored of the MPC.Result struct in the proto file.
	// Redundant but otherwise we would get warning from the compiler.
	Id     int32
	Result int32
}
