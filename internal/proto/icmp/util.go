package icmp

// computeChecksum computes the checksum of a package, by splitting it up into 16 Bit words,
// adding those words together and performing an end around carry until the sum is also a 16 Bit word.
// In the Case of ICMP, while the Checksum is not computed, a Placeholder should be used of which the 16 Bit
// word value is 0
func computeChecksum(request []byte) uint16 {
	sum := uint32(0)

	// Turn the bytes into 16 Bit Words and add them up
	for i := 0; i+1 < len(request); i += 2 {
		sum += (uint32(request[i]) << 8) + uint32(request[i+1])
	}

	if len(request)%2 != 0 {
		sum += uint32(request[len(request)-1]) << 8
	}

	// sum needs to be a valid uint16, otherwise an end around carry is performed
	for sum>>16 != 0 {
		sum = uint32(uint16(sum)) + sum>>16
	}

	// One's complement
	var checksum = ^uint16(sum)

	return checksum
}
