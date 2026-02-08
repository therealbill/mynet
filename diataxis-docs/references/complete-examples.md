# Complete Documentation Examples

Read this file to see finished examples of each Diataxis doc type. Use these as quality benchmarks when writing.

## How-to Guide Example: Email Notifications

```markdown
---
title: "How to Set Up Email Notifications"
summary: "Configure email alerts for system events"
doc_type: how-to
prerequisites: ["Application installed", "SMTP credentials"]
est_time: "10 min"
roles: ["developer", "operator"]
stability: stable
---

# How to Set Up Email Notifications

**Goal**: Configure your application to send email notifications for system events.

**Time**: Approximately 10 minutes

## Prerequisites

- Application installed and running — see [Installation tutorial](../tutorials/00-onboarding.md)
- SMTP server credentials (host, port, username, password)

## Steps

### 1. Create email configuration file

Create `config/email.json` with your SMTP settings:

\`\`\`json
{
  "smtp": {
    "host": "smtp.example.com",
    "port": 587,
    "secure": true,
    "auth": {
      "user": "your-username",
      "pass": "your-password"
    }
  },
  "from": "notifications@example.com"
}
\`\`\`

### 2. Enable email notifications

Edit `config/app.json` and enable notifications:

\`\`\`json
{
  "notifications": {
    "enabled": true,
    "channels": ["email"],
    "events": ["error", "warning"]
  }
}
\`\`\`

### 3. Restart the application

\`\`\`bash
npm restart
\`\`\`

Expected output:

\`\`\`
Server restarted
Email notifications: enabled
\`\`\`

## Verify it works

Trigger a test notification:

\`\`\`bash
npm run test:notification
\`\`\`

You should receive an email with subject "Test Notification". Check your inbox (and spam folder).

## Troubleshooting

### Problem: No emails received

**Symptom**: Application runs but emails don't arrive
**Cause**: Incorrect SMTP credentials or firewall blocking
**Solution**:

1. Check credentials in `config/email.json`
2. Verify port 587 is not blocked
3. Check application logs: `tail -f logs/email.log`

### Problem: "Authentication failed" error

**Symptom**: `Error: SMTP authentication failed`
**Cause**: Wrong username or password
**Solution**:

- Verify credentials with your email provider
- Some providers require app-specific passwords
- Check for spaces or special characters in password

## See also

- [How to customize notification templates](./customize-notifications.md)
- [Email API reference](../reference/email-api.md)
- [Understanding the event system](../explanation/events.md)
```

**Why this works:** Single goal, numbered steps, complete code, verification, troubleshooting, cross-links. Under 1800 words.

---

## Explanation Example: Event System

