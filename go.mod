module github.com/bearname/videohost

go 1.16

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/bearname/videohost/videoserver v0.0.0-20210517041336-153c0acc2878
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.2.0
	github.com/gorilla/context v1.1.1
	github.com/gorilla/mux v1.8.0
	github.com/jlaffaye/ftp v0.0.0-20210307004419-5d4190119067
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/streadway/amqp v1.0.0
	github.com/swaggo/http-swagger v1.0.0
	github.com/swaggo/swag v1.7.0
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
)
