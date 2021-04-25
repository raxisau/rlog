# RLog - Yet another Logger

My logger writes the log to file or stdout on goroutine so that there is no impact on the program.
So feel free to send LOTS of debug

```go
...
    rlog.LogSetup("/var/log/myserverlog.log", rlog.DEBUGLEVELNAME)
    rlog.Info("Routing loading, Server started, waiting for requests...")
    http.ListenAndServe(":9000", router)
    rlog.Info("Exiting")
    rlog.CloseLogging()
...
```
