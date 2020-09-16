lru
===

[![Build Status](https://github.com/hdfchain/hdfd/workflows/Build%20and%20Test/badge.svg)](https://github.com/hdfchain/hdfd/actions)
[![ISC License](https://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hdfchain/hdfd/lru)

Package lru implements a generic least-recently-used cache with near O(1) perf.

## LRU Cache

A least-recently-used (LRU) cache is a cache that holds a limited number of
items with an eviction policy such that when the capacity of the cache is
exceeded, the least-recently-used item is automatically removed when inserting a
new item.  The meaning of used in this implementation is either accessing the
item via a lookup or adding the item into the cache, including when the item
already exists.

## External Use

This package has intentionally been designed so it can be used as a standalone
package for any projects needing to make use of a well-tested and concurrent
safe least-recently-used cache with near O(1) performance characteristics for
lookups, inserts, and deletions.

## Installation and Updating

```bash
$ go get -u github.com/hdfchain/hdfd/lru
```

## Examples

* [Basic Usage](https://godoc.org/github.com/hdfchain/hdfd/lru#example-package--BasicUsage)
  Demonstrates creating a new cache instance, inserting items into the cache,
  causing an eviction of the least-recently-used item, and removing an item.

## License

Package lru is licensed under the [copyfree](http://copyfree.org) ISC License.
