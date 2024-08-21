```mermaid
erDiagram
    ShortUrl {
        string rash pk "hash is the url 'path' "
        string destination "full destination url"
        %% string status "enum:EXPIRED|ACTIVE"
        int ttl "Unix epoch time format date to expire"
        date updatedAt
        date createdAt
    }
```