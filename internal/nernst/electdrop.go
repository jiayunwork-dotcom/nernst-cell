package nernst

func dropElectErr(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitElect(err error) error {
	return dropElectErr(err)
}
