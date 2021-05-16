@echo off

cd videoserver && go build && cd ..
cd thumbgenerator && go build && cd ..
cd videoprocessor && go build && cd ..
cd frontend && go build && cd ..