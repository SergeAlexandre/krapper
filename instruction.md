

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

Resync with Codebase

Add handlers in back/cmd/serve.go

- GET .../api/v1/wraps/{wrap-name}

    Return a wrap definition

- GET .../api/v1/resources/{wrap-name}

    Retrieve the associated k8s object list. Object is defined by wrap.source.apiVersion/kind
    Also handle selector and namespace if defined.
    clusterScoped means not namespaced.

Small change on the returned GET .../api/v1/resources/{wrap-name} request:
- Remove the k8s list wrapper and return a list of individual objects.
- Remove the managedFields attribute.

Small change on the wrapStore interface: Make getWrap() returning a *wrap.Wrap. And nil of wrap does not exist

In front folder, complete the react and primereact skeleton to build an admin like dashboard with:
- A left menu hosting an entry per item of a catalog. Such catalog to be fetched from http://localhost:7777/api/v1/wraps
- Leave right panel empty for the moment
- A top banner hosting a logo on the left and on the right a 'my account' and settings logo Keep it as simple as possible

For each menu entry where 'menuMode' is 'subMenu', when clicked:

- Fetch the corresponding wrap data with http://localhost:7777/api/v1/resources/ where wrapName is the name in the catalog entry
- Create a sub-menu entry with the name of the fetched resource. Use the backend server which is already running.

When menuMode is 'grid', leave it as it. Will see later

Uploading a screenshot, showing no menu entry.

Menu item now OK. But no action when clicking on Groups entry.
For information, I log the backend and I don't see any request to /api/v1/resources. (But I see request on /api/v1/wraps )

Make the backend resource fetching to occurs on each opening of submenu. Don't cache resources date, as they may change.

Refactor the front code layout by using components 


