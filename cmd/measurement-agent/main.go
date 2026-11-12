package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"server_sfp_sla/internal/measurement"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out io.Writer) error {
	c := measurement.DefaultConfig()
	flags := flag.NewFlagSet("measurement-agent", flag.ContinueOnError)
	flags.SetOutput(out)
	output := flags.String("output", "results/go-agent", "result package parent directory")
	config := flags.String("config", "", "strict JSON configuration file (duration fields use integer ns)")
	flags.IntVar(&c.Packets, "packets", c.Packets, "modeled packet count, 1..10000")
	flags.IntVar(&c.Workers, "workers", c.Workers, "bounded workers, 1..16")
	flags.Uint64Var(&c.Seed, "seed", c.Seed, "deterministic model seed")
	flags.DurationVar(&c.Timeout, "timeout", c.Timeout, "session computation timeout")
	flags.DurationVar(&c.Interval, "interval", c.Interval, "modeled packet interval")
	flags.IntVar(&c.LossEvery, "loss-every", c.LossEvery, "model every Nth packet lost; 0 disables loss")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected positional arguments")
	}
	if *config != "" {
		// Avoid ambiguous precedence: a file supplies the complete measurement config.
		conflict := false
		flags.Visit(func(f *flag.Flag) {
			if f.Name != "config" && f.Name != "output" {
				conflict = true
			}
		})
		if conflict {
			return fmt.Errorf("measurement flags cannot be combined with --config")
		}
		f, err := os.Open(*config)
		if err != nil {
			return err
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil {
			return err
		}
		if info.Size() > 65536 {
			return fmt.Errorf("configuration exceeds 64 KiB")
		}
		dec := json.NewDecoder(io.LimitReader(f, 65537))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&c); err != nil {
			return err
		}
		var extra any
		if err := dec.Decode(&extra); err != io.EOF {
			return fmt.Errorf("configuration must contain one JSON object")
		}
	}
	result, runErr := measurement.Run(ctx, c)
	if result.ID == "" {
		return runErr
	}
	dest, err := measurement.Export(*output, result)
	if err != nil {
		return fmt.Errorf("export: %w", err)
	}
	if _, err := fmt.Fprintln(out, dest); err != nil {
		return err
	}
	return runErr
}
