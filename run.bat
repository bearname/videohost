@echo off
cd web && start npm run serve && cd ..
cd bin/videoserver && start videoserver.exe && cd ../../
cd bin/thumbgenerator && start thumbgenerator.exe && cd ../../
cd bin/video-scaler && start video-scaler.exe && cd ../../