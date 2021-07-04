package find

import "sync"

func mergeStringChannels(channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)
	output := func(c <-chan string) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
