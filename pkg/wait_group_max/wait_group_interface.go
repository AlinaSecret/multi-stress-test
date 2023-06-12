package wait_group_max

type IWaitGroup interface {
	Add(delta int)
	Wait()
	Done()
}
