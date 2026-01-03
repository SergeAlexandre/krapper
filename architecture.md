

# URLS


### GET .../api/v1/wraps

Return the catalog

### GET .../api/v1/wraps/{wrap-name}

Return a wrap definition

### GET .../api/v1/resources/{wrap-name}

Retrieve the associated k8s object set

### PUT .../api/v1/resources/{wrap-name}

Create or update the associated k8s object. 
Data are fields content. 
Templating is performed at server level.
