module github.com/Ruletk/OnlineClinic/apps/doctor

go 1.23.8

require (
    github.com/gin-gonic/gin       v1.10.0
    github.com/joho/godotenv       v1.5.1
    github.com/nats-io/nats.go     v1.17.3     // if you’re using NATS
    gorm.io/gorm                   v1.26.0
    gorm.io/driver/postgres        v1.6.4      // v1.6.1 doesn’t exist, v1.6.4 is the latest
)
