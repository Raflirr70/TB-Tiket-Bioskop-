let allFilms = [];
let carouselTimer = null;

async function fetchFilms() {
    try {
        const response = await fetch("/api/v1/films?sort=ratting");

        if (!response.ok) {
            throw new Error(`HTTP error: ${response.status}`);
        }

        const result = await response.json();

        allFilms = result.data ?? result;

        renderCarousel(allFilms.slice(0, 5));
        renderExplore(allFilms.slice(0, 4));

    } catch (error) {
        console.error("Gagal mengambil data film:", error);

        const track = document.querySelector("#film-carousel .carousel-track");
        const explore = document.getElementById("explore-grid");

        if (track) track.innerHTML = "<p>Gagal mengambil data film.</p>";
        if (explore) explore.innerHTML = "<p>Gagal mengambil data film.</p>";
    }
}

/* ==================== CAROUSEL ==================== */

let carouselState = {
    slides: [],
    active: 0,
    base: 0,
    track: null,
    wrap: null,
    dragStartX: 0,
    baseTransform: 0,
    dragging: false
};

function renderCarousel(films) {
    const track = document.querySelector("#film-carousel .carousel-track");
    const wrap = document.querySelector("#film-carousel .carousel-viewport");
    if (!track || !wrap) return;

    if (!films || films.length === 0) {
        wrap.innerHTML = "<p>Belum ada film.</p>";
        return;
    }

    // render 3x data supaya window tak pernah kehabisan (infinite)
    const tripled = films.concat(films).concat(films);

    track.innerHTML = tripled.map((film) => `
        <div class="carousel-slide" data-id="${film.id}">
            <img src="${film.irl_img}" alt="${film.name}">
            <div class="carousel-overlay">
                <h3>${film.name}</h3>
                <span class="carousel-rating">${formatRating(film.ratting)}</span>
            </div>
        </div>
    `).join("");

    carouselState.slides = [...track.querySelectorAll(".carousel-slide")];
    carouselState.track = track;
    carouselState.wrap = wrap;
    carouselState.base = films.length;
    carouselState.active = carouselState.base; // mulai dari salinan tengah

    carouselState.slides.forEach(slide => {
        slide.addEventListener("click", () => {
            if (carouselState.dragging) return;
            FilmDetail(slide.dataset.id);
        });
    });

    trackTransition(false);
    positionCarousel();
    requestAnimationFrame(() => requestAnimationFrame(() => trackTransition(true)));

    setupCarouselDrag();

    if (carouselTimer) clearInterval(carouselTimer);
    carouselTimer = setInterval(() => goCarousel(1), 5000);
}

function trackStep() {
    const s = carouselState.slides[0];
    if (!s || !carouselState.track) return 0;
    // offsetWidth tak terpengaruh scale(transform), basis asli flex
    const gap = parseFloat(getComputedStyle(carouselState.track).gap || "0");
    return s.offsetWidth + gap;
}

function positionCarousel() {
    const slides = carouselState.slides;
    const track = carouselState.track;
    const wrap = carouselState.wrap;
    if (!slides.length || !track) return;

    const step = trackStep();
    const center = wrap.clientWidth / 2.15;
    const slideW = slides[carouselState.active].offsetWidth;

    track.style.transform = `translateX(${center - slideW / 2.15 - carouselState.active * step}px)`;

    slides.forEach((s, i) => s.classList.toggle("active", i === carouselState.active));
}

function trackTransition(on) {
    if (carouselState.track) carouselState.track.style.transition = on ? "" : "none";
}

function goCarousel(dir) {
    const total = carouselState.slides.length;
    if (!total) return;

    if (carouselState.active + dir < carouselState.base ||
        carouselState.active + dir >= total - carouselState.base) {
        // teleport balik ke tengah tanpa animasi (infinite)
        trackTransition(true);
        carouselState.active = carouselState.base;
        positionCarousel();
        requestAnimationFrame(() => requestAnimationFrame(() => trackTransition(true)));
        return;
    }

    carouselState.active += dir;
    positionCarousel();
}

function resetCarouselTimer() {
    if (carouselTimer) clearInterval(carouselTimer);
    carouselTimer = setInterval(() => goCarousel(1), 5000);
}

function setupCarouselDrag() {
    const track = carouselState.track;
    const wrap = carouselState.wrap;

    const onDown = (e) => {
        if (e.button !== undefined && e.button !== 0) return;
        carouselState.dragging = false;
        carouselState.dragStartX = (e.touches ? e.touches[0].clientX : e.clientX);
        carouselState.baseTransform = parseFloat(
            track.style.transform.replace(/[^-\d.]/g, "") || "0"
        );
        trackTransition(false);
        wrap.classList.add("dragging");
    };

    const onMove = (e) => {
        if (!wrap.classList.contains("dragging")) return;

        const x = (e.touches ? e.touches[0].clientX : e.clientX);
        const dx = x - carouselState.dragStartX;

        if (Math.abs(dx) > 6) carouselState.dragging = true;

        track.style.transform = `translateX(${carouselState.baseTransform + dx}px)`;
    };

    const onUp = (e) => {
        if (!wrap.classList.contains("dragging")) return;
        wrap.classList.remove("dragging");
        trackTransition(true);

        if (carouselState.dragging) {
            const x = (e.touches ? e.changedTouches[0].clientX : e.clientX);
            const dx = x - carouselState.dragStartX;

            if (dx < -40) goCarousel(1);
            else if (dx > 40) goCarousel(-1);
            else positionCarousel();

            resetCarouselTimer();
        }
    };

    wrap.addEventListener("pointerdown", onDown);
    wrap.addEventListener("pointermove", onMove);
    wrap.addEventListener("pointerup", onUp);
    wrap.addEventListener("pointercancel", onUp);

    wrap.addEventListener("wheel", (e) => {
        if (Math.abs(e.deltaY) < Math.abs(e.deltaX)) return; // horizontal scroll tidak dipakai
        e.preventDefault();
        goCarousel(e.deltaY > 0 ? 1 : -1);
        resetCarouselTimer();
    }, { passive: false });
}

/* ==================== EXPLORE FILM ==================== */

function renderExplore(films) {
    const grid = document.getElementById("explore-grid");
    if (!grid) return;

    if (!films || films.length === 0) {
        grid.innerHTML = "<p>Belum ada film.</p>";
        return;
    }

    grid.innerHTML = films.map((film, i) => filmCardHTML(film, i)).join("");
    attachCardClick(grid);
}


function filmCardHTML(film, i) {
    return `
        <div class="film-card" data-id="${film.id}">
            <span class="film-rank">${i + 1}</span>
            <img src="${film.irl_img}" alt="${film.name}">
            <div class="film-card-body">
                <h3>${film.name}</h3>
                <span class="film-rating">${formatRating(film.ratting)}</span>
            </div>
        </div>
    `;
}

function attachCardClick(container) {
    container.querySelectorAll(".film-card").forEach(card => {
        card.addEventListener("click", () => FilmDetail(card.dataset.id));
    });
}

/* ==================== HELPERS ==================== */

function formatRating(r) {
    return `${Number(r || 0).toFixed(1)} ★`;
}

function FilmDetail(filmId) {
    window.location.href = `/film/${filmId}`;
}

/* ==================== INIT ==================== */

document.addEventListener("DOMContentLoaded", () => {
    fetchFilms();

    window.addEventListener("resize", () => {
        if (carouselState.slides.length) positionCarousel();
    });
});
