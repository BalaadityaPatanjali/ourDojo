const API_BASE = "http://localhost:8080"; // change after deploy
const CLOUD_NAME = "dfqqy3br";
const UPLOAD_PRESET = "ourDojo";

let socket;

// LOGIN
async function login() {
  const username = document.getElementById("username").value;
  const password = document.getElementById("password").value;

  const res = await fetch(`${API_BASE}/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password })
  });

  const data = await res.json();
  if (!data.token) {
    document.getElementById("error").innerText = "Login failed";
    return;
  }

  localStorage.setItem("token", data.token);
  window.location.href = "chat.html";
}

// CHAT INIT
if (window.location.pathname.includes("chat.html")) {
  const token = localStorage.getItem("token");
  if (!token) window.location.href = "index.html";

  socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

  socket.onmessage = (e) => {
    const msg = JSON.parse(e.data);
    renderMessage(msg, "other");
  };
}

function sendText() {
  const input = document.getElementById("messageInput");
  if (!input.value) return;

  const msg = { type: "text", content: input.value };
  socket.send(JSON.stringify(msg));
  renderMessage(msg, "me");
  input.value = "";
}

function toggleEmoji() {
  document.getElementById("emojiPanel").classList.toggle("hidden");
}

document.getElementById("emojiPanel")?.addEventListener("click", (e) => {
  if (e.target.innerText) {
    const msg = { type: "emoji", content: e.target.innerText };
    socket.send(JSON.stringify(msg));
    renderMessage(msg, "me");
  }
});

function pickFile() {
  document.getElementById("fileInput").click();
}

async function uploadFile() {
  const file = document.getElementById("fileInput").files[0];
  if (!file) return;

  const formData = new FormData();
  formData.append("file", file);
  formData.append("upload_preset", UPLOAD_PRESET);

  const res = await fetch(`https://api.cloudinary.com/v1_1/${CLOUD_NAME}/auto/upload`, {
    method: "POST",
    body: formData
  });

  const data = await res.json();
  const type = file.type.startsWith("video") ? "video" : "image";

  const msg = { type, media_url: data.secure_url };
  socket.send(JSON.stringify(msg));
  renderMessage(msg, "me");
}

function renderMessage(msg, who) {
  const div = document.createElement("div");
  div.className = `message ${who}`;

  if (msg.type === "text" || msg.type === "emoji") {
    div.innerText = msg.content;
  } else if (msg.type === "image" || msg.type === "gif" || msg.type === "sticker") {
    const img = document.createElement("img");
    img.src = msg.media_url;
    div.appendChild(img);
  } else if (msg.type === "video") {
    const video = document.createElement("video");
    video.src = msg.media_url;
    video.controls = true;
    div.appendChild(video);
  }

  document.getElementById("messages").appendChild(div);
  div.scrollIntoView();
}