package x

func FailFast(err error) {
	if err != nil {
		panic(err)
	}
}
