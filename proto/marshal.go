package proto

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

func Unmarshal[T Packet](data []byte, zero T) error {
	h := zero.GetHeaders()

	// TODO: Before unmarshaling the Headers, the Size of the Headers must be determined
	// May add a function that gets the size of the Headers from a byte array

	// Header Interface must implement a calculate Header Length function which takes []byte as input

	if err := h.Unmarshal(data); err != nil {
		return err
	}
	if err := zero.Unmarshal(data[h.Len():zero.Len()]); err != nil {
		return err
	}
	return nil
}
