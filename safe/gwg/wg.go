package gwg

import "sync"

/**
When many programmers use WaitGroup, they often use structures incorrectly,
and the correct way is to use pointers.
*/

func New() *sync.WaitGroup {
	return &sync.WaitGroup{}
}
