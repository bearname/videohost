@echo
migrate -database mysql://root:123@tcp(127.0.0.1:3306)/video -path data/mysql/migration %1 1