package worker

type Stage interface {
	Name() string
	Run(ctx *Context) error
}
