package main

//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path=.. Alastor/alastor.proto

// for build:
//  fyne package --release -os windows -icon ../gvaw-sq-16.png --appID net.dalva.GvawSkinSync --appVersion 0.3.1

import (
	"net.dalva.GvawSkinSync/logger"
	"net.dalva.GvawSkinSync/ui"
)

func init() {
	logger.InitializeLoggerOnceToFile("skinsync.log")
}

func main() {

	ui.ShowUI()
	/*

		addr := serverAddr + ":" + strconv.Itoa(conf.Cfg.Port)

		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)

		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.GetMessage())
	*/
}
