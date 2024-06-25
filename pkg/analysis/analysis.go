package analysis

type State struct {
	documents map[string]string
}

func NewState() State {
	return State{documents: map[string]string{}}
}

func (s *State) OpenDocument( uri, text string) {
    s.documents[uri] = text
}

func (s *State) UpdateDocument( uri, text string) {
    s.documents[uri] = text
}

