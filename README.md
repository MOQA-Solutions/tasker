# Task Manager

This repository contains **partial and illustrative code samples** used to showcase **Tasker**.

> ⚠️ **Important Notice**  
> This is **not the complete or production version** of Tasker.  
> The code here is provided **for demonstration purposes only** to showcase how Tasker distributes and balances tasks across the worker pool using priority-based channels.

---

## 🧩 Purpose

The main goal of this repository is to demonstrate:

- How Tasker interacts with the **Worker Pool**.  
- How to Select a **Specific Worker** to handle a task.
- **One Routine per Task** protection mechanism (concurrency limiting).
- Real-time, on-demand **Task Cancellation** (using Context).
- The overall **Architecture Concept** behind the project.

This code is meant to help understand the **Technical Depth** of Tasker, not to represent a deployable system.

---

## 🚫 Limitations

- This repository does not include the full backend, frontend, or database code.  
- It may not compile or run directly without missing modules or environment setup.  
- Certain values (such as API keys, email payload, ...) may be mocked or adjusted for clarity and control.  
- Some parts are simplified or redacted to protect intellectual property and security.

---

## 🧠 About Tasker

**Tasker** is a distributed Task Manager written in Go, designed to manage user tasks on demand with fair scheduling, using a pool of persistent and fault-tolerant workers that adapt to real-time actions, to ensure:
- AI-powered email optimization, sending, and broadcasting.  
- AI-powered image generation with prompt optimization
- Cloud API requests

## 📖 How to Use

- Build cmd/server/main.go and start the server 
- Build cmd/client/main.go and start the client 
- From the client terminal, you can run four commands: 
```
1> new_task [task type]
where task type is either email or image 
this command returns the taskID

2> task_status {taskID} 
this command returns the corresponding task status

3> cancel_task {taskID} 
this command cancels the corresponding task 

4> pool_status
This command returns the number of active (busy) workers, available workers, and the overall pool activity rate
```

---

