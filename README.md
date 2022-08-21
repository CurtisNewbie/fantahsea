# fantahsea

Fantahsea is a simple gallery/image hosting backend service. It's built using GO, GIN, GORM and other beautiful libraries. It's a simple, personal project, so it won't be very useful for you :D

Fantahsea depends on my other projects: the auth-service for authentication and user management; auth-gateway is just a gateway; and file-service for transferring the already hosted files to fantahsea for browsing.

Thumbnails are built using linux's `convert` program.

```sh
# 256x means 256 pixels
convert original.png -resize 256x original-thumbnail.png
```
