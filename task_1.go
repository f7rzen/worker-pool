package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    in := make(chan string)
    var wg sync.WaitGroup
    workerCount := 0
    done := make(chan struct{})

    addWorker := func(id int, in chan string, done chan struct{}) {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for {
                select {
                case data, ok := <-in:
                    if !ok {
                        return
                    }
                    fmt.Printf("Рабочий %d получил: %s\n", id, data)
                    time.Sleep(500 * time.Millisecond)
                case <-done:
                    fmt.Printf("Рабочий %d остановлен\n", id)
                    return
                }
            }
        }()
    }

    for i := 0; i < 3; i++ {
        workerCount++
        addWorker(workerCount, in, done)
    }

    for i := 1; i <= 5; i++ {
        in <- fmt.Sprintf("Задача %d", i)
    }

    time.Sleep(1 * time.Second)
    workerCount++
    fmt.Println("Добавление нового рабочего...")
    addWorker(workerCount, in, done)

    for i := 6; i <= 10; i++ {
        in <- fmt.Sprintf("Задача %d", i)
    }

    time.Sleep(2 * time.Second)
    fmt.Println("Остановка одного рабочего...")
    done <- struct{}{}

    time.Sleep(2 * time.Second)

    close(in)
    wg.Wait()
    fmt.Println("Все рабочие остановлены.")
}