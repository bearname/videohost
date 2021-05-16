@echo off

cd front && start npm run serve && cd ..
cd frontend && start frontend.exe && cd ..
cd videoserver && start videoserver.exe && cd ..
cd thumbgenerator && start thumbgenerator.exe && cd ..