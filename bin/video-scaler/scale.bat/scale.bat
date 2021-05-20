@echo off
ffmpeg -y -i %1 -vf scale=-2:%2  %3