```markdown
---
title: "Understanding the Event System"
summary: "How event-driven architecture works in the platform"
doc_type: explanation
prerequisites: []
est_time: "15 min"
roles: ["developer", "architect"]
stability: stable
---

# Understanding the Event System

Our platform uses event-driven architecture to decouple components
and enable extensibility. This explanation covers how the event
system works, why we designed it this way, and the trade-offs involved.

## The problem

Traditional request-response systems create tight coupling between
components. When Service A needs to notify Service B of a change,
it must know about Service B explicitly:

\`\`\`javascript
// Tight coupling — every new reaction requires modifying this code
userService.createUser(data);
emailService.sendWelcome(data);
analyticsService.trackSignup(data);
\`\`\`

This becomes problematic as the system grows:

- Adding new reactions requires modifying existing code
- Services can't be developed independently
- Testing requires mocking many dependencies
- Failures cascade between services

## How the event system works

Instead, we use a publish-subscribe pattern:

\`\`\`
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│ User Service │────>│  Event Bus   │────>│Email Service │
└──────────────┘     │              │     └──────────────┘
                     │              │────>┌──────────────┐
                     │              │     │  Analytics   │
                     │              │     └──────────────┘
                     │              │────>┌──────────────┐
                     │              │     │  Audit Log   │
                     └──────────────┘     └──────────────┘
\`\`\`

**Publishers** emit events without knowing who will handle them:

\`\`\`javascript
eventBus.publish('user.created', { userId: 123, email: 'user@example.com' });
\`\`\`

**Subscribers** register interest in event types:

\`\`\`javascript
eventBus.subscribe('user.created', async (event) => {
  await sendWelcomeEmail(event.email);
});
\`\`\`

**Event Bus** routes events from publishers to subscribers, maintaining
registrations, delivering asynchronously, handling retries, and
guaranteeing at-least-once delivery.

## Why this design?

We chose event-driven architecture for three key reasons:

### 1. Decoupling

Services don't need to know about each other. The User Service
just publishes "user.created" — it doesn't care who reacts.

### 2. Extensibility

Adding new behavior is trivial — just add a new subscriber.
No changes to existing code.

### 3. Resilience

If the Email Service is down, the User Service keeps working.
Events are queued and retried automatically.

## Trade-offs

| Benefit | Trade-off |
|---------|-----------|
| Loose coupling | Harder to trace execution flow |
| Easy to extend | More moving parts to operate |
| Resilient to failures | Eventual consistency (not immediate) |
| Independent services | Debugging requires correlation IDs |

### When it's worth it

Event-driven architecture shines when you need to extend behavior
without modifying code, services are developed by different teams,
or resilience to partial failures matters.

### When it's not

For simple, synchronous operations, request-response is simpler:
direct authentication checks, read operations with no side effects,
or time-sensitive operations requiring immediate response.

## Common misconceptions

**"Events are always faster than requests"**
False. Events add latency from async delivery. Use events for
decoupling, not performance.

**"Events guarantee exactly-once processing"**
False. We provide at-least-once. Handlers must be idempotent.

**"All operations should use events"**
False. Synchronous request-response is simpler when you don't
need decoupling or async processing.

## Related

- **Do it:** [How to publish events](../how-to/publish-events.md)
- **Do it:** [How to create subscribers](../how-to/create-subscribers.md)
- **Look it up:** [Event catalog](../reference/events.md)
- **Related concept:** [Understanding service architecture](./architecture.md)
```

**Why this works:** Opens with the problem, builds a mental model, explains the design rationale, discusses trade-offs honestly, addresses misconceptions, links to practical content. No step-by-step instructions.

---

## Tutorial Example: Key Patterns

A tutorial doesn't need a complete example here because the template in `tutorial-template.md` is already structured as a fill-in-the-blanks guide. Instead, here are the key quality patterns to follow:

### Strong opening

```markdown
# Build Your First REST API

In this tutorial, you'll build a working REST API for a task manager
in about 90 minutes. By the end, you'll have an API that can create,
read, update, and delete tasks — running on your local machine.
```

Concrete outcome, specific time, clear scope.

### Good checkpoint

```markdown
### Checkpoint

Test your endpoint:

\`\`\`bash
curl http://localhost:3000/tasks
\`\`\`

Expected output:

\`\`\`json
{"tasks":[]}
\`\`\`

If you see errors, check that your server is running (Step 3) and
that the route is registered correctly.
```

Verification command, expected output, troubleshooting hint.

### Effective "what just happened" box

```markdown
> **What just happened?**
>
> That command created a database migration file. We'll run it in
> Step 8 to create the tasks table. Migrations let you version
> your database schema — but don't worry about that now.
```

Brief context, forward reference, explicit "don't worry about the details."

---

## Reference Example: Key Patterns

Reference doesn't benefit from a full example document because the template in `reference-template.md` covers the structure. Here are the key quality patterns:

### No-advice rewrite

Before (contains advice):
```markdown
# sendEmail

We recommend using the async version for better performance.
You should always validate email addresses before calling this.
```

After (pure specification):
```markdown
# sendEmail

Sends an email message synchronously. Blocks until completion or timeout.

An async version is available as `sendEmailAsync`.
Input validation is the caller's responsibility.
```

### Syntax-only examples

Before (shows workflow):
```markdown
## Example: Building a complete authentication flow

First, create a user, then verify their email, then log them in...
```

After (shows syntax):
```markdown
## Example

\`\`\`javascript
authenticate(username, password)
// Returns: { token: "abc123", expiresIn: 3600 }
\`\`\`
```
