@echo off
ffprobe -v error -show_entries stream=width,height -of default=noprint_wrappers=1 %1