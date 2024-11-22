# FTP to web static
🚀 Simple static files hosting server, with FTP 💾 as a backend.

### 🛫 How to use it?

1. Provide environmental variables
2. Start with docker: `docker run -p 80:80 <image>`

### 🔨 Avaliable environmental variables for your config:

```env
FTP_HOSTNAME=ftp.example.com
FTP_USERNAME=ftpuser
FTP_PASSWORD=ftppassword
FTP_PORT=21
FTP_TLS=false
```