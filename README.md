# NOTE: This software is under construction and doesn't do much that this README says it does
# Transcodebot 🎞🤖
While there are plenty of solutions around for distributing video transcoding (or large workloads in general), they are designed with linux as the only target platform for distributing too. And while that is fine for corprate entities, most people don't have access to large numbers of computers that they can run a nonstandard OS on and dedicate wholly to transcoding. Many people do, however, have a few computers that generally run very light loads on non-Linux operating systems. Transcodebot is designed to leverage these kinds of machines towards batch transcoding tasks with minimal interuption in their primary use as web browsing machines.

It also provides sane defaults for video transcoding to minimize the amount of research needed to transcode video.

### Quick setup
[Install go](https://golang.org/doc/install)

    go get -u github.com/yourfin/transcodebot/...
    transcodebot build

## Usage
### `build`
Build the self-contained client binaries.

### `watch`
Watch a folder for new files to transcode, and push them out to be transcoded as they come in.
Also runs a web server to download clients from.

### `one-shot`
Like watch, but only the files passed in on the command line are transcoded

## Design
Transcodebot is designed for client machines that have generally have something better to do.

### Why go (and not java)?
 - Single binary client distribution:
It's much easier to only have a single program to run than having to install a java runtime on all clients
 - Go's network stack
 - Concurrency primitives

## TODO:
### General
 - Add these as issues
### Client
 - Run only during some hours
 - Disable backrounding for dedicated hardware
 - OSX
 - Download ffmpeg from external source
 - Default client file
 - Show documentation on startup for clients: skratchdot/opengolang
 - Pause when listed processes seen
   - Widows: DebugActiveProcess
   - Unix: SIGSTOP SIGCONT
### ffmpeg
 - Find reliable downloads for statically compiled ffmpeg versions
 - Clean up ffmpeg calls
 - Provide reasonable defaults for things other than VP9
   - h.264
   - h.265
   - audio
 - Point towards documentation of how ffmpeg works
### Server
 - Figure out where to put settings if running as root
 - Verify static compiles of clients
 - Docker!

## Wishful thinking:
 - Have some system for detecting if clients have someone being a user on them
 - Learn client by client ETA's for files
 - Pausing ffmpeg
 - GUI's
 - Windows server
   - Figure out when superuser
 - Prettify server download page
 - Run only on wall power
   - github.com/distatus/battery

## Licence
MIT
