package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net"
    "os"
    "strings"
    "time"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()
    addr := conn.RemoteAddr().String()
    log.Printf("[%s] Connected at %s", addr, time.Now().Format(time.RFC3339))

    // Per‑client log file named by IP
    ip := strings.Split(addr, ":")[0]
    f, err := os.OpenFile(ip+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Could not open log for %s: %v", addr, err)
        return
    }
    defer f.Close()
    clientLog := log.New(f, "", log.LstdFlags)

    reader := bufio.NewReader(conn)
    for {
        // Enforce 30s inactivity timeout
        conn.SetReadDeadline(time.Now().Add(30 * time.Second))
        line, err := reader.ReadString('\n')
        if err != nil {
            if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                conn.Write([]byte("Disconnected due to inactivity\n"))
            }
            log.Printf("[%s] Disconnected at %s", addr, time.Now().Format(time.RFC3339))
            return
        }

        // Trim whitespace
        msg := strings.TrimSpace(line)
        // Truncate overly long messages
        if len(msg) > 1024 {
            msg = msg[:1024]
            conn.Write([]byte("Message too long, truncated\n"))
        }

        // Log the cleaned message
        clientLog.Println(msg)

        // Personality / commands
        var response string
        switch {
        case msg == "":
            response = "Say something...\n"
        case strings.EqualFold(msg, "hello"):
            response = "Hi there!\n"
        case strings.EqualFold(msg, "bye"):
            response = "Goodbye!\n"
            conn.Write([]byte(response))
            log.Printf("[%s] Disconnected at %s", addr, time.Now().Format(time.RFC3339))
            return
        case strings.HasPrefix(msg, "/"):
            response = handleCommand(msg)
            // If /quit, close after sending
            if strings.HasPrefix(strings.ToLower(msg), "/quit") {
                conn.Write([]byte(response))
                log.Printf("[%s] Disconnected at %s", addr, time.Now().Format(time.RFC3339))
                return
            }
        default:
            // Standard echo
            response = msg + "\n"
        }

        // Send back
        if _, err := conn.Write([]byte(response)); err != nil {
            log.Printf("Write error to %s: %v", addr, err)
            return
        }
    }
}

func handleCommand(msg string) string {
    parts := strings.Fields(msg)
    cmd := strings.ToLower(parts[0])
    switch cmd {
    case "/time":
        return time.Now().Format(time.RFC3339) + "\n"
    case "/quit":
        return "Goodbye!\n"
    case "/echo":
        if len(parts) > 1 {
            return strings.Join(parts[1:], " ") + "\n"
        }
        return "\n"
    default:
        // Unknown commands just get echoed back
        return msg + "\n"
    }
}

// Command‑line flag for the port (default 4000)
var port = flag.Int("port", 4000, "Port to listen on")

func main() {
    flag.Parse()
    addr := fmt.Sprintf(":%d", *port)

    listener, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("Error starting server: %v", err)
    }
    defer listener.Close()
    log.Printf("Server listening on %s", addr)

    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Accept error: %v", err)
            continue
        }
        go handleConnection(conn)
    }
}

