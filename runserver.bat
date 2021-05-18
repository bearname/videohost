@echo off

FOR /L %%A IN (%1,1,%2) DO (
    ECHO %%A
    start videoserver.exe %%A
)