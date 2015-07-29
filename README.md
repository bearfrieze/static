# Static

**Static has been retired and is no longer maintained. Batch feed requests of the type that Static used to handle are now handled by Nimbus that uses Redis for caching.**

Static is the caching layer between [Litenin](https://github.com/bearfrieze/litenin) and [Nimbus](https://github.com/bearfrieze/nimbus).

## Motivation

Nimbus is backed by a relational database and is tailored to handle single feed queries. This provides great flexibility and simplicity, but lacks speed and effeciency. Static was planned along with Nimbus, but it was decided not to include the functionality into Nimbus for the following reasons:

- Simplicity and modularization: Nimbus and Static can both focus on doing one thing really well.
- Distribution and scalability: Availablity of Nimbus resources can easily be increased by launching additional instances of static. Latency can be kept at a minimum by having instances of Static runnning in relevant regions.

## Goal

The goal of Static is to provide a caching layer between Litenin and Nimbus in order to minimize the load on Nimbus.

- Service batch requests of feeds from Litenin through a simple API (just POST a serialized array of feed URL's).
- Intelligent caching of resources from Nimbus based on expiration and frequency of requests.
- Attempt to have the best possible cache coverage of queries from Litenin given the hardware limitations.
