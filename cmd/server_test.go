package main

import (
	"bytes"
	"context"
	"encoding/json"
	pb "go-vuejs/proto/helloworld"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func TestOpenAPI(t *testing.T) {
	gwmux := runtime.NewServeMux()
	httpSrv := httptest.NewServer(NewWebServer(gwmux))
	defer httpSrv.Close()

	resp, err := http.Get(httpSrv.URL + "/api/swagger.json")

	if err != nil {
		t.Fatalf("err %v\n", err)
	}

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("err %v. was expecting application/json but got %v\n", err, resp.Header.Get("Content-Type") != "application/json")
	}
}

func TestFileServer(t *testing.T) {
	gwmux := runtime.NewServeMux()
	httpSrv := httptest.NewServer(NewWebServer(gwmux))
	defer httpSrv.Close()

	resp, err := http.Get(httpSrv.URL + "/api/swagger-ui.css")

	if err != nil {
		t.Fatalf("err %v\n", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expecting status code 200 but got %d\n", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != "text/css; charset=utf-8" {
		t.Fatalf("err %v. was expecting application/json but got %v\n", err, resp.Header.Get("Content-Type"))
	}

}

func dialer() func(context.Context, string) (net.Conn, error) {
	lis := bufconn.Listen(1024 * 1024)
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	go func() {
		log.Fatal(s.Serve(lis))
	}()
	return func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

}

func TestHelloWorld_REST(t *testing.T) {
	//curl -X POST -k localhost:8080/hello_world -d '{"name":"yo"}'
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	conn, err := grpc.DialContext(
		ctx,
		"",
		grpc.WithInsecure(),
		grpc.WithContextDialer(dialer()),
	)
	if err != nil {
		log.Fatal("err ", err)
	}
	defer conn.Close()

	gwmux := runtime.NewServeMux()
	err = pb.RegisterGreeterHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}

	httpSrv := httptest.NewServer(NewWebServer(gwmux))
	defer httpSrv.Close()

	payload, err := json.Marshal(map[string]string{
		"name": "yo"})

	resp, err := http.Post(httpSrv.URL+"/hello_world", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatalf("err %v\n", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err %v\n", err)
	}

	type Result struct {
		Message string `json:"message"`
	}
	var rr Result
	err = json.Unmarshal(bodyBytes, &rr)
	if err != nil {
		t.Fatalf("err %v", err)
	}
	if rr.Message != "yo world" {
		t.Fatalf("err. was expecting: yo world but got %v\n", rr.Message)
	}
}
