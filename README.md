# udpserver

A simple echo UDP server; useful for mocking a backend required for tests. You
could use it for sending StatsD metrics like this:

    s := udpserver.New(8125)
    s.Serve()

    m := metrics.New("127.0.0.1:8125")
    m.Count("hello.bar", 79, metrics.Tag("env", "dev"))

    s.Close()

Which would result in the following on stdout:

    2020/02/26 12:18:01 server listening on: 0.0.0.0:8125
    2020/02/26 12:18:01 received: hello.bar:79|c|#env:dev
    2020/02/26 12:18:01 shutting down...
