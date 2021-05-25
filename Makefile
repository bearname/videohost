.PHONY: build
build: build-video build-user build-thumbgenerator build-videoprocessor build-video-scaler build-notifier

.PHONY: build-video
build-video:
	go build  -o .\bin\videoserver\videoserver.exe .\cmd\videoserver\main.go

.PHONY: build-user
build-user:
	go build  -o .\bin\user\userserver.exe .\cmd\user\main.go

.PHONY: build-thumbgenerator
build-thumbgenerator:
	go build  -o .\bin\thumbgenerator\thumbgenerator.exe .\cmd\thumbgenerator\main.go

.PHONY: build-video
build-videoprocessor:
	go build  -o .\bin\videoprocessor\videoprocessor.exe .\cmd\videoprocessor\main.go

.PHONY: build-notifier
build-notifier:
	go build  -o .\bin\notifier\notifier.exe .\cmd\notifier\main.go

.PHONY: build-video-scaler
build-video-scaler:
	xcopy /f .\cmd\video-scaler\scale.bat .\bin\video-scaler\scale.bat  /Y
	xcopy /f .\cmd\video-scaler\resolution.bat .\bin\video-scaler\resolution.bat /Y
	go build  -o .\bin\video-scaler\video-scaler.exe .\cmd\video-scaler\main.go

.PHONY: run
run:
	run.bat

.PHONY: run-videoserver
run-videoserver:
	.\bin\videoserver\videoserver.exe

.PHONY: run-user
run-user:
	.\bin\user\userserver.exe

.PHONY: run-thumbgenerator
run-thumbgenerator:
	.\bin\thumbgenerator\thumbgenerator.exe

.PHONY: run-video-scaler
run-video-scaler:
	.\bin\video-scaler\video-scaler.exe

.PHONY: run-notifier
run-notifier:
	.\bin\notifier\notifier.exe

.PHONY: run-web
run-web:
	cd ./web/ && start npm run serve && cd ../

.PHONY: stop
stop:
	stop.bat
