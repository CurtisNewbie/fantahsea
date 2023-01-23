# fantahsea V1.0.3.5

Fantahsea is a simple gallery/image hosting backend service. It's built using GO, GIN, GORM and other beautiful libraries. It's a simple, personal project, so it won't be very useful for you :D

Fantahsea depends on my other projects: the auth-service for authentication and user management; auth-gateway is just a gateway; and file-service for transferring the already hosted files to fantahsea for browsing.

### Requirements

- Consul
- MySQL
- Redis
- RabbitMQ
- file-service: [file-service >= v1.2.5.4](https://github.com/CurtisNewbie/file-server/tree/v1.2.5.4)
- file-service-front (Angular Frontend): [file-service-front >= v1.2.0](https://github.com/CurtisNewbie/file-service-front/tree/v1.2.0)

### About Thumbnails Generation

Thumbnails are built using linux's `convert` program.

```sh
# 256x means 256 pixels
convert original.png -resize 256x original-thumbnail.png
```

### Configurations

See https://github.com/CurtisNewbie/gocommon for more information.


### Changes

Since v1.0.3.5, fantahsea nolonger serves the original images. The original images before compression are served by file-service. Fantahsea only services the generated thumbnails.