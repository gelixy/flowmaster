package blocks

type Block interface {
	LinkDown(any) Block
}
