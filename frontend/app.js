const API_BASE = "http://localhost:8080";
const CLOUD_NAME = "dfqqy3br";
const UPLOAD_PRESET = "ourDojo"; // Your unsigned preset

let socket;
let isConnected = false;

// === LOGIN LOGIC ===
async function login() {
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value;
  const errorEl = document.getElementById("error");

  if (!username || !password) {
    errorEl.innerText = "Please fill all fields";
    return;
  }

  try {
    const res = await fetch(`${API_BASE}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username, password })
    });

    const data = await res.json();
    
    if (!res.ok || !data.token) {
      errorEl.innerText = data.error || "Login failed";
      return;
    }

    localStorage.setItem("token", data.token);
    localStorage.setItem("username", username);
    window.location.replace("chat.html");
  } catch (err) {
    errorEl.innerText = "Connection error";
    console.error(err);
  }
}

// Handle Enter key on login
if (window.location.pathname.includes("index.html") || window.location.pathname === "/") {
  document.addEventListener("keypress", (e) => {
    if (e.key === "Enter") login();
  });
}

// === CHAT LOGIC ===
if (window.location.pathname.includes("chat.html")) {
  const token = localStorage.getItem("token");
  if (!token) {
    window.location.replace("index.html");
  } else {
    initChat();
    
    // Set the other user's name based on who logged in
    const currentUser = localStorage.getItem("username");
    const otherUserEl = document.getElementById("otherUser");
    
    if (currentUser && otherUserEl) {
      if (currentUser.toLowerCase() === "secretary") {
        otherUserEl.innerText = "Boss";
      } else if (currentUser.toLowerCase() === "boss") {
        otherUserEl.innerText = "Secretary";
      } else {
        // Fallback for any other username
        otherUserEl.innerText = "Other User";
      }
    }
    
    // Prevent going back to login when logged in
    window.history.pushState(null, "", window.location.href);
    window.onpopstate = function() {
      window.history.pushState(null, "", window.location.href);
    };
  }
}

function initChat() {
  const token = localStorage.getItem("token");
  socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

  socket.onopen = () => {
    isConnected = true;
    updateStatus(true);
  };

  socket.onclose = () => {
    isConnected = false;
    updateStatus(false);
  };

  socket.onerror = (err) => {
    console.error("WebSocket error:", err);
    updateStatus(false);
  };

  socket.onmessage = (e) => {
    try {
      const msg = JSON.parse(e.data);
      
      // Handle chat history (array of messages)
      if (Array.isArray(msg)) {
        msg.forEach(historyMsg => {
          const isMe = historyMsg.from_self === "true";
          renderMessage(historyMsg, isMe ? "me" : "other");
        });
      } else {
        // Handle single incoming message
        const isMe = msg.from_self === "true";
        renderMessage(msg, isMe ? "me" : "other");
      }
    } catch (err) {
      console.error("Failed to parse message:", err);
    }
  };

  // Handle Enter key to send
  const input = document.getElementById("messageInput");
  input?.addEventListener("keypress", (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendText();
    }
  });
}

function updateStatus(connected) {
  const statusDot = document.querySelector(".status-dot");
  const statusText = document.querySelector(".status-text");
  
  if (statusDot && statusText) {
    if (connected) {
      statusDot.style.background = "#4ecca3";
      statusText.innerText = "Connected";
      statusText.style.color = "#4ecca3";
    } else {
      statusDot.style.background = "#ff6b6b";
      statusText.innerText = "Disconnected";
      statusText.style.color = "#ff6b6b";
    }
  }
}

function sendText() {
  const input = document.getElementById("messageInput");
  if (!input.value.trim() || !isConnected) return;

  const msg = { 
    type: "text", 
    content: input.value.trim() 
  };
  
  socket.send(JSON.stringify(msg));
  renderMessage(msg, "me");
  input.value = "";
}

function toggleEmoji() {
  const panel = document.getElementById("emojiPanel");
  panel?.classList.toggle("hidden");
}

// Emoji click handler
document.getElementById("emojiPanel")?.addEventListener("click", (e) => {
  if (e.target.classList.contains("emoji-item")) {
    const emoji = e.target.innerText;
    const input = document.getElementById("messageInput");
    
    // Add emoji to existing text instead of sending immediately
    if (input) {
      input.value += emoji;
      input.focus(); // Keep focus on input
    }
    
    toggleEmoji(); // Close emoji panel
  }
});

function pickFile() {
  document.getElementById("fileInput")?.click();
}

async function uploadFile() {
  const fileInput = document.getElementById("fileInput");
  const file = fileInput?.files[0];
  if (!file) return;

  try {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("upload_preset", UPLOAD_PRESET);

    console.log("Uploading to Cloudinary...", {
      cloud_name: CLOUD_NAME,
      upload_preset: UPLOAD_PRESET,
      file_name: file.name,
      file_type: file.type,
      file_size: file.size
    });

    const res = await fetch(
      `https://api.cloudinary.com/v1_1/${CLOUD_NAME}/upload`,
      {
        method: "POST",
        body: formData
      }
    );

    const data = await res.json();
    
    if (!res.ok) {
      console.error("Cloudinary error response:", data);
      throw new Error(`Upload failed: ${data.error?.message || JSON.stringify(data)}`);
    }

    console.log("Upload successful:", data);
    
    let type = "image";
    if (file.type.startsWith("video")) {
      type = "video";
    } else if (file.name.toLowerCase().endsWith(".gif")) {
      type = "gif";
    }

    const msg = { 
      type, 
      media_url: data.secure_url 
    };
    
    socket.send(JSON.stringify(msg));
    renderMessage(msg, "me");
    
    // Reset file input
    fileInput.value = "";
  } catch (err) {
    console.error("Upload failed:", err);
    alert(`Failed to upload file: ${err.message}`);
  }
}

function renderMessage(msg, who) {
  const container = document.getElementById("messages");
  const div = document.createElement("div");
  div.className = `message ${who}`;

  if (msg.type === "text" || msg.type === "emoji") {
    div.innerText = msg.content;
  } else if (msg.type === "image" || msg.type === "gif" || msg.type === "sticker") {
    const img = document.createElement("img");
    img.src = msg.media_url;
    img.alt = "Image";
    img.loading = "lazy";
    div.appendChild(img);
  } else if (msg.type === "video") {
    const video = document.createElement("video");
    video.src = msg.media_url;
    video.controls = true;
    video.preload = "metadata";
    div.appendChild(video);
  }

  container?.appendChild(div);
  div.scrollIntoView({ behavior: "smooth", block: "end" });
}

// Close emoji panel when clicking outside
document.addEventListener("click", (e) => {
  const emojiPanel = document.getElementById("emojiPanel");
  const emojiBtn = document.querySelector(".emoji-btn");
  
  if (emojiPanel && !emojiPanel.contains(e.target) && e.target !== emojiBtn) {
    emojiPanel.classList.add("hidden");
  }
});