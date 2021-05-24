@echo off
@REM cd web && start npm run serve && cd ..
cd bin/videoserver && start videoserver.exe && cd ../../
cd bin/user && start userserver.exe && cd ../../
cd bin/thumbgenerator && start thumbgenerator.exe && cd ../../
cd bin/video-scaler && start video-scaler.exe && cd ../../
cd bin/notifier && start notifier.exe && cd ../../