package ip

func Unmarshal[T IPPacket](b []byte, zero T) error {
	if err := zero.GetHeaders().Unmarshal(b); err != nil {
		return err
	}

	if err := zero.Unmarshal(b[zero.HeaderSize():zero.Size()]); err != nil {
		return err
	}

	return nil

}
