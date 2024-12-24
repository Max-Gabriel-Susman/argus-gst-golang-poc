# Argus GStreamer Golang Proof of Concept

The Argus GStreamer Golang Proof of Concpet is a service I wrote to explore the use of gstreamer in the Argus Stream Engine Service

## Setup

sudo apt-get install libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio

go get "github.com/go-gst/go-gst/gst" "github.com/go-gst/go-glib/glib"

## Usage 

make build

make run

Run the program with an RTMP input and an HTTP output port:
```
go run main.go rtmp://source.endpoint/live/stream 8080
```

Access the HLS stream via:
```
http://localhost:8080/playlist.m3u8
```