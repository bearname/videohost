.PHONY: build
build: build-video build-user build-thumbgenerator build-videoprocessor build-video-scaler

build-video:
	go build  -o .\bin\videoserver\videoserver.exe .\cmd\videoserver\main.go

build-user:
	go build  -o .\bin\user\userserver.exe .\cmd\user\main.go

build-thumbgenerator:
	go build  -o .\bin\thumbgenerator\thumbgenerator.exe .\cmd\thumbgenerator\main.go

build-videoprocessor:
	go build  -o .\bin\videoprocessor\videoprocessor.exe .\cmd\videoprocessor\main.go

build-video-scaler:
	xcopy /f .\cmd\video-scaler\scale.bat .\bin\video-scaler\scale.bat  /Y
	xcopy /f .\cmd\video-scaler\resolution.bat .\bin\video-scaler\resolution.bat /Y
	go build  -o .\bin\video-scaler\video-scaler.exe .\cmd\video-scaler\main.go

run:
	run.bat

run-videoserver:
	.\bin\videoserver\videoserver.exe

run-user:
	.\bin\user\userserver.exe

run-thumbgenerator:
	.\bin\thumbgenerator\thumbgenerator.exe

run-video-scaler:
	.\bin\video-scaler\video-scaler.exe

run-web:
	cd ./web/ && start npm run serve && cd ../

stop:
	stop.bat
