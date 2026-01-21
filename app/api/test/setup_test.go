package test

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"testing"
	"time"

	"github.com/PiaoAdmin/pmall/app/api/biz/dal"
	"github.com/PiaoAdmin/pmall/app/api/biz/router"
	mwError "github.com/PiaoAdmin/pmall/app/api/md/error"
	"github.com/PiaoAdmin/pmall/app/api/md/jwt"
	"github.com/PiaoAdmin/pmall/app/api/rpc"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

var testBaseURL string

// TestMain initializes the test environment and starts the test server
func TestMain(m *testing.M) {
	fmt.Println("=== TestMain: Starting test server setup ===")

	// Change to parent directory so config files can be found
	if err := os.Chdir("../"); err != nil {
		fmt.Printf("failed to change directory: %v", err)
		os.Exit(1)
	}

	// Initialize shared dependencies just like main.go does.
	dal.Init()
	rpc.Init()
	jwt.Init()

	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := ln.Addr().String()
	ln.Close() // Close it so Hertz can bind to it

	f, err := os.OpenFile("./output_test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fileWriter := io.MultiWriter(f)
	hlog.SetOutput(fileWriter)

	h := server.New(server.WithHostPorts(addr))
	h.Use(mwError.GlobalErrorHandler())
	router.GeneratedRegister(h)

	go h.Spin()

	testBaseURL = fmt.Sprintf("http://%s", addr)
	fmt.Printf("=== TestMain: Server started at %s ===\n", testBaseURL)
	time.Sleep(200 * time.Millisecond)

	// Run all tests
	code := m.Run()

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	h.Shutdown(ctx)

	os.Exit(code)
}
