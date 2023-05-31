package blocks

type Path struct {
	Id        string
	In        chan any
	Out       chan any
	point     chan any
	converter func(any) any
}

func NewPath(id string, converter func(any) any) *Path {
	block := &Path{
		Id:        id,
		point:     make(chan any),
		converter: converter,
	}

	block.In = block.point
	block.Out = block.point

	return block
}

func (path *Path) LinkDown(lower any) Block {
	switch block := lower.(type) {
	case *Path:
		go func() {
			for {
				message := <-path.point
				if path.converter != nil {
					block.point <- path.converter(message)
				} else {
					block.point <- message
				}
			}
		}()
		return block
	}

	return nil
}
