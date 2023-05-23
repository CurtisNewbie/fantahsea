# fantahsea V1.0.6

Fantahsea is a simple gallery/image hosting backend service. It's built using GO, GIN, GORM and other beautiful libraries. It's a simple, personal project, so it won't be very useful for you :D

### Requirements

- Consul
- MySQL
- Redis
- RabbitMQ
- [goauth >= v1.0.0](https://github.com/CurtisNewbie/goauth/tree/v1.0.0)
- [bolobao (Angular Frontend) >= v0.0.1](https://github.com/CurtisNewbie/bolobao/tree/v0.0.1)
- [vfm >= v0.0.1](https://github.com/CurtisNewbie/vfm/tree/v0.0.1)
- [mini-fstore >= v0.0.1](https://github.com/CurtisNewbie/mini-fstore/tree/v0.0.1)

### About Thumbnails Generation

Thumbnails are built using linux's `convert` program.

```sh
# 256x means 256 pixels
convert original.png -resize 256x original-thumbnail.png
```

### Configurations

See https://github.com/CurtisNewbie/gocommon for more information.


### Changes

- Since v1.0.3.5, fantahsea nolonger serves the original images. The original images before compression are served by file-service. Fantahsea only services the generated thumbnails.