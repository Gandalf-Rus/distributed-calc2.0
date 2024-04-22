package work

import (
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// Интерфейс надо реализовать объектам, которые будут обрабатываться параллельно
type Task interface {
	// таски должны уметь выполнятся
	Do()
}

// Пул для выполнения
type Pool struct {
	// из этого канала будем брать задачи для обработки
	tasks chan Task
	// для информации сколько задач может принять в обработку пул
	freeGorutines int32
	maxGorutines  int
	// для синхронизации работы (в нашем проекте он не обязателен но как хороший тон)
	wg sync.WaitGroup
}

// при создании пула передадим максимальное количество горутин
func New(maxGoroutines int) *Pool {
	p := Pool{
		tasks:         make(chan Task), // канал, откуда брать задачи
		freeGorutines: int32(maxGoroutines),
		maxGorutines:  maxGoroutines,
	}
	return &p
}

func (p *Pool) Run() {
	// для ожидания завершения
	p.wg.Add(p.maxGorutines)
	for i := 0; i < p.maxGorutines; i++ {
		// создадим горутины по указанному количеству maxGoroutines
		go func(ch <-chan Task) {
			// забираем задачи из канала
			for w := range ch {
				// задержка чтоб каналы не читали одну и ту же записть
				time.Sleep(time.Millisecond)
				//уменьшим счетчик
				atomic.AddInt32(&p.freeGorutines, -1)
				// и выполняем
				log.Println("приступим")
				w.Do()
				log.Println("закончил")
				// выполнили => горутина свободна
				atomic.AddInt32(&p.freeGorutines, 1)
			}
			// после закрытия канала нужно оповестить наш пул
			p.wg.Done()
		}(p.tasks)
	}

	p.wg.Wait()
	log.Println("Task pool: 'ну все я пошел...'")
}

func (p *Pool) CountOfFreeGorutines() int {
	return int(p.freeGorutines)
}

// Передаем объект, который реализует интерфейс Task
func (p *Pool) AddTask(w Task) {
	// добавляем задачи в канал, из которого забирает работу пул
	p.tasks <- w
}

func (p *Pool) Shutdown() {
	// закроем канал с задачами
	close(p.tasks)
}
