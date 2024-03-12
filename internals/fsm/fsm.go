package fsm

type state struct {
	transition map[byte][]*state

	terminal bool
	start    bool
}

func (s *state) Check() bool {
    return false
}
