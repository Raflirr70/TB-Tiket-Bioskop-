let roomsState = [];
let editingRoomId = null;

const rowLabels = "ABCDEFGHIJKLMNOPQRSTUVWXYZ";

/* ==================== GRID ==================== */

async function fetchRooms() {
    const grid = document.getElementById("rooms-grid");
    try {
        const res = await fetch("/api/v1/rooms");
        if (!res.ok) throw new Error(`HTTP error: ${res.status}`);
        const result = await res.json();
        roomsState = result.data ?? [];
        renderRooms();
    } catch (err) {
        console.error("Gagal mengambil ruangan:", err);
        if (grid) grid.innerHTML = "<p>Gagal mengambil data ruangan.</p>";
    }
}

function renderRooms() {
    const grid = document.getElementById("rooms-grid");
    if (!grid) return;

    const list = Array.isArray(roomsState) ? roomsState : [];

    if (!list.length) {
        grid.innerHTML = "<p>Belum ada ruangan.</p>";
        return;
    }

    grid.innerHTML = list.map(room => {
        const ready = room.status === "ready";
        return `
            <div class="room-card ${ready ? "is-ready" : "is-not-ready"} " data-id="${room.id}">
                <div class="room-card-top">
                    <h3>${room.name}</h3>
                    <span class="room-status ${ready ? "ok" : "no"}">${ready ? "READY" : "NOT READY"}</span>
                </div>
                <div class="room-card-info">
                    <p><strong>${room.total_seats}</strong> kursi</p>
                    <p>Baris: ${room.seat_rows && room.seat_rows.length ? room.seat_rows.join(", ") : "-"}</p>
                    <p>Jadwal aktif: ${room.active_schedules}</p>
                </div>
            </div>
        `;
    }).join("");

    grid.querySelectorAll(".room-card").forEach(card => {
        card.addEventListener("click", () => toggleRoomActions(card));
    });
}

function toggleRoomActions(card) {
    const existing = card.querySelector(".room-actions");
    if (existing) {
        existing.remove();
        return;
    }
    closeAllActions();

    const room = roomsState.find(r => r.id == card.dataset.id);
    const div = document.createElement("div");
    div.className = "room-actions";
    div.innerHTML = `
        <button type="button" class="room-act edit">Edit</button>
        <button type="button" class="room-act del">Hapus</button>
    `;
    div.querySelector(".edit").addEventListener("click", (e) => {
        e.stopPropagation();
        onEdit(room);
    });
    div.querySelector(".del").addEventListener("click", (e) => {
        e.stopPropagation();
        onDelete(room);
    });
    card.appendChild(div);
}

function closeAllActions() {
    document.querySelectorAll(".room-actions").forEach(a => a.remove());
}

/* ==================== EDIT / HAPUS ==================== */

async function onDelete(room) {
    closeAllActions();
    if (!confirm(`Hapus ruangan "${room.name}"?`)) return;

    const res = await fetch(`/api/v1/rooms/${room.id}`, { method: "DELETE" });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
        alert(data.message || "Gagal menghapus ruangan.");
        return;
    }
    await fetchRooms();
}

function onEdit(room) {
    closeAllActions();

    const rows = document.getElementById("seats-rows");
    rows.innerHTML = "";
    document.getElementById("room-name").value = room.name;
    document.getElementById("room-status").value = room.status === "ready" ? "ready" : "not ready";
    document.getElementById("room-form-title").textContent = "Ubah Ruangan";
    document.getElementById("room-name-error").hidden = true;
    document.getElementById("room-name").classList.remove("input-error");

    (room.seat_rows || []).forEach((label, i) => {
        addSeatRow(label, (room.seat_counts || [])[i] || 1);
    });
    if (!rows.children.length) addSeatRow();

    editingRoomId = room.id;
    document.getElementById("room-form-section").hidden = false;
}

/* ==================== FORM TAMBAH ==================== */

