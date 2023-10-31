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

	// Connects to the other patients and the hospital
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

	// Setup is now done and our secure aggregation with TLS implementation begins
	fmt.Print("\n")
	if ownPort != hospitalID { // The following logic is only applicable for the patients
		patientLogic(p)
	} else { // The following logic is only applicable for the hospital
		HospitalLogic(p)
	}

}

func patientLogic(p *peer) {
	// The following logic is at the patient side
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Patient (ID: %v): Enter a number that you want to share with the other patients. This can overflow so be nice :)\n", p.id) // This can overflow, so be nice
	for scanner.Scan() {
		secret, _ := strconv.ParseInt(scanner.Text(), 10, 32)
		fmt.Printf("Patient (ID: %v): my secret number is %v. Performing secure aggregation\n", p.id, secret)

		p.ShareChunks(int(secret)) // Break the secret into chunks and send them to the other patients
		break
	}

	// Wait for the other patients to send their chunks
	current_chunks := -1 // Ignore this: just some logic for the print statement while we wait
	for current_chunks < len(p.clients) {
		if current_chunks != len(p.chunks) {
			current_chunks = len(p.chunks)
			fmt.Printf("Patient (ID: %v): Waiting for data from other patients (%v out of %v)\n", p.id, len(p.chunks)-1, len(p.clients)-1)
		}
		time.Sleep(2 * time.Second)
	}

	// Now we have all the chunks, so we can sum them and send them to the hospital
	sumChunks := p.SumChunks()
	fmt.Printf("Patient (ID: %v): I got all the data from the other patients!\n", p.id)
	fmt.Printf("Patient (ID: %v): My aggregated value is: %v\n", p.id, sumChunks)
	fmt.Printf("Patient (ID: %v): Sending my secret to the hospital at id: %v\n", p.id, p.HospitalID)
	hospitalClient := p.clients[p.HospitalID]
	request := &MPC.Result{Id: p.id, Result: int32(sumChunks)}
	reply, err := hospitalClient.SendResult(p.ctx, request)

	if err != nil {
		log.Fatalf("Patient (ID: %v): Could not send chunk to id:%v : %v", p.id, p.HospitalID, err)
	}
	fmt.Printf("Patient (ID: %v): I got reply from the hospital (id: %v):\n %v\n", p.id, p.HospitalID, reply)
	time.Sleep(6 * time.Second)
}

func HospitalLogic(p *peer) {
	// The following logic is aimed at the hospital side
	scanner := bufio.NewScanner(os.Stdin)
	for (len(p.chunks)) < len(p.clients) {
		// Wait for all the patients to send their chunks. Using a channel to avoid race conditions
		fmt.Printf("Hospital (ID: %v): Waiting for data from other patients. Got %v out of %v\n", p.id, len(p.chunks), len(p.clients))
		result := <-p.channel
		fmt.Printf("Hospital (ID: %v): I just got result %v from id %v\n", p.id, result.Result, result.Id)
		p.chunks[result.Id] = int(result.Result)
	}
	superChunk := int(0)
	for id, chunk := range p.chunks {
		fmt.Printf("Hospital (ID: %v): Patient (Id: %v) sent the following result: %v\n", p.id, id, chunk)
		superChunk += chunk
	}
	fmt.Printf("Hospital (ID: %v): The total sum for all the patients is %v\n", p.id, superChunk)

	fmt.Print("Press any button to quit\n")
	scanner.Scan()
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
	rString := fmt.Sprintf("Hospital (ID: %v): I got %v from Patient (ID: %v)", p.id, req.Result, req.Id)
	return &MPC.Reply{Reply: rString}, nil
}

func splitChunks(secret int, fromId int, n int) map[int32]int {
	// Split the secret into chunks, while avoiding floating point numbers but still keeping the sum of the chunks equal to the secret
	n = n - 1 // Minus one because our implementation is 0 indexed
	remainder := secret
	chunks := make(map[int32]int)
	for i := 0; i < n; i++ {
		rando := 0
		if secret < n { 
			// 0 is ok to send if the user has chosen a secret less than n
			// 0 can lead to an edge case where the secret value can be deduced by the hospital/patients when every patient gets unlucky and sends 0 chunks to n-1 patients. It actually happens pretty often. Try inputting 1 from each patient.
			rando = rand.Intn(remainder)
		} else {
			// This is preferred but not usable if the value is less than n. 
			// I probably should have used floating point numbers.
			rando = improved_rand(remainder)
		}
		chunks[int32((fromId + i))] = rando
		fmt.Printf("Iteration %v (for id:%v): Generated following chunk: %v from the remainder: %v and our secret was: %v\n", i, fromId+i, rando, remainder, secret)
		remainder -= rando
	}
	fmt.Printf("Iteration %v (for id:%v): Parsing the remainder: %v and our secret was: %v\n", n, fromId+n, remainder, secret)
	chunks[int32((fromId + n))] = remainder // The remainder gets parsed to the last peer
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
	patient_amount := len(p.clients) // We don't want a chunk for the hospital, but as our peer is not in the client list, it evens out.
	split := splitChunks(secret, int(p.HospitalID)+1, patient_amount)

	fmt.Printf("Patient (ID: %v): i take my share (chunk: %v) of the chunks\n", p.id, split[p.id])
	p.chunks[p.id] = split[p.id] // I take my share

	for id, client := range p.clients {
		if id == p.HospitalID {
			fmt.Printf("Patient (ID: %v): I skip sending a chunk to the hospital at id: %v\n", p.id, id)
			continue
		}

		request := &MPC.Request{Chunk: int32(split[id]), Id: p.id}
		reply, err := client.SendChunk(p.ctx, request)
		if err != nil {
			log.Fatalf("Patient (ID: %v): Could not send chunk to patient (id:%v): %v", p.id, id, err)
		}
		fmt.Printf("Patient (ID: %v): I got reply from patient (id %v):\n %v\n", p.id, id, reply)

	}
}

func improved_rand(secret int) int {
	randomized := 0
	for randomized < 1 {
		randomized = rand.Intn(secret)
	}
	return randomized

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
