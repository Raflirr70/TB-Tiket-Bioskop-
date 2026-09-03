let allFilms = [];

async function fetchFilms() {
    try {
        const response = await fetch("/api/v1/films");

        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }

        const result = await response.json();

        allFilms = result.data ?? result;

        renderFilms(currentFilter());
    } catch (error) {
        console.error("Gagal mengambil data film:", error);

        const grid = document.getElementById("films-grid");
        if (grid) grid.innerHTML = "<p>Gagal mengambil data film.</p>";
    }
}

function currentFilter() {
    const active = document.querySelector(".films-tab.active");
    return active ? active.dataset.filter : "now_playing";
}

function filterFilms(f) {
    return allFilms.filter(film => {
        if (f === "now_playing") return film.status === "regular";
        return film.status === "pre sale" || film.status === "comming soon";
    });
}

function renderFilms(filter) {
    const grid = document.getElementById("films-grid");
    if (!grid) return;

    const films = filterFilms(filter);

    if (!films.length) {
        grid.innerHTML = "<p>Belum ada film.</p>";
        return;
    }

    grid.innerHTML = films.map((film, i) => `
        <div class="film-card" data-id="${film.id}">
            <span class="film-rank">${i + 1}</span>
            <img src="${film.irl_img}" alt="${film.name}">
            <div class="film-card-body">
                <h3>${film.name}</h3>
                <span class="film-rating">${formatRating(film.ratting)}</span>
                <p class="film-genre">${film.genres?.map(g => g.name).join(", ") || "-"}</p>
            </div>
        </div>
    `).join("");

    grid.querySelectorAll(".film-card").forEach(card => {
        card.addEventListener("click", () => FilmDetail(card.dataset.id));
    });
}

function formatRating(r) {
    return `${Number(r || 0).toFixed(1)} ★`;
}

function FilmDetail(filmId) {
    window.location.href = `/film/${filmId}`;
}

document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll(".films-tab").forEach(tab => {
        tab.addEventListener("click", () => {
            document.querySelectorAll(".films-tab").forEach(t => t.classList.remove("active"));
            tab.classList.add("active");
            renderFilms(tab.dataset.filter);
        });
    });

    fetchFilms();
});
