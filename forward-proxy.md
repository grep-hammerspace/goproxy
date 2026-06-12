# HTTPS Support Build Plan

## Phase 7 — Detect CONNECT requests

- When a request comes in, check if the method is `CONNECT`
- If it is, pass it to a new `handleTunnel` function
- If it isn't, handle it as before with the existing plain HTTP logic
- You don't need to do anything in `handleTunnel` yet, just log that you received a CONNECT request

**Test:** `curl -v -x localhost:8080 https://example.com` — you should see the CONNECT request logged

---

## Phase 8 — Establish the tunnel

- In `handleTunnel`, extract the host and port from the CONNECT request (e.g. `example.com:443`)
- Open a TCP connection to that host and port
- Send `HTTP/1.1 200 Connection Established\r\n\r\n` back to the client
- This tells the client the tunnel is open and it can begin its TLS handshake

**Test:** Same curl command — curl should no longer hang after connecting, though you won't get a full response yet

---

## Phase 9 — Copy bytes bidirectionally

- Once the tunnel is established, copy bytes in both directions simultaneously:
    - Client → target
    - Target → client
- These must run concurrently in two separate goroutines, since either side can send at any time
- Use `io.Copy` for each direction
- Wait for both goroutines to finish before closing connections

**Test:** `curl -v -x localhost:8080 https://example.com` should return real HTML over HTTPS

---

## Phase 10 — Handle tunnel teardown cleanly

- When either side closes the connection, the other side should be closed too
- If the client disconnects, close the target connection
- If the target disconnects, close the client connection
- Without this, one goroutine may hang waiting for bytes that will never arrive

**Test:** `curl -v -x localhost:8080 https://example.com https://example.com` — two HTTPS requests, both succeed, no hanging goroutines after completion