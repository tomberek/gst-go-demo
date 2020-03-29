package main

/*
#cgo pkg-config: gstreamer-1.0 gstreamer-app-1.0
#include <stdio.h>
#include <gst/gst.h>

void cb_proxy_padadd(GstElement* v, GstPad *v2,gpointer v3);
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"

	//"github.com/notedit/gst"
	"github.com/tomberek/gst"
)

//export cb_proxy_padadd
func cb_proxy_padadd(v *C.GstElement, v2 *C.GstPad, v3 unsafe.Pointer) {
	// Recreate a gst.Pad without having access to private .pad field
	// TODO: allow construction of gst.Pad from C.GstPad and remove unsafe
	pad := (*gst.Pad)(unsafe.Pointer(&v2))
	element := (*gst.Element)(v3)

	fmt.Printf("[ USR ] enter cb_proxy_padadd\n")
	fmt.Printf("[ USR ] element: %+v\n", element)

	capstr := pad.GetCurrentCaps().ToString()
	if strings.HasPrefix(capstr, "audio") {
		sinkpad := convert.GetStaticPad("sink")
		pad.Link(sinkpad)
	}
}

var convert *gst.Element

func main() {
	pipeline, err := gst.PipelineNew("test-pipeline")
	if err != nil {
		panic(err)
	}

	source, _ := gst.ElementFactoryMake("uridecodebin", "source")
	convert, _ = gst.ElementFactoryMake("audioconvert", "convert")
	sink, _ := gst.ElementFactoryMake("autoaudiosink", "sink")

	pipeline.Add(source)
	pipeline.Add(convert)
	pipeline.Add(sink)

	convert.Link(sink)

	source.SetObject("uri", "http://dl5.webmfiles.org/big-buck-bunny_trailer.webm")
	source.SetCallback("pad-added", C.cb_proxy_padadd)

	pipeline.SetState(gst.StatePlaying)
	bus := pipeline.GetBus()
	for {
		message := bus.Pull(gst.MessageError | gst.MessageEos)
		fmt.Println("message:", message.GetName())
		fmt.Printf("message: %+v\n", message.GetType())
		fmt.Printf("message: %s\n", message.GetStructure().ToString())
		if message.GetType() == gst.MessageEos {
			break
		}
	}
}
