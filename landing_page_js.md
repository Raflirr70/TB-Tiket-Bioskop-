# landing_page.js — Panduan Penyesuaian dengan Landing Page

Dokumen spesifikasi implementasi `web/static/js/landing_page.js` agar sesuai dengan struktur HTML yang sudah ada di `web/templates/pages/landing_pages.html`.

**Struktur HTML yang sudah ada** (yang harus diisi/digenerate JS):

| Selector | Konten |
|---|---|
| `#film-carousel .carousel-track` | 5 slide film (carousel berputar) |
| `#explore-grid.film-row` | 4 kartu film horizontal |
| `#explore-view-all` | link, klik → ubah `#explore-grid` jadi `film-grid` + isi 6 film |
| `#film-list.film-grid` | *(legacy dari kode lama — reuse jika masih dipakai)* |

**Sumber data**: `GET /api/v1/films?sort=ratting` → array film, rating rata-rata tertinggi duluan.
Tiap film punya field dari `FilmResponse` (usecase): `id, name, synopsis, duration, price, status, ratting (float rata-rata), genres[], schedules[]`.

> **Catatan gambar**: film **tidak punya field gambar** di API response. Semua `<img>` film pakai placeholder. Lihat bagian "Gambar Placeholder" di bawah.

---

## 1. Fetch Data (satu kali, reuse)

Fetch sekali dengan `?sort=ratting`, simpan ke variabel global, lalu render ke semua section.

```js
let allFilms = [];

async function fetchFilms() {
    try {
        const res = await fetch("/api/v1/films?sort=ratting");
        if (!res.ok) throw new Error(`HTTP error: ${res.status}`);
        const result = await res.json();
        allFilms = result.data ?? result;

        renderCarousel(allFilms.slice(0, 5));   // 5 slide
        renderExplore(allFilms.slice(0, 4));    // 4 horizontal
    } catch (error) {
        console.error("Gagal mengambil data film:", error);
    }
}
```

---

## 2. Carousel (5 Slide, Auto 5 Detik, Tengah Menonjol)

**Logika**:
- 5 slide. Ada indeks aktif `activeIndex`.
- Posisi: render kelima slide sebagai flex di `.carousel-track`. Slide `.active` → CSS scale besar + opacity 1.
- Setiap 5 detik `activeIndex` bertambah (wrap ke 0), re-generate class aktif.

```js
let carouselTimer = null;

function renderCarousel(films) {
    const track = document.querySelector("#film-carousel .carousel-track");
    if (!track) return;

    track.innerHTML = films.map((film, i) => `
        <div class="carousel-slide${i === 0 ? " active" : ""}" data-index="${i}">
            <img src="${filmImage(film)}" alt="${film.name}">
            <p class="carousel-title">${film.name}</p>
            <p class="carousel-rating">${formatRating(film.ratting)}</p>
        </div>
    `).join("");

    startCarousel();
}

function startCarousel() {
    if (carouselTimer) clearInterval(carouselTimer);

    carouselTimer = setInterval(() => {
        const slides = document.querySelectorAll(".carousel-slide");
        if (!slides.length) return;

        // cari indeks aktif sekarang, pindah ke berikutnya (wrap)
        let active = [...slides].findIndex(s => s.classList.contains("active"));
        const next = (active + 1) % slides.length;

        slides[active].classList.remove("active");
        slides[next].classList.add("active");
    }, 5000);
}
```

**Catatan tampilan**: render ulang full track tiap pergantian slide (paling sederhana) ATAU posisikan lewat transform. Untuk "3 terlihat, tengah besar" — CSS `.carousel-slide { flex: 0 0 30% }` + `.active { transform: scale(1.15) }`. Urutan flex alami: tengah = indeks 2. Opsional ekspansi — kalau mau slide tengah benar-benar center, tambah state offset; untuk v1 cukup active-scale.

---

## 3. Explore Film (4 ⇒ Grid 6 via "view all")

**Logika**:
- Render 4 kartu ke `#explore-grid` (class `film-row`).
- Link `#explore-view-all` klik → ganti class container jadi `film-grid`, render 6 film.
- Kartu: gambar, judul, rating. Klik → `FilmDetail(id)`.

```js
function renderExplore(films) {
    const grid = document.getElementById("explore-grid");
    if (!grid) return;
    grid.innerHTML = films.map(filmCardHTML).join("");
    attachCardClick(grid);
}

function renderExploreAll() {
    const grid = document.getElementById("explore-grid");
    if (!grid) return;
    grid.classList.remove("film-row");
    grid.classList.add("film-grid");          // switch ke grid 6 (2x3)
    grid.innerHTML = allFilms.slice(0, 6).map(filmCardHTML).join("");
    attachCardClick(grid);
}

function filmCardHTML(film) {
    return `
        <div class="film-card" data-id="${film.id}">
            <img src="${filmImage(film)}" alt="${film.name}">
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

document.getElementById("explore-view-all")?.addEventListener("click", (e) => {
    e.preventDefault();          // cegah lompat ke "#"
    renderExploreAll();
});
```

---

## 4. Helper

Format rating rata-rata (1 desimal) + bintang:
```js
function formatRating(r) {
    return `${Number(r || 0).toFixed(1)} ★`;
}
```

Gambar placeholder — film tak punya field gambar, pakai placeholder statis:
```js
function filmImage(film) {
    // placeholder: ganti dengan URL/gambar nyata bila kolom Image ditambah
    return "/images/[film-placeholder].avif";
}
```

---

## 5. FilmDetail & Init (tetap)

```js
function FilmDetail(filmId) {
    window.location.href = `/film/${filmId}`;
}

document.addEventListener("DOMContentLoaded", () => {
    fetchFilms();
});
```

---

## Verifikasi

- Buka `/`, carousel berisi 5 slide, slide tengah menonjol (scale besar), berpindah tiap 5 detik.
- Explore Film tampil 4 horizontal (gambar, judul, rating).
- Klik "view all >" → jadi grid 6 film.
- Klik kartu → navigasi `/film/{id}`.

## Poin yang perlu diputuskan saat implementasi

1. **Gambar film** — API tak kirim field gambar. Opsi: (a) placeholder statis semua (cepat, untuk demo), (b) tambah kolom `Image` di entity `Film` + `FilmResponse` + update `/api/v1/films` (butuh perubahan backend + migration + seeder).
2. **`#film-list` lama** — section list film lama di bawah jugadipakai `renderFilms`. Bisa dihapus bila semua konten pindah ke Explore, atau dibiarkan. Rekomendasi: hapus `renderFilms` lama, gunakan `renderExplore`.
3. **Carousel "3 terlihat"** — v1 pakai active-scale (5 flex, yang tengah lebih besar). Kalau mau tepat 3 terlihat (kiri/kanan terpotong), butuh CSS overflow hidden + offset transform.
