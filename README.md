# NetMap

## About


NetMap is a network topology exploration tool that uses different Protocol Scans (ICMP, ARP, ...) to infer connections between devices within a network.

## Features

NetMap only supports IPv4 as of now

- ARP-Scan

  Iterates over all Interfaces and sends an ARP Request for all possible IPs in the Subnet of the Interface IP Adress
  
- Ping Sweep

  Iterates over all Interfaces and Ping Sweeps each Subnet of the Interface IP Adress

## Usage
1. Clone the repository:
   ```bash
   git clone https://github.com/MoritzMy/NetMap
    ```
   
2. Navigate to the project directory:
   ```bash
   cd NetMap
   ```
   
3. Build the project:
   ```bash
   go build -o netmap netmap.go
   ```
   
4. Run the tool with appropriate permissions: requires `sudo` for raw socket access:
   ```bash
   sudo ./netmap <flags>
    ```
## Flags

- `-arp-scan` : Perform an ARP scan on all interfaces
- `-ping-sweep` : Perform a ping sweep on all interfaces
- `-dot-file <output-file-name>` : Name of the exported .dot file for graphviz
   
## Project Status

**WIP**

NetMap is under active development.  
Features, APIs, and behavior may change as the project evolves.

## Limitations

NetMap infers topology based on active probing and therefore has inherent limitations:

- Firewalls and routers may block or rate-limit ICMP traffic
- Not all devices respond to traceroute probes
- NAT, VLANs, and Layer 2 topology are not visible
- Routing paths may be asymmetric or change over time

As a result, the generated topology represents a **best-effort approximation**, not a definitive or complete network map.


## Intended Use

- Learning low-level networking concepts
- Experimenting with ICMP and traceroute mechanics
- Exploring network measurement techniques
- Educational demos and lab environments


## Ethical & Legal Notice

Only run NetMap on networks you **own or are explicitly authorized to test**.  
Active network probing may be considered intrusive or malicious in some environments.
