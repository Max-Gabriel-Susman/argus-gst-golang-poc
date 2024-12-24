package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-gst/go-glib/glib"
	"github.com/go-gst/go-gst/gst"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <input_rtmp_url> <http_output_port>")
		os.Exit(1)
	}

	// Initialize GStreamer
	gst.Init(nil)

	// Create a main loop for handling events
	mainLoop := glib.NewMainLoop(glib.MainContextDefault(), false)

	// Get input RTMP URL and HTTP output port
	inputRTMP := os.Args[1]
	outputPort := os.Args[2]

	// Define the GStreamer pipeline
	// This pipeline takes the RTMP input, processes it, and serves via HTTP using HLS
	pipelineString := fmt.Sprintf(
		"rtmpsrc location=%s ! flvdemux name=demux "+
			"demux.video ! decodebin ! videoconvert ! x264enc bitrate=1000 ! mpegtsmux ! hlssink location=/tmp/segment%%05d.ts playlist-location=/tmp/playlist.m3u8 playlist-length=3 target-duration=2 max-files=5",
		inputRTMP,
	)

	// Create the GStreamer pipeline
	pipeline, err := gst.NewPipelineFromString(pipelineString)
	if err != nil {
		fmt.Println("Failed to create pipeline:", err)
		os.Exit(2)
	}

	// Add a message handler to the pipeline bus
	pipeline.GetPipelineBus().AddWatch(func(msg *gst.Message) bool {
		switch msg.Type() {
		case gst.MessageEOS: // End of stream
			fmt.Println("End of Stream reached")
			pipeline.BlockSetState(gst.StateNull)
			mainLoop.Quit()
		case gst.MessageError: // Error handling
			err := msg.ParseError()
			fmt.Println("ERROR:", err.Error())
			if debug := err.DebugString(); debug != "" {
				fmt.Println("DEBUG:", debug)
			}
			mainLoop.Quit()
		default:
			// Optional: Log other messages for debugging
			fmt.Println("Message:", msg.String())
		}
		return true
	})

	// Set the pipeline state to PLAYING
	fmt.Println("Starting RTMP to HTTP streaming service...")
	err = pipeline.SetState(gst.StatePlaying)
	if err != nil {
		fmt.Println("Failed to start pipeline:", err)
		os.Exit(3)
	}

	// Serve the HLS output via a simple HTTP server
	go func() {
		fmt.Printf("Serving HLS stream on http://localhost:%s/playlist.m3u8\n", outputPort)
		if err := serveHLS(outputPort); err != nil {
			fmt.Println("Failed to serve HLS:", err)
			mainLoop.Quit()
		}
	}()

	// Run the main loop to process events
	mainLoop.Run()

	// Cleanup
	pipeline.SetState(gst.StateNull)
	fmt.Println("Pipeline stopped")
}

// serveHLS starts a simple HTTP server to serve the HLS files
func serveHLS(port string) error {
	httpDir := "/tmp" // Directory containing HLS files
	return http.ListenAndServe(":"+port, http.FileServer(http.Dir(httpDir)))
}
