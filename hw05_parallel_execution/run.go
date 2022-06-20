package hw05parallelexecution

import "errors"

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (err error) {
	cTasks := make(chan Task, 1)
	cErrors := make(chan error)
	active := n // Количество задачь в статусе выполнения.
	l := len(tasks)
	if l < active {
		active = l
	}
	for i := 0; i < active; i++ {
		cTasks <- tasks[i] // Загрузка конвейера.
		go func() {        // Запуск конкурентных исполнителей.
			for t := range cTasks {
				cErrors <- t()
			}
		}()
	}
	for i := active; i < l; i++ { // Конвейер обслуживания конкурентных исполнителей.
		cTasks <- tasks[i]
		if <-cErrors != nil {
			m--
			if m == 0 {
				err = ErrErrorsLimitExceeded
				break
			}
		}
	}
	close(cTasks) // Остановка конвейера.
	for ; active > 0; active-- {
		<-cErrors // Разгрузка конвейера
	}
	return
}
