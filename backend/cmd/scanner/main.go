package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MoritzMy/NetMap/backend/cmd/arp_scan"
	"github.com/MoritzMy/NetMap/backend/cmd/ping"
	"github.com/MoritzMy/NetMap/backend/internal/graphing"
)

func main() {
	arp := flag.Bool("arp-scan", false, "Run ARP Discovery ScanInterface")
	icmp := flag.Bool("ping-sweep", false, "Run ICMP Sweep")
	dot_file := flag.String("dot-file", "", "Output the resulting graph to a DOT file")
	json_file := flag.String("json-file", "", "Output the resulting json to a file")

	flag.Parse()

	graph := graphing.NewGraph()

	if *arp {
		arp_scan.RunARPScan(graph)
	}

	if *icmp {
		ping.RunICMPSweep(graph)
	}

	if !*arp && !*icmp {
		fmt.Println("Please specify a scan type. Use -h for help.")
	}

	if *json_file != "" {
		json, err := graph.MarshalJSON()

		if err != nil {
			fmt.Printf("failed to marshal graph: %s", err)
		}

		wd, err := os.Getwd()

		outPath := filepath.Join(wd, *json_file)

		file, err := os.Create(outPath)
		if err != nil {
			fmt.Printf("could not create file: %s", err)
		}
		defer file.Close()

		file.Chmod(0644)
		file.Write(json)
	}

	if *dot_file != "" {
		err := graph.ExportToDOT(*dot_file)
		if err != nil {
			fmt.Println("Error exporting graph to DOT file:", err)
		} else {
			fmt.Printf("Graph exported to %s\n", *dot_file)
		}
	}

	fmt.Println(graph)
}
