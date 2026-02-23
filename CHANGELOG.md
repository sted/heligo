# Changelog

## 0.3.0

### Security
- ReadJSON: limit body to 1MB to prevent DoS
- Param count bounds check (MAXPARAMS)
- WriteJSON: marshal before writing header to avoid partial responses on error
- AdaptMiddleware: propagate wrapped ResponseWriter and Request
- Group: fix middleware slice aliasing between parent and sub-groups

### Added
- `Recover(onPanic)` middleware for panic recovery
- `CleanPaths()` middleware for URL sanitization
- `HasPath(method, path)` for 405 Method Not Allowed support
- `TrailingSlash` option: register both `/path` and `/path/` at zero runtime cost
- `ErrorHandler` now receives `http.ResponseWriter` and `*http.Request`
