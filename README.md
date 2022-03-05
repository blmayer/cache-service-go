# Cache service

> An example cache service using MongoDB ChangeStreams to keep
> it in sync.


## Dependencies

- MongoDB v4.2 or later


## Using

First change the names of the database and collection, and then
set env variable `CONNSTRING` pointing to your mongodb service.
And run the program.


## Deploying

This service is good to be uploaded to cloud, I normally use it
as a cloud run service, but it can be used as a cloud function with
just small changes.

Check my other repo [go-templates](https://github.com/blmayer/go-templates.git)
for examples on how to deploy it.


## MongoDB docs

- [https://docs.mongodb.com/manual/changeStreams/](changeStreams)
- [https://docs.mongodb.com/manual/reference/change-events/](change-events)

