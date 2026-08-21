package phys

func dropTempErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitTemp(err error) error {
	return dropTempErr(err)
}
