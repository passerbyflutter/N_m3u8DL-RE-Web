## Unpublish

### Features

- Add real-time progress updates via SSE, #20
- Add function for deleting deleted tasks, #16

### Fixes

- Fix race condition in command output processing, #22
  - Implement proper synchronization for stdout/stderr streams
  - Add comprehensive error detection and handling
  - Improve task status management for failed downloads

## 0.0.1

### Features

- Initialize Project
