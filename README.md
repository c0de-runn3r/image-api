# HTTP API for uploading, optimizing, and serving images

## Usage:
To use the server API you need:
- Installed and running RabbitMQ (default port :5672)
- Installed and running Redis (default port :6379)
- run main.go

## Supported requests:
`/start` (GET)
returns basic html page to test API requests

`/upload` (POST)
for uploading your images, returns ID

`/download` (GET)
for getting images back

| Parameter | Required |  Description |
| ----------- | ----------- |----------- |
| id | Yes | image ID |
| quality | Optional | image quality (*100, 75, 50, 25*) |

## Supported features:
- receives and stores images
- optimizes them to three smaller-size image variants (75%, 50%, 25%)
- receives and optimizes images in goroutines
- sends upon request original or optimized image
- stores image names and id's in redis