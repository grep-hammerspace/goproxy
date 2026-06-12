# Forward Proxy Build Plan

## Phase 1 — Accept a single connection

- Start a TCP listener on a port (e.g. 8080)
- Accept one incoming connection
- Read the raw bytes and print them to stdout

**Test:** `curl -x localhost:8080 http://example.com` and see the raw HTTP request printed

---

## Phase 2 — Parse the destination

- From the bytes you printed in phase 1, write code to extract the host and port from the request line
- Print the extracted host/port to confirm it's correct

**Test:** Same curl command, confirm you're printing `example.com:80`

---

## Phase 3 — Connect to the destination

- Open a TCP connection to the extracted host/port
- Don't forward anything yet, just confirm the connection succeeds

**Test:** Confirm no connection error, then close both connections

---

## Phase 4 — Forward the request

- Send the bytes you read from the client to the destination connection
- Read the response from the destination and send it back to the client

**Test:** `curl -x localhost:8080 http://example.com` should return actual HTML
---

## Phase 5 — Handle multiple connections

- Right now it handles one connection then stops
- Wrap the accept loop so it handles connections concurrently

**Test:** Run two curl commands simultaneously, both should succeed
Done up to here.
---

## Phase 6 — Handle persistent connections

- A single TCP connection can carry multiple HTTP requests (keep-alive)
- After forwarding one request, keep the connections open and handle the next request

**Test:** `curl -x localhost:8080 --keepalive http://example.com http://example.com/about` — both requests should succeed over one connection