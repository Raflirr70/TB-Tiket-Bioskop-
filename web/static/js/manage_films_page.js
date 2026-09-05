let filmsData = [];
let currentStatusFilter = "all";
let currentSearchQuery = "";

async function fetchManageFilms() {
    const tbody = document.getElementById("manage-films-tbody");
    try {
        const response = await fetch("/api/v1/films");
        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }

        const result = await response.json();
        filmsData = result.data ?? result ?? [];
        renderFilmsTable(currentStatusFilter);
    } catch (error) {
        console.error("Gagal mengambil data film:", error);
        if (tbody) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="7" class="table-empty">Gagal mengambil data film.</td>
                </tr>
            `;
        }
    }
}

function renderFilmsTable(statusFilter) {
    const tbody = document.getElementById("manage-films-tbody");
    if (!tbody) return;

    let filtered = filmsData;
    if (statusFilter !== "all") {
        filtered = filtered.filter(film => {
            const status = (film.status || "").toLowerCase();
            return status === statusFilter.toLowerCase();
        });
    }

    if (currentSearchQuery) {
        const query = currentSearchQuery.toLowerCase();
        filtered = filtered.filter(film => {
            const nameMatch = (film.name || "").toLowerCase().includes(query);
            const synopsisMatch = (film.synopsis || "").toLowerCase().includes(query);
            return nameMatch || synopsisMatch;
        });
    }

    if (!filtered || filtered.length === 0) {
        tbody.innerHTML = `
            <tr>
                <td colspan="7" class="table-empty">Tidak ada film yang cocok dengan pencarian / status ini.</td>
            </tr>
        `;
        return;
    }

    tbody.innerHTML = filtered.map((film, index) => {
        const synopsisText = truncateText(film.synopsis || "-", 90);
        const posterImg = film.irl_img || "https://placehold.co/100x150/6c3ce9/ffffff?text=Poster";
        const formattedPrice = Number(film.price || 0).toLocaleString("id-ID");
        const statusBadge = getStatusBadge(film.status);

        return `
            <tr>
                <td class="text-center font-bold">${index + 1}</td>
                <td class="text-center">
                    <img src="${posterImg}" alt="${film.name}" class="film-poster-thumb">
                </td>
                <td class="font-semibold">${film.name || "-"}</td>
                <td class="synopsis-cell" title="${film.synopsis || ""}">${synopsisText}</td>
                <td>${film.duration || 0} menit</td>
                <td>Rp ${formattedPrice}</td>
                <td>${statusBadge}</td>
            </tr>
        `;
    }).join("");
}

function truncateText(text, maxLength) {
    if (!text || text.length <= maxLength) {
        return text || "";
    }
    return text.substring(0, maxLength).trim() + "...";
}

function getStatusBadge(status) {
    const st = (status || "").toLowerCase();
    let badgeClass = "status-regular";
    let label = status || "Regular";

    if (st === "pre sale") {
        badgeClass = "status-presale";
        label = "Pre Sale";
    } else if (st === "comming soon" || st === "coming soon") {
        badgeClass = "status-comingsoon";
        label = "Coming Soon";
    } else if (st === "regular") {
        badgeClass = "status-regular";
        label = "Regular";
    }

    return `<span class="film-status-badge ${badgeClass}">${label}</span>`;
}

document.addEventListener("DOMContentLoaded", () => {
    // Setup tabs
    const tabsContainer = document.getElementById("manage-films-tabs");
    if (tabsContainer) {
        const tabs = tabsContainer.querySelectorAll(".film-tab");
        tabs.forEach(tab => {
            tab.addEventListener("click", () => {
                tabs.forEach(t => t.classList.remove("active"));
                tab.classList.add("active");
                currentStatusFilter = tab.dataset.status || "all";
                renderFilmsTable(currentStatusFilter);
            });
        });
    }

    // Setup search input
    const searchInput = document.getElementById("film-search-input");
    if (searchInput) {
        searchInput.addEventListener("input", (e) => {
            currentSearchQuery = e.target.value.trim();
            renderFilmsTable(currentStatusFilter);
        });
    }

    // Setup tombol tambah film
    const addBtn = document.getElementById("add-film-btn");
    if (addBtn) {
        addBtn.addEventListener("click", () => {
            window.location.href = "/manage-films/add";
        });
    }

    fetchManageFilms();
});
