package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"server_sfp_sla/internal/bench"
	"server_sfp_sla/internal/measurement"
)

func runServe(ctx context.Context, args []string, out io.Writer) error {
	flags := flag.NewFlagSet("measurement-agent serve", flag.ContinueOnError)
	flags.SetOutput(out)
	id := flags.String("id", "agent", "agent identity")
	listen := flags.String("listen", "127.0.0.1:8851", "loopback HTTP endpoint")
	probe := flags.String("probe", "127.0.0.1:9851", "loopback UDP endpoint")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	address, err := measurement.LoopbackAddress(*listen)
	if err != nil {
		return err
	}
	a, err := bench.NewAgent(ctx, *id, *probe, os.Getenv("NQA_BENCH_TOKEN"))
	if err != nil {
		return err
	}
	defer a.Close()
	listener, err := net.Listen("tcp", address.String())
	if err != nil {
		return err
	}
	defer listener.Close()
	if _, err := fmt.Fprintf(out, "agent=%s control=http://%s probe=%s timestamp=software\n", *id, listener.Addr(), a.Info().ProbeAddress); err != nil {
		return err
	}
	return bench.Serve(ctx, listener, a)
}
