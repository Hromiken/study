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

	workers := 2
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(job, &wg, i)
	}

	data := 0
	timeout := time.After(time.Duration(*x)) // канал, который "сработает" через duration

	for {
		select {
		case <-timeout:
			fmt.Println("⏰ Время вышло")
			close(job) // закрываем канал, чтобы воркеры завершились
			wg.Wait()  // ждём всех воркеров
			fmt.Println("✅ Все воркеры завершены")
			return
		default:
			job <- data // отправляем данные
			data++
			time.Sleep(500 * time.Millisecond) // чтобы не перегружать CPU
		}
	}
}

func worker(ch chan int, wg *sync.WaitGroup, id int) {
	defer wg.Done()
	for v := range ch {
		log.Printf("Рабочий #%d получил: %d", id, v)
	}
	log.Printf("Рабочий #%d завершил работу", id)
}
