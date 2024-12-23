package main

import (
	"fmt"

	"github.com/ziutek/gst"
)

func main() {
	// Initialize GStreamer
	gst.Init(nil)

	// Create a GStreamer pipeline
	pipeline := gst.NewPipeline("pipeline")

	// Create elements
	src := gst.NewElement("videotestsrc", "source")
	enc := gst.NewElement("x264enc", "encoder")
	mux := gst.NewElement("hlsmux", "muxer")
	sink := gst.NewElement("filesink", "sink")

	// Set properties
	sink.SetProperty("location", "output.m3u8")

	// Add elements to the pipeline
	pipeline.Add(src, enc, mux, sink)

	// Link elements
	src.Link(enc)
	enc.Link(mux)
	mux.Link(sink)

	// Start playing the pipeline
	pipeline.SetState(gst.StatePlaying)

	// Wait until the pipeline finishes
	bus := pipeline.GetBus()
	for {
		msg := bus.TimedPopFiltered(gst.CLOCK_TIME_NONE, gst.MESSAGE_EOS|gst.MESSAGE_ERROR)
		if msg != nil {
			switch msg.Type {
			case gst.MESSAGE_EOS:
				fmt.Println("End of Stream")
				pipeline.SetState(gst.StateNull)
				return
			case gst.MESSAGE_ERROR:
				err := msg.ParseError()
				fmt.Printf("Error: %s\\\\n", err)
				pipeline.SetState(gst.StateNull)
				return
			}
		}
	}
}
