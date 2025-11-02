// Copyright 2018-2025 Celer Network

package testing

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/celer-network/agent-pay/utils"
	"google.golang.org/grpc"
)

type ServerController struct {
	process *os.Process
}

// note address isn't used
func StartServerController(path string, args ...string) *ServerController {
	args = append(args,
		"-routerbcastinterval", "4",
		"-routerbuildinterval", "5",
		"-routeralivetimeout", "8",
		"-ospclearpayinterval", "8")
	process := StartProcess(path, args...)
	// Try to detect the gRPC serving port from args and wait until it's reachable
	var port string
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "-port" {
			port = args[i+1]
			break
		}
	}
	deadline := time.Now().Add(30 * time.Second)
	if port != "" {
		for time.Now().Before(deadline) {
			// If process died, exit early
			if err := process.Signal(syscall.Signal(0)); err != nil {
				fmt.Println("server process exited before listening:", err)
				break
			}
			// First ensure TCP port is open
			if conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 300*time.Millisecond); err == nil {
				conn.Close()
				// Then ensure gRPC/TLS handshake succeeds (server is actually serving)
				dctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				_, gerr := grpc.DialContext(
					dctx,
					"127.0.0.1:"+port,
					utils.GetClientTlsOptionPermissive(),
					grpc.WithBlock(),
				)
				cancel()
				if gerr == nil {
					break
				}
			}
			time.Sleep(250 * time.Millisecond)
		}
	} else {
		// Fallback sleep if port unknown
		time.Sleep(2 * time.Second)
	}
	return &ServerController{process}
}

func (sc *ServerController) Kill() error {
	KillProcess(sc.process)
	return nil
}
