

Resync with Codebase

In back, complete the internal/wrapstore.wrapStore implementation, by loading all wrap.Wrap objects from the path argument.
- Load recursively
- Watch file change to reload corresponding file
- use the wrap.Load function
- Consider all '*.yaml' file are Wrap entity. In case of error, issue a Warning in the provided logger, but continue to process other files
- Load can also return nil, nil. This will means the file is not a wrap. Skip it, issuing a warning
- Implement WrapStore interface


I think it will be more effective to store the catalog and rebuild it on each file modification.

Log an 'Info' message for each loaded or reloaded file

wrapstore_test.go has been removed ?


Could you explain why I got the following error:
```
cd ack
go test ./internal/httpsrv/
# krapper/internal/httpsrv
internal/httpsrv/httpsrv.go:143:2: call to slog.Logger.Info missing a final value
FAIL    krapper/internal/httpsrv [build failed]
FAIL
```

