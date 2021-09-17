# Why?

Need a file watcher component for Exo, but there are some problems:
1. development instance of exo (dexo) restarts frequently (ie it uses a file
   watcher!). need an isolated, stateful subcomponent that isn't frequently
   restarted.
2. exo and dexo are often monitoring the same files; multiple
   components/workspaces may be monitoring overlapping filesets too. This is
   inefficient. Maybe not a problem for 2 or 3 watchers on a directory, but
   efficient multiplexing of the underlying OS resources enables us to freely
   create many "watchers" with varying rules.
3. we want to make the above efficencies available to non-core exo components
   coded in a variety of languages.
4. It sounds like fun?

Goals:
1. Use the most resource efficient monitoring mechanism on each platform,
   especially for recursive watches.  (FSEvents on MacOS)
2. HTTP-based interface: retained events for pull, Server-Sent Events (SSE) for
   push. Retained events pull makes for reliable consumers & push makes for
   fast response. Specified with an OpenAPI schema.
3. sensible defaults: easily watch directories recursively, ignore hidden
   files, coallese/debounce bulk changes, etc.

Non-goals:
1. dynamic rulesets: creating/deleting watchers should be cheap enough that
   clients can accomplish this via wholesale replacement.

# Example Usage

```
> http localhost:3000/watchers path='/path/to/watch'
< 201 Created
< Location: abc123

> http localhost:3000/watchers/abc123
< { "path": "/path/to/watch" }

> http localhost:3000/watchers/abc123/events
< { items: [{"id": "001", ... }], ... }

> http localhost:3000/watchers/abc123/events?after=001
< { items: [{"id": "002", ...

> http --stream localhost:3000/watchers/abc123/events?after=001
< data: {"id": "002", ... 
< data: {"id": "003", ... 

> http delete localhost:3000/watchers/abc123
```

# Operations

- Create watcher
- Describe watchers, filterable by ids and tags. (TODO: Paginate?)
- Delete watchers, selected by ids or tags.
- Get events, potentially since some timestamp, with or without streaming.

# Formats

## Watcher Description

JSON with the following keys:

- `path: string`: root of the tree to watch.
- `tags: Record<string, string>`: arbitrary key value pairs for later lookup.

Hidden files are ignored by default.

### Planned Extensions

**Sequence of rules:**

Instead of `path`, more complex rule sets.

```json
{
  "rules": [
    ["include-path", "/path/to/watch"],
    ["exclude-hidden"],
    ["exclude-glob", "/path/to/watch/**.json"]
  ]
}
```

**Coallessing Interval/Resolution/Duration**

**Timestamps, IDs for least and most recent events**

## Events

JSON object with the following keys:

- `id`: sequential event id string. Compare lexographically.
- `action`: some string describing the kind of change. created, updated,
  deleted, plus some metadata events, things like mounts, etc. Details TBD.
- `path`: file path affected.

# Security

File paths are potentially privacy sensitive information. Events should be
filtered by what the user has access to. Need Kerberos or somethign like that.

# Platform

## MacOS

FSEvents already works correctly! It provides persistent events, replay,
etc. Adapting will be very easy.

https://github.com/fsnotify/fsevents provides decent bindings.

Any abstraction layer would be a step backwards for MacOS.

Event IDs are int64, which is too big for a float64 in JSON. Must pad with zeros to preserve lexographical ordering.

Retention policy is up to the operating system.

## Other Platforms

Potential abstraction layer:
- https://github.com/emcrisostomo/fswatch
- http://emcrisostomo.github.io/fswatch/doc/1.16.0/libfswatch.html/

Sadly, nothing in Go-land seems as robust as libfswatch.

If we implement our own event database, we'll need a retention policy.
