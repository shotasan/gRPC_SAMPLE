package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"grpc_sample/api/gen/api"
	"grpc_sample/api/handler"
)

func main() {
	port := 50051
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failde to listen: %v", err)
	}

	server := grpc.NewServer()
	// 自動生成された関数にgRPCサーバーとリクエストを処理するハンドラを渡す
	api.RegisterPancakeBakerServiceServer(
		server,
		handler.NewBakerHandler(),
	)

	// サーバーリフレクションを使うとprotocで生成したコードを使わなくてもProtocolBufferの定義をサーバーから直接取得したりRPCメソッドを実行できる
	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port: %v", port)
		server.Serve(lis)
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()
}
