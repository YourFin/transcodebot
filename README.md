# NOTE: This software is under construction and doesn't do much that this README says it does
# Transcodebot
Transcodebot is a "batteries included" system for distributing large video transcoding tasks to computers that have other jobs, like being workstations or browsing machines.

### Quick setup
[Install go](INSERT LINK)
    go install github.com/yourfin/transcodebot

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
