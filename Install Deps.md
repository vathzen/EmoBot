Install DCA First

    git clone https://github.com/liclac/dca.git
    sudo apt install cmake build-essentials opus-tools libav-tools libopus-dev libavformat-dev
    cmake .
    make

Convert mp4 to mp3:
ffmpeg -i JR.mp4 -vn -ac 2 -ab 160k -ar 48000 audio.mp3

Cut mp3 to correct part:
ffmpeg -i audio.mp3 -ss 23.5 -t 2.5 -c copy audio2.mp3

Convert to DCA using golib:
ffmpeg -i /mnt/g/go/dca/bin/audio2.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | ./dca > jrm.dca