package hw05parallelexecution

import "errors"

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) (err error) {
	cTasks := make(chan Task, 1)
	cErrors := make(chan error)
	active := n // Количество задач в статусе выполнения.
	l := len(tasks)
	if l < active {
		active = l
	}

	// Фаза 1 - загрузка конвейера обслуживания конкурентных исполнителей.
	for i := 0; i < active; i++ {
		go func() { // Запуск конкурентных исполнителей.
			for t := range cTasks {
				cErrors <- t()
			}
		}()
		cTasks <- tasks[i] // Загрузка конвейера: каждому исполнителю по задаче.
	}

	// Фаза 2 - работа конвейера: на каждую принятую задачу одна отправленная.
	for i := active; i < l; i++ { // Конвейер обслуживания конкурентных исполнителей.
		cTasks <- tasks[i] // Задачу в буфер, перед тем как проверять зданную, что бы исполнитель не ждал.
		if <-cErrors != nil {
			m--
			if m == 0 {
				err = ErrErrorsLimitExceeded
				break
			}
		}
	}

	// Фаза 3 - разгрузка конвейера.
	close(cTasks)                // Остановка конвейера.
	for ; active > 0; active-- { // Осталось ровно столько, сколько загрузили на первой фазе.
		<-cErrors // Разгрузка конвейера
	}
	return
}
