# NetMap

NetMap is a small, educational network discovery tool that builds a network map for each available network interface by scanning the local network. It combines ARP scanning and ICMP ping sweeps (with a custom, minimal ICMP implementation) to detect live hosts and collect basic information about them. NetMap is intended as a learning project for people who want to understand raw sockets, packet construction, and low-level network I/O.

- Repository: [MoritzMy/NetMap](https://github.com/MoritzMy/NetMap)

Table of contents
- About
- Features
- Safety & Legal
- Requirements
- Quick start
- Typical usage
- How it works (high level)
- Limitations
- Contributing
- License
- Acknowledgements

---

## About
NetMap is not meant to replace production-grade network scanners. Instead, it is focused on education: showing how ARP and ICMP can be implemented and used to discover hosts on a LAN using raw packet handling in Go.

---

## Features
- Per-interface network mapping (scans each available interface)
- ARP-based host discovery
- ICMP ping sweeps using a custom ICMP implementation (packet construction and parsing done manually)
- Hands-on, low-level packet handling suitable for learning and experimentation
- Minimal external dependencies (primarily Go standard library)

---

## Safety & Legal
- Use this tool only on networks you own or on which you have explicit permission to perform scanning.
- Network scanning can be disruptive; run during maintenance windows or on isolated lab networks whenever possible.
- The author and contributors are not responsible for misuse.


---
## Requirements
- Go compiler (Go 1.16+ recommended)
- Root / elevated privileges are required to open raw sockets and send crafted packets on most OSes (Linux, macOS may require sudo). Running as non-root will typically fail when trying to create raw sockets or perform ARP/ICMP operations.
- A machine with at least one network interface connected to the network you want to scan.

---

## Quick start (build and run)
1. Clone the repository:
   git clone https://github.com/MoritzMy/NetMap.git
   cd NetMap

2. Build:
   go build -o netmap ./...

3. Run (examples below; most operations require root):
   sudo ./netmap
   sudo ./netmap --help

---

## Typical usage examples
- Scan all available interfaces (default behavior):
  sudo ./netmap

- Scan a single interface (replace `eth0` with your interface name):
  sudo ./netmap --interface eth0

- Get help / options:
  ./netmap --help

Note: The exact flags and options available in the binary can be seen with the `--help` flag. Example flags commonly provided by network tools include interface selection, timeouts, concurrency limits, and output formatting; check the binary for supported flags.

---

## How it works (high level)
- Interface discovery: NetMap enumerates available network interfaces and determines their network addresses and masks.
- ARP scanning:
  - The scanner crafts ARP requests for addresses in the local subnet and listens for ARP replies.
  - ARP discovery is fast and reliable for hosts on the same Ethernet segment (it does not require the target to respond to ICMP).
- ICMP ping sweeps:
  - NetMap implements a minimal ICMP echo request and echo reply flow using raw sockets (rather than relying on an external library).
  - This allows learning about ICMP packet structure (headers, checksums) and low-level send/receive behavior.
- Results aggregation:
  - For each interface the tool produces a map of discovered hosts (IP, MAC where available, and which mechanism discovered them).

---

## Limitations
- Elevated privileges required: raw sockets and ARP operations need root permissions.
- Local-network only: ARP discovery works only on the same broadcast domain (LAN).
- ICMP-based discovery can be blocked by host firewalls or devices that do not respond to ping.
- Platform differences: raw socket behavior and available features differ across operating systems — testing has typically been done on Unix-like systems; behavior on Windows is not guaranteed.
- Not a hardened scanner: NetMap is built for experimentation and education rather than performance, stealth, or resilience against active defenses.

Security considerations
- Sending ARP requests and ping sweeps can be noisy; avoid running this on production or sensitive networks without permission.
- Do not run NetMap on networks where active scanning is prohibited by policy.

---

## Contributing
Contributions, fixes, and improvements are welcome. Suggested ways to help:
- Improve documentation (usage examples, architecture diagrams, CLI docs)
- Add tests for packet building/parsing and scanning logic
- Improve cross-platform support and robust error handling
- Add safe export formats (JSON, CSV) and options to control verbosity/concurrency

When contributing:
- Open an issue to discuss larger changes first.
- Follow idiomatic Go style and add tests for new functionality where possible.
- Keep changes focused and well-documented.

---

## License
NetMap is provided for educational purposes. Check the repository for a LICENSE file for full license terms.

---

## Acknowledgements
- Built to explore raw sockets, ICMP, and ARP in Go.
- The project uses Go's standard library; packet checksum and binary handling uses packages such as `encoding/binary`.

---

## Contact / Author
- Repository owner: MoritzMy
- File an issue on GitHub for bugs, feature requests, or questions: https://github.com/MoritzMy/NetMap/issues

---

## Disclaimer
This tool is intended for authorized network analysis and education only. Do not use it on networks you do not own or have explicit permission to scan. The author is not responsible for misuse or any consequences arising from unauthorized use.
