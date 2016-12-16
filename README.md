# remote-poller
[![GoDoc](https://godoc.org/github.com/JoelW-S/remote-poller?status.svg)](https://godoc.org/github.com/JoelW-S/remote-poller)

## Usage

This library is used to listen on directories typically remote ones for events such as:

* File added
* File modified
* File deleted

This Library handles the polling and notifying the user of an event. 

The client needs to provide a means of listing files at the remote location. This can be done by 
creating an implementation of [PolledDirectory](https://github.com/JoelW-S/remote-poller/blob/master/cycler.go#L15).

Instead of providing multiple adapters/implementations to cater for library specific types, 
this library provides [Element](https://github.com/JoelW-S/remote-poller/blob/master/element.go).

remote-poller comes with [FileElement](https://github.com/JoelW-S/remote-poller/blob/master/element.go#L26).
This type can be used to convert golang's os.FileInfo interface to Element 
and be used to create simple filesystem polling examples.

[Examples](https://github.com/JoelW-S/remote-poller/blob/master/examples)

