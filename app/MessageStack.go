package main

// MessageStack is a stack of message (FILO)
type MessageStack []Message

// Push adds to the top of the stack
func (s MessageStack) Push(v Message) MessageStack {
	return append(s, v)
}

// Pop returns a stack without the top element and returns the top elemtns
func (s MessageStack) Pop() (MessageStack, Message) {

	if s == nil || len(s) == 0 {
		return s, Message{}
	}

	l := len(s)
	return s[:l-1], s[l-1]
}

// Everything grabs everything from the stack
func (s *MessageStack) Everything() <-chan Message {
	out := make(chan Message)
	go func() {
		stack := *s
		for {
			st, msg := stack.Pop()
			if msg.Chat == nil {
				close(out)
				return
			}
			stack = st
			out <- msg
		}
	}()
	return out
}
