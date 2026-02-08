---
title: "Getting Started with Language Agents"
description: "Learn the workflow of using language-specific agents to review and improve code by working through hands-on examples with Go and C++ agents"
weight: 1
---

# Getting Started with Language Agents

In this tutorial, you will learn the core workflow for using the programming-languages plugin to review, refine, and improve code. You will write code in two languages, invoke the appropriate language agent, and observe how each agent applies language-specific idioms and best practices.

By the end of this tutorial, you will understand how to:

- Identify which agent handles which language
- Ask an agent to review existing code
- Interpret the agent's changes and explanations
- Use the review-and-refine cycle to iteratively improve code quality

## Prerequisites

- The programming-languages plugin installed in your Claude Code environment
- A working directory where you can create test files
- Basic familiarity with at least one of the supported languages (C++, Go, JavaScript, TypeScript, or Zsh)

## Step 1: Write Some Go Code

Start by creating a simple Go HTTP handler that has room for improvement. Create a file called `handler.go` with the following content:

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io/ioutil"
)

func handleUser(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        users := getUsers()
        data, err := json.Marshal(users)
        if err != nil {
            w.WriteHeader(500)
            w.Write([]byte("error"))
            return
        }
        w.Header().Set("Content-Type", "application/json")
        w.Write(data)
    } else if r.Method == "POST" {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            w.WriteHeader(400)
            w.Write([]byte("bad request"))
            return
        }
        var user User
        err = json.Unmarshal(body, &user)
        if err != nil {
            w.WriteHeader(400)
            w.Write([]byte("invalid json"))
            return
        }
        saveUser(user)
        w.WriteHeader(201)
        fmt.Fprintf(w, "created")
    } else {
        w.WriteHeader(405)
        w.Write([]byte("method not allowed"))
    }
}
```

This code works, but it uses several non-idiomatic Go patterns. You will use the go-simplifier agent to clean it up.

**Checkpoint:** You should have a `handler.go` file with approximately 40 lines of working but verbose Go code.

## Step 2: Ask go-simplifier to Review the Code

Now ask Claude Code to simplify the Go code. The go-simplifier agent activates automatically when you ask for Go-specific improvements. Try a prompt like:

```
Clean up the Go code in handler.go and make it more idiomatic
```

The go-simplifier agent will read the file, identify non-idiomatic patterns, and refine the code. Watch for these typical improvements:

- **`ioutil.ReadAll` replaced with `io.ReadAll`** -- `ioutil` has been deprecated since Go 1.16
- **`if/else if` chain replaced with a `switch` statement** -- `switch` is the idiomatic way to dispatch on method in Go
- **`json.NewEncoder` and `json.NewDecoder` used instead of Marshal/Unmarshal** -- streaming is more efficient for HTTP handlers
- **Early returns to reduce nesting** -- the agent favors flat control flow
- **`http.StatusOK` constants instead of magic numbers** -- named constants are clearer than `500`, `400`, `201`

**Checkpoint:** After the agent finishes, your `handler.go` should be shorter, use a switch statement, and contain no deprecated APIs. The agent will summarize every change it made and explain why.

## Step 3: Understand the Agent's Reasoning

The go-simplifier agent always explains its changes in a summary after applying them. Review the summary carefully. Each change will include:

- **What changed** -- the specific code transformation
- **Why** -- the idiomatic Go convention or best practice behind it

This is not a generic code formatter. The go-simplifier uses the opus model (the strongest reasoning model) specifically because simplification requires judgment about which patterns are genuinely clearer, not just different. It applies Effective Go guidelines and the Go Code Review Comments standard.

If you disagree with any change, you can ask the agent to revert specific modifications while keeping others.

**Checkpoint:** You should be able to explain why each change the agent made is an improvement.

## Step 4: Write Some C++ Code

Now try a different language to see how agent specialization works. Create a file called `resource.cpp`:

```cpp
#include <iostream>
#include <cstring>

class Buffer {
    char* data;
    size_t size;
public:
    Buffer(size_t s) : size(s) {
        data = new char[size];
        memset(data, 0, size);
    }

    ~Buffer() {
        delete[] data;
    }

    void copy_from(const Buffer& other) {
        delete[] data;
        size = other.size;
        data = new char[size];
        memcpy(data, other.data, size);
    }

    char* get() { return data; }
};

Buffer* create_buffer(size_t size) {
    Buffer* buf = new Buffer(size);
    return buf;
}

void process(Buffer* buf) {
    // ... use buf ...
    std::cout << "Processing " << buf->get() << std::endl;
}

int main() {
    Buffer* buf = create_buffer(1024);
    process(buf);
    delete buf;
    return 0;
}
```

This code has several safety issues: raw `new`/`delete`, a missing copy constructor, no move semantics, and raw pointer ownership.

## Step 5: Ask cpp-pro to Modernize the Code

Ask Claude Code to modernize the C++ code:

```
Modernize this C++ code to use smart pointers and RAII
```

The cpp-pro agent will identify and fix issues specific to modern C++:

- **Raw `new`/`delete` replaced with `std::unique_ptr`** -- the buffer factory function returns `unique_ptr` instead of a raw pointer
- **Rule of Five applied** -- copy constructor, copy assignment, move constructor, and move assignment added or defaulted
- **`std::vector<char>` instead of manual allocation** -- eliminates the need for a custom destructor entirely
- **`memset`/`memcpy` replaced with STL algorithms or constructors** -- using the initializing constructor of `vector` handles zeroing

**Checkpoint:** The modernized code should have zero raw `new` or `delete` calls. The agent uses C++17 as its baseline standard and will note if any changes require C++20 features.

## Step 6: Compare Agent Behaviors

Notice the differences between the two agents:

| Aspect | go-simplifier | cpp-pro |
|--------|---------------|---------|
| **Focus** | Clarity and idiomaticity | Safety and performance |
| **Model** | Opus (stronger reasoning) | Sonnet (fast iteration) |
| **Verification** | Runs `go vet` | Runs build with sanitizers |
| **Philosophy** | Simpler is better | Modern standards compliance |

Both agents share a core principle: they change how code is expressed, not what it does. Neither agent will alter the behavior of your code unless you explicitly ask.

## Step 7: Try the Review Cycle Again

The review-and-refine workflow is iterative. After the first pass, you can ask follow-up questions:

- "Are there any remaining issues in this file?"
- "Can you also add error handling for the edge case where size is zero?"
- "What would this look like with C++20 concepts?"

Each follow-up narrows the scope and deepens the improvements. The agent remembers context from the current session and builds on its previous analysis.

**Checkpoint:** You have completed at least two rounds of review on one file and can see how the iterative cycle produces progressively better code.

## What You Learned

In this tutorial you:

- Used go-simplifier to transform verbose Go into idiomatic, clean code
- Used cpp-pro to modernize unsafe C++ into RAII-based modern C++
- Observed how each agent applies language-specific knowledge that a general-purpose reviewer would miss
- Practiced the review-and-refine cycle that makes these agents most effective

## Next Steps

- Read the [Modernize Legacy Code](../../howto/modernize-legacy-code/) guide for specific legacy modernization recipes
- See the [Agents](../../reference/agents/) for the complete specification of all five agents
- Learn [Language Specialization](../../explanation/language-specialization/) to understand what each agent knows that a generalist does not
