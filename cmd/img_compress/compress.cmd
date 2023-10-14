mkdir processed
for %%F in (*.jpg) do (
    ffmpeg.exe -i "%%F" -q:v 10 "processed\%%F"
)