function closeForm() {
    editingRoomId = null;
    document.getElementById("room-form-section").hidden = true;
    document.getElementById("room-name").value = "";
    document.getElementById("room-status").value = "ready";
    document.getElementById("room-form-title").textContent = "Tambah Ruangan";
    document.getElementById("room-name-error").hidden = true;
    document.getElementById("room-name").classList.remove("input-error");
    const rows = document.getElementById("seats-rows");
    rows.innerHTML = "";
    addSeatRow();
}

function addSeatRow(label, count) {
    const rows = document.getElementById("seats-rows");
    const index = rows.children.length;
    const l = label || rowLabels[index];

    const row = document.createElement("div");
    row.className = "seat-row";
    row.innerHTML = `
        <span class="seat-label">${l}</span>
        <input type="number" class="seat-count" min="1" value="${count || 10}" placeholder="Jumlah kursi">
        <button type="button" class="seat-remove" title="Hapus baris">&times;</button>
    `;
    row.querySelector(".seat-remove").addEventListener("click", () => {
        row.remove();
        refreshRowLabels();
    });
    rows.appendChild(row);
}

function refreshRowLabels() {
    const rows = document.getElementById("seats-rows");
    [...rows.children].forEach((r, i) => {
        r.querySelector(".seat-label").textContent = rowLabels[i];
    });
}

function collectRoomData() {
    const name = document.getElementById("room-name").value.trim();
    const status = document.getElementById("room-status").value;
    const rows = [...document.querySelectorAll(".seat-row")];

    const seats = rows.map(r => ({
        label: r.querySelector(".seat-label").textContent,
        count: parseInt(r.querySelector(".seat-count").value, 10) || 0
    })).filter(s => s.count > 0);

    return { name, status, seats };
}

function showVerify() {
    const { name, status, seats } = collectRoomData();
    const error = document.getElementById("room-name-error");
    const input = document.getElementById("room-name");

    if (!name) {
        input.classList.add("input-error");
        error.hidden = false;
        input.focus();
        return;
    }
    input.classList.remove("input-error");
    error.hidden = true;

    if (!seats.length) {
        alert("Tambahkan minimal satu baris dengan jumlah kursi.");
        return;
    }

    const total = seats.reduce((s, r) => s + r.count, 0);
    const rowList = seats.map(s => s.label).join(", ");

    document.getElementById("verify-name").textContent = name;
    document.getElementById("verify-rows").textContent = rowList;
    document.getElementById("verify-total").textContent = total;
    document.getElementById("verify-status").textContent = status;

    document.getElementById("room-verify-modal").hidden = false;
}

function closeVerifyModal() {
    document.getElementById("room-verify-modal").hidden = true;
}

async function submitRoom() {
    const { name, status, seats } = collectRoomData();

    const isEdit = editingRoomId !== null;
    const url = isEdit ? `/api/v1/rooms/${editingRoomId}` : "/api/v1/rooms";
    const method = isEdit ? "PUT" : "POST";

    const res = await fetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name, status, seats })
    });
    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
        alert(data.message || "Gagal menyimpan ruangan.");
        return;
    }

    closeVerifyModal();
    closeForm();
    await fetchRooms();
}

/* ==================== INIT ==================== */

document.addEventListener("DOMContentLoaded", () => {
    document.getElementById("add-room-btn").addEventListener("click", () => {
        closeForm();
        document.getElementById("room-form-section").hidden = false;
    });
    document.getElementById("room-form-close").addEventListener("click", closeForm);
    document.getElementById("room-form-cancel").addEventListener("click", closeForm);
    document.getElementById("seats-add-btn").addEventListener("click", addSeatRow);
    document.getElementById("room-form-save").addEventListener("click", showVerify);
    document.getElementById("verify-cancel").addEventListener("click", closeVerifyModal);
    document.getElementById("verify-close-modal").addEventListener("click", closeVerifyModal);
    document.getElementById("verify-confirm").addEventListener("click", submitRoom);

    const nameInput = document.getElementById("room-name");
    nameInput.addEventListener("input", () => {
        document.getElementById("room-name-error").hidden = true;
        nameInput.classList.remove("input-error");
    });

    const modal = document.getElementById("room-verify-modal");
    modal.addEventListener("click", (e) => {
        if (e.target === modal) closeVerifyModal();
    });

    fetchRooms();
});
