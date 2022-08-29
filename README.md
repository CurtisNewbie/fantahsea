# fantahsea

Fantahsea is a simple gallery/image hosting backend service. It's built using GO, GIN, GORM and other beautiful libraries. It's a simple, personal project, so it won't be very useful for you :D

Fantahsea depends on my other projects: the auth-service for authentication and user management; auth-gateway is just a gateway; and file-service for transferring the already hosted files to fantahsea for browsing.

Plus, the frontend that talks to this app (as well as other services) is in repository [file-service-front](https://github.com/CurtisNewbie/file-service-front), and this app is compatible with the v1.1.10 verison. 

# About Thumbnails Generation

Thumbnails are built using linux's `convert` program.

```sh
# 256x means 256 pixels
convert original.png -resize 256x original-thumbnail.png
```

# Todo

- [ ] Integrate with Redis to do some basic *'distributed'* locking 
- [ ] Make Transferring images async (after we have *'distributed'* locking available) 
- [ ] Add some validator
