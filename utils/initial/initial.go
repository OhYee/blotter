package initial

var (
	initialFunction = make([]func(), 0)
)

func Register(f func()) {
	initialFunction = append(initialFunction, f)
}

func Run() {
	for _, f := range initialFunction {
		go f()
	}
}
