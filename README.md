# Static

Static is the caching layer between Litenin and Nimbus.

## Motivation

Nimbus is backed by a relational database and is tailored to handle single feed queries. This provides great flexibility and simplicity, but lacks speed and effeciency.

## Goal

The goal of Static is to provide a caching layer between Litenin and Nimbus in order to minimize the load on Nimbus.

- Service batch requests of feeds from Litenin through a simple API (just POST a serialized array of feed URL's).
- Intelligent caching of resources from Nimbus based on expiration and frequency of requests.
- Attempt to have the best possible cache coverage of queries from Litenin given the hardware limitations.
