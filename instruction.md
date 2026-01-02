

Resync with Codebase

Complete the wrapstore.wrapStore implementation, by loading all wrap.Wrap objects from the path argument.
- Load recursively
- Watch file change to reload corresponding file
- use the wrap.Load function
- Consider all '*.yaml' file are Wrap entity. In case of error, issue a Warning in the provided logger, but continue to process other files
- Load can also return nil, nil. This will means the file is not a wrap. Skip it, issuing a warning
- Implement WrapStore interface
