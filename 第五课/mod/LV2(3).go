package mod

import "sync"

type task func()

type mod struct {
	tasksh chan task
	wg     sync.WaitGroup
}

func NewMod(workerCount int, taskBuffer int) *mod {
	m := &mod{
		tasksh: make(chan task, taskBuffer),
	}
	for i := 0; i < workerCount; i++ {
		go func() {
			for t := range m.tasksh {
				t()
				m.wg.Done()
			}
		}()
	}
	return m
}

func (m *mod) Submit(t task) {
	m.wg.Add(1)
	m.tasksh <- t
}

func (m *mod) Wait() {
	m.wg.Wait()
}

func (m *mod) Close() {
	close(m.tasksh)
}
