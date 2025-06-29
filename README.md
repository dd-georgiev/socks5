# Overview
This is very primitive implementation of the SOCKS5 protocol. It contains:
1) A client capable of executing CONNECT, BIND and UDP_ASSOCIATE commands
2) A server capable of serving CONNECT, BIND and UDP_ASSOCIATE commands

It supports only `No Auth` method and `IPv4`.

The implementation is based on [RFC-1928](https://datatracker.ietf.org/doc/html/rfc1928) and [Dante](https://www.inet.no/dante/). 
In the docs folder there is a series of [labs](https://github.com/dd-georgiev/socks5/tree/main/docs/labs/index.md) which contain the rough code, implemented piece by piece as I was writing it without any refactoring. 
The first lab sets the foundations of how **as per my understanding** the socks5 protocol functions. 


The server lacks some fundamental features such as:
1) Authentication
2) Timeouts(i.e. when client is inactive for X amount of time)
3) Proper error handling for edge cases

The client is very basic, it lacks proper error handling for edge cases as well