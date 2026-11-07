package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"server_sfp_sla/internal/bench"
	"server_sfp_sla/internal/measurement"
	"strings"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
func run(ctx context.Context, args []string, out io.Writer) error {
	flags := flag.NewFlagSet("measurement-bench", flag.ContinueOnError)
	flags.SetOutput(out)
	endpoints := flags.String("agents", "http://127.0.0.1:8851,http://127.0.0.1:8852", "two numeric loopback agent URLs")
	output := flags.String("output", "bin/bench-results", "result parent directory")
	mode := flags.String("mode", "udp", "udp or synthetic")
	packets := flags.Int("packets", 20, "packets per agent")
	interval := flags.Duration("interval", 10*time.Millisecond, "interval between completed UDP probes")
	timeout := flags.Duration("timeout", 10*time.Second, "overall controller timeout, <=30s")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	if *timeout < time.Millisecond || *timeout > 30*time.Second {
		return fmt.Errorf("controller timeout must be 1ms..30s")
	}
	c := measurement.DefaultConfig()
	c.Mode = *mode
	c.Packets = *packets
	c.Interval = *interval
	c.Timeout = 5 * time.Second
	if c.Mode == "udp" {
		c.Interface = "loopback"
		c.Workers = 1
		c.LossEvery = 0
	}
	client, err := bench.NewClient(os.Getenv("NQA_BENCH_TOKEN"))
	if err != nil {
		return err
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()
	report, runErr := client.Run(ctx, strings.Split(*endpoints, ","), c)
	if report.ID != "" {
		dest, err := bench.Export(*output, report)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintln(out, dest); err != nil {
			return err
		}
	}
	return runErr
}
