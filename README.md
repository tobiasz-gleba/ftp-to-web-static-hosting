# FTP to web static
🚀 Simple static files hosting server, with FTP 💾 as a backend.

![latest version build](https://github.com/tobiasz-gleba/ftp-to-web-static-hosting/actions/workflows/docker-publish.yml/badge.svg)

- ✅ host your files, and explore them in the browser 💻
- ✅ cache in application memory, you don't have to reach the FTP server each time 🏄‍♀️
- ✅ small docker image (under 10Mb) 🗃️
- ✅ super fast app due to golang usage 🏎️

### 🛫 How to use it?

1. Place your static files into FTP server.
2. Start app with docker:

```sh
docker run -p 80:80 \
-e FTP_HOSTNAME=${your-value} \
-e FTP_USERNAME=${your-value} \
-e FTP_PASSWORD=${your-value} \
ghcr.io/tobiasz-gleba/ftp-to-web-static-hosting
```

### 🔨 Avaliable environmental variables for your config:

```env
FTP_HOSTNAME=localhost
FTP_USERNAME=admin
FTP_PASSWORD=admin
FTP_PORT=21
FTP_BASEDIR="/public"
CACHE_TTL_MINUTES=30
CACHE_SIZE_MB=300
```