# Asumsi Livestream
This is a temporary name for our next product, Asumsi livestreaming platform. A platform to help medias, brands, and content creators to go live and interact with their audiences.


This is a starting point for all documents, references, designs, etc.

### Discussion
- [Competitor references](https://github.com/asumsi/platform/discussions/1021)


## Deploying in local
1. clone this repo
2. run `go mod vendor`
3. Copy `.env.example` to `.env` and set the environment variables in `.env` (ask other engineers if you don't already know the variables)
4. run `docker-compose up` (or `docker-compose up --build` if you need to rebuild the images)

For the time being, the docker instance does not support live-reloading. So, never forget the reload/restart the server every time there's code changes happening (this can be done with `docker-compose restart livestream`)