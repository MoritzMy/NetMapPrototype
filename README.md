# NetMap

## About

> [!CAUTION]
> This tool is built for educational purposes and to gain a deeper understanding of low-level networking concepts.
> It is **not intended for production or enterprise environments**.

NetMap is a network topology exploration tool that uses **ICMP-based traceroutes** to infer connections between devices within a network.

The primary goal of this project is **learning and experimentation**. Core networking mechanisms such as traceroute logic, ICMP packet marshaling/unmarshaling, response parsing, and route construction are **implemented manually** to provide deeper insight into how these protocols work at a low level.

Rather than relying on existing high-level libraries, NetMap focuses on understanding:
- How ICMP Echo and Time Exceeded messages work
- How TTL-based probing reveals routing paths
- How hop-by-hop paths can be aggregated into a graph representation

## Features

NetMap only supports IPv4 as of now

- ARP-Scan

  Iterates over all Interfaces and sends an ARP Request for all possible IPs in the Subnet of the Interface IP Adress
  
- Ping Sweep

  Iterates over all Interfaces and Ping Sweeps each Subnet of the Interface IP Adress

  
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
