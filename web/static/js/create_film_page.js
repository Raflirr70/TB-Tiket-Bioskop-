let availableRooms = [];

document.addEventListener("DOMContentLoaded", () => {
    setupDragAndDrop();
    fetchRoomsForSchedule();
    setupScheduleManager();
    setupFormSubmit();
});

/* ==================== DRAG & DROP POSTER ==================== */
function setupDragAndDrop() {
    const dropZone = document.getElementById("drag-drop-zone");
    const fileInput = document.getElementById("poster-file-input");
    const content = document.getElementById("drop-zone-content");
    const previewContainer = document.getElementById("poster-preview-container");
    const previewImg = document.getElementById("poster-preview-img");
    const removeBtn = document.getElementById("btn-remove-poster");
    const hiddenUrlInput = document.getElementById("irl-img-url");

    if (!dropZone || !fileInput) return;

    dropZone.addEventListener("click", (e) => {
        if (e.target !== removeBtn && !removeBtn.contains(e.target)) {
            fileInput.click();
        }
    });

    ["dragenter", "dragover"].forEach(eventName => {
        dropZone.addEventListener(eventName, (e) => {
            e.preventDefault();
            e.stopPropagation();
            dropZone.classList.add("drag-active");
        });
    });

    ["dragleave", "drop"].forEach(eventName => {
        dropZone.addEventListener(eventName, (e) => {
            e.preventDefault();
            e.stopPropagation();
            dropZone.classList.remove("drag-active");
        });
    });

    dropZone.addEventListener("drop", (e) => {
        const dt = e.dataTransfer;
        const files = dt.files;
        if (files && files.length > 0) {
            handleFileSelected(files[0]);
        }
    });

    fileInput.addEventListener("change", (e) => {
        if (fileInput.files && fileInput.files.length > 0) {
            handleFileSelected(fileInput.files[0]);
        }
    });

    removeBtn.addEventListener("click", (e) => {
        e.stopPropagation();
        fileInput.value = "";
        hiddenUrlInput.value = "";
        previewImg.src = "";
        previewContainer.hidden = true;
        content.hidden = false;
    });

    async function handleFileSelected(file) {
        if (!file.type.startsWith("image/")) {
            alert("File yang dipilih harus berupa gambar.");
            return;
        }

        // 1. Local Preview
        const reader = new FileReader();
        reader.onload = (e) => {
            previewImg.src = e.target.result;
            content.hidden = true;
            previewContainer.hidden = false;
        };
        reader.readAsDataURL(file);

        // 2. Upload file ke server
        const formData = new FormData();
        formData.append("poster", file);

        try {
            const response = await fetch("/api/v1/films/upload", {
                method: "POST",
                body: formData
            });

            if (response.ok) {
                const resData = await response.json();
                if (resData.data && resData.data.url) {
                    hiddenUrlInput.value = resData.data.url;
                }
            } else {
                console.warn("Gagal mengupload gambar ke server, menggunakan preview lokal.");
            }
        } catch (err) {
            console.warn("Error upload file poster:", err);
        }
    }
}

/* ==================== FETCH ROOMS ==================== */
async function fetchRoomsForSchedule() {
    try {
        const response = await fetch("/api/v1/rooms");
        if (response.ok) {
            const result = await response.json();
            availableRooms = result.data ?? result ?? [];
        }
    } catch (err) {
        console.error("Gagal mengambil data ruangan:", err);
    }
}

/* ==================== SCHEDULE MANAGER ==================== */
function setupScheduleManager() {
    const btnAdd = document.getElementById("btn-add-schedule");
    if (btnAdd) {
        btnAdd.addEventListener("click", () => {
            addScheduleRow();
        });
    }
}

function addScheduleRow(dateVal = "", timeVal = "", roomIdVal = "") {
    const list = document.getElementById("schedules-list");
    if (!list) return;

    const row = document.createElement("div");
    row.className = "schedule-item-row";

    const roomOptionsHTML = availableRooms.map(r => `
        <option value="${r.id}" ${r.id == roomIdVal ? "selected" : ""}>${r.name} (${r.total_seats || r.capacity || 0} Kursi)</option>
    `).join("");

    row.innerHTML = `
        <div class="schedule-field">
            <label>Tanggal</label>
            <input type="date" class="schedule-date" value="${dateVal}" required>
        </div>
        <div class="schedule-field">
            <label>Waktu</label>
            <input type="time" class="schedule-time" value="${timeVal}" required>
        </div>
        <div class="schedule-field">
            <label>Ruangan Studio</label>
            <select class="schedule-room" required>
                <option value="">-- Pilih Studio --</option>
                ${roomOptionsHTML}
            </select>
        </div>
        <button type="button" class="btn-remove-schedule" title="Hapus Jadwal">&times;</button>
    `;

    row.querySelector(".btn-remove-schedule").addEventListener("click", () => {
        row.remove();
    });

    list.appendChild(row);
}

/* ==================== SUBMIT FORM ==================== */
function setupFormSubmit() {
    const form = document.getElementById("create-film-form");
    if (!form) return;

    form.addEventListener("submit", async (e) => {
        e.preventDefault();

        const name = document.getElementById("film-name").value.trim();
        const status = document.getElementById("film-status").value;
        const duration = parseInt(document.getElementById("film-duration").value, 10) || 0;
        const price = parseInt(document.getElementById("film-price").value, 10) || 0;
        const synopsis = document.getElementById("film-synopsis").value.trim();
        const irlImg = document.getElementById("irl-img-url").value.trim() || document.getElementById("poster-preview-img").src;

        if (!name) {
            alert("Judul film wajib diisi.");
            return;
        }

        // Collect Schedules
        const scheduleRows = document.querySelectorAll(".schedule-item-row");
        const schedules = [];

        scheduleRows.forEach(row => {
            const date = row.querySelector(".schedule-date").value;
            const time = row.querySelector(".schedule-time").value;
            const roomId = parseInt(row.querySelector(".schedule-room").value, 10) || 0;

            if (date && time && roomId > 0) {
                schedules.push({
                    room_id: roomId,
                    date: date,
                    time: time
                });
            }
        });

        const payload = {
            name: name,
            status: status,
            duration: duration,
            price: price,
            synopsis: synopsis,
            irl_img: irlImg,
            schedules: schedules
        };

        const submitBtn = document.getElementById("btn-save-film");
        if (submitBtn) {
            submitBtn.disabled = true;
            submitBtn.innerHTML = `<i class="uil uil-spinner"></i> Menyimpan...`;
        }

        try {
            const response = await fetch("/api/v1/films", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload)
            });

            if (response.ok) {
                alert("Film berhasil ditambahkan!");
                window.location.href = "/manage-films";
            } else {
                const resData = await response.json().catch(() => ({}));
                alert(resData.message || "Gagal menambah film baru.");
                if (submitBtn) {
                    submitBtn.disabled = false;
                    submitBtn.innerHTML = `<i class="uil uil-check"></i> Simpan Film Baru`;
                }
            }
        } catch (err) {
            console.error("Error submit film:", err);
            alert("Terjadi kesalahan koneksi.");
            if (submitBtn) {
                submitBtn.disabled = false;
                submitBtn.innerHTML = `<i class="uil uil-check"></i> Simpan Film Baru`;
            }
        }
    });
}
