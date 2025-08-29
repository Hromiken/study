package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	x := flag.Int("t", 3, "TimeLiving (seconds)")
	flag.Parse()

	var wg sync.WaitGroup
	job := make(chan int, 100)

	// таймер на завершение
	timer := time.NewTimer(time.Duration(*x) * time.Second)
	defer timer.Stop()

	// запуск воркеров
	workers := 2
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go timerWorkers(job, &wg, i)
	}

	// генерация данных раз в 500 мс (для ослабления нагрузки на CPU
	data := 0
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timer.C:
			fmt.Println("⏰ Время вышло")
			close(job)
			wg.Wait()
			fmt.Println("✅ Все воркеры завершены")
			return
		case <-ticker.C:
			job <- data
			data++
		}
	}
}

func timerWorkers(ch chan int, wg *sync.WaitGroup, i int) {
	defer wg.Done()
	for v := range ch {
		log.Printf("Рабочий #%d получил: %d", i, v)
	}
	log.Printf("Рабочий #%d завершил работу", i)
}
