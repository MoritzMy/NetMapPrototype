package proto

// Marshal marshals the provided Packet of type T into a byte slice.
func Marshal[T Packet](packet T) ([]byte, error) {
	context, err := packet.Marshal()

	if err != nil {
		return nil, err
	}

	headers, err := packet.GetHeaders().Marshal(context)

	if err != nil {
		return nil, err
	}

	b := append(headers, context...)

	return b, nil
}

// Unmarshal unmarshals data into the provided zero Packet of type T. If ctx is non-zero, it indicates the length of the header context to consider during unmarshalling.
func Unmarshal[T Packet](data []byte, zero T) error {
	h := zero.GetHeaders()

	if err := h.Unmarshal(data); err != nil {
		return err
	}

	if err := zero.Unmarshal(data[h.Len():]); err != nil {
		return err
	}
	return nil
}
