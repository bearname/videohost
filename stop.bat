@echo off

taskkill /F /IM videoserver.exe
taskkill /F /IM userserver.exe
taskkill /F /IM streamingserver.exe
taskkill /F /IM thumbgenerator.exe
taskkill /F /IM video-scaler.exe
taskkill /F /IM notifier.exe
