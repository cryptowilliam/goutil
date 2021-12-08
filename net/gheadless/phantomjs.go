package gheadless

/**
If run `phantomjs` report error:
qt.qpa.xcb: could not connect to display
qt.qpa.plugin: Could not load the Qt platform plugin "xcb" in "" even though it was found.
This application failed to start because no Qt platform plugin could be initialized. Reinstalling the application may fix this problem.

Available platform plugins are: eglfs, linuxfb, minimal, minimalegl, offscreen, vnc, xcb.

PhantomJS has crashed. Please read the bug reporting guide at
<http://phantomjs.org/bug-reporting.html> and file a bug report.
Aborted (core dumped)


Fix problem with run `QT_QPA_PLATFORM=offscreen phantomjs` instead of `phantomjs`.
*/

import (
	"fmt"
	"github.com/cryptowilliam/goutil/container/grand"
	"github.com/cryptowilliam/goutil/sys/gcmd"
	"github.com/cryptowilliam/goutil/sys/gfs"
	"os"
)

const (
	// modified from https://github.com/ariya/phantomjs/blob/master/examples/rasterize.js
	// reference https://blog.csdn.net/qq_41534566/article/details/83147103
	rasterizeJS = `"use strict";
var page = require('webpage').create(),
    system = require('system'),
    address, output, size, pageWidth, pageHeight;

page.settings.resourceTimeout = 70000; // 70 seconds
page.onResourceTimeout = function(e) {
    console.log(e.errorCode);   // it'll probably be 408
    console.log(e.errorString); // it'll probably be 'Network timeout on resource'
    console.log(e.url);         // the url whose request timed out
    phantom.exit(1);
};

if (system.args.length < 3 || system.args.length > 5) {
    console.log('Usage: rasterize.js URL filename [paperwidth*paperheight|paperformat] [zoom]');
    console.log('  paper (pdf output) examples: "5in*7.5in", "10cm*20cm", "A4", "Letter"');
    console.log('  image (png/jpg output) examples: "1920px" entire page, window width 1920px');
    console.log('                                   "800px*600px" window, clipped to 800x600');
    phantom.exit(1);
} else {
    address = system.args[1];
    output = system.args[2];
    page.viewportSize = { width: 600, height: 600 };
    if (system.args.length > 3 && system.args[2].substr(-4) === ".pdf") {
        size = system.args[3].split('*');
        page.paperSize = size.length === 2 ? { width: size[0], height: size[1], margin: '0px' }
                                           : { format: system.args[3], orientation: 'portrait', margin: '1cm' };
    } else if (system.args.length > 3 && system.args[3].substr(-2) === "px") {
        size = system.args[3].split('*');
        if (size.length === 2) {
            pageWidth = parseInt(size[0], 10);
            pageHeight = parseInt(size[1], 10);
            page.viewportSize = { width: pageWidth, height: pageHeight };
            page.clipRect = { top: 0, left: 0, width: pageWidth, height: pageHeight };
        } else {
            console.log("size:", system.args[3]);
            pageWidth = parseInt(system.args[3], 10);
            pageHeight = parseInt(pageWidth * 3/4, 10); // it's as good an assumption as any
            console.log ("pageHeight:",pageHeight);
            page.viewportSize = { width: pageWidth, height: pageHeight };
        }
    }
    if (system.args.length > 4) {
        page.zoomFactor = system.args[4];
    }
    page.open(address, function (status) {
        if (status !== 'success') {
            console.log('Unable to load the address!');
            phantom.exit(1);
        } else {
            window.setTimeout(function () {
                page.render(output);
                phantom.exit();
            }, 200);
        }
    });
}`
)

func ScreenshotPhantomJS(urlStr, path string, offscreen bool) error {
	fmt.Println(urlStr)
	jsFile := grand.RandomString(10) + ".js"
	if err := gfs.StringToFile(rasterizeJS, jsFile); err != nil {
		return err
	}
	defer os.Remove(jsFile)
	app := "phantomjs"
	if offscreen {
		app = "QT_QPA_PLATFORM=offscreen phantomjs"
	}
	return gcmd.ExecWaitPrintScreen("bash", "-c", fmt.Sprintf("%s %s %s %s", app, jsFile, urlStr, path))
}
