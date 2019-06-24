package foo

// mockgen -source=foo.go --destination=../mock_foo/foo.go

type Foo interface {
	Bar(x int) int
	Get()
	Set(x int)
}

func SUT(f Foo) {
	i := f.Bar(99)
	_ = i
}
