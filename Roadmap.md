# cucu: TUI HTTP Client in C++ - One Month Roadmap

## 🗓️ Overview Milestones

| Week | Goals | Date Started |
|------|-------|------|
| 1 | Set up project skeleton: TUI, build system, UI layout | 2025-05-26 |
| 2 | Implement HTTP request logic: GET/POST with headers + body | - |
| 3 | Hook up UI to request engine, build JSON/text views | - |
| 4 | Polish UI, handle edge cases, package MVP build | - |

---

### Week 1: Project Setup + UI Skeleton

**Goals**
- Project structure using `CMake`
- Build FTXUI-based TUI
- Define app layout (sidebar, method input, URL input, send button, response view)

**Tasks**
- [x] Set up `CMakeLists.txt` with `FTXUI` (via `FetchContent`)
- [x] Create a `main.cpp` with:
  - Input box for URL
  - Dropdown or radio box for method (GET, POST)
  - Text box for headers and body
  - Button to "Send"
  - Scrollable area for response
- [x] Add `nlohmann/json` (for JSON formatting)
- [ ] Set up basic state model (URL, method, body, response)

**Deliverable**
A navigable TUI that accepts user input, no HTTP requests yet.

```
╭───────────────────────────── Context Bar ────────────────────────────╮
│ [GET ▼]  https://api.example.com/v1/users                            │
├──────────────────────────── Request / Response ──────────────────────┤
│                                                                      │
│  (Depending on mode: headers or response view appears here)          │
│                                                                      │
├────────────────────────────── Footer Bar ────────────────────────────┤
│ [i] Insert  [x] Delete  [j/k] Navigate  [:] Command  [/] Filter      │
╰──────────────────────────────────────────────────────────────────────╯

```

---

### Week 2: HTTP Layer (Backend)

**Goals**
- Make HTTP requests with headers/body
- Handle JSON and plain text responses
- Log errors cleanly

**Tasks**
- [ ] Integrate `cpr` (or `libcurl`)
- [ ] Implement `send_request(method, url, headers, body)`
- [ ] Return parsed response or error
- [ ] Unit-test logic independently of UI

**Deliverable**
A working request engine that handles GET/POST requests.

---

### Week 3: UI ↔ Backend Integration

**Goals**
- Connect UI inputs to HTTP logic
- Show status while sending
- Display response in UI

**Tasks**
- [ ] Wire "Send" button to `send_request(...)`
- [ ] Display response body, status, headers
- [ ] Highlight JSON if possible
- [ ] Show error messages in UI

**Deliverable**
Working input → HTTP → response pipeline.

---

### Week 4: Polish and MVP Completion

**Goals**
- Input validation
- User-friendly interface
- Easy build/run

**Tasks**
- [ ] Validate URL/method
- [ ] Add scroll/resize support
- [ ] Handle long responses
- [ ] Optional: Save/load recent requests
- [ ] Write usage instructions

**Deliverable**
MVP binary: sends HTTP requests and shows responses.

---

cucu/
├── CMakeLists.txt
├── main.cpp
├── ui/
│   └── layout.cpp, layout.hpp
├── net/
│   └── request.cpp, request.hpp
├── models/
│   └── app_state.hpp
└── third_party/
    └── (cpr, ftxui, json via FetchContent)

---
