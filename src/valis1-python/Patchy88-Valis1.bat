@echo off
cd /d "%~dp0"
where py >nul 2>nul
if %errorlevel%==0 (
  py -3 Patchy88_Valis1.py
  exit /b %errorlevel%
)
python Patchy88_Valis1.py
