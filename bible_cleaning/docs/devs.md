
# BFS Webscraper

Initially the main thread has one URL to do:

```go
urlCh := make(chan string, 100) // buffered
tasks.Add(1)
urlCh <- startURL
```

Worker goes through the list.

```go
for url := range ctx.urlCh {
    //...
    ctx.tasks.Add(1)
    ctx.urlCh <- nextURL
}
```

Go routine channels are akin to pipes, yet they are synchronized across threads. So, the first worker that reads the URL depletes the channel.
This forces other threads to sleep until the channel has more content or it closes.

Once the first worker finishes processing, it adds another to a queue for another worker to process.

## Multi-threading notes

Here are some notes when using multiple multithreading strategies to optimize processing

- N Language Choose 2 Workers: 136ms
- 8 workers: 132ms
