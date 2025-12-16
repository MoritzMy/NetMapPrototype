# NetMap

NetMap is a simple network discovery tool that creates a map for each available network interface by scanning the local network.

It uses a combination of ARP scans and ICMP ping sweeps to detect live hosts. The ICMP ping functionality is self implemented, allowing for direct interaction with ICMP packets instead of relying on system utilities.

## Features

- Network mapping per network interface

- ARP-based host discovery

- ICMP ping sweeps with a custom ICMP implementation

- Low-level packet handling for learning and experimentation

## Purpose

NetMap is intended for educational use.
The custom ICMP implementation exists to better understand network protocols, raw sockets, and packet-level communication.
To better understand how to build ICMP Packets, no additional libraries were used, except "encodings/binary" from https://pkg.go.dev/encoding/binary


## Disclaimer

This tool is meant for learning and authorized network analysis only.
Do not use it on networks you do not own or have permission to scan.
