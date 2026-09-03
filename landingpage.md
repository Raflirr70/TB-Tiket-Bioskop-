# Landing Page — Panduan Pembuatan Tampilan

Dokumen spesifikasi untuk membangun tampilan content landing page BOOCINS.
Bagian yang sudah ada saat ini: hero (`landing_pages.html` + `style.css`), list film sederhana (`film-list` div).
Bagian di bawah = penambahan baru.

**Sumber data**: semua konten film dari `GET /api/v1/films?sort=ratting` (rating rata-rata tertinggi → terendah).
**Template target**: `web/templates/pages/landing_pages.html`
**CSS target**: `web/static/css/style.css`
**JS target**: `web/static/js/landing_page.js`

---

## 1. Carousel Film — 5 Gambar Berputar (Rated)

**Tampilan**: carousel horizontal berisi 5 slide film. Setiap slide = gambar film. Setiap 5 detik slide berpindah otomatis. Sekaligus terlihat 3 slide: slide kiri, slide **tengah (yang paling menonjol / besar / aktif)**, slide kanan. Slide tengah lebih besar dari kiri-kanan.

**Data**: 5 film pertama dari `?sort=ratting` (top 5 rating).

**Struktur HTML**:
```html
<section class="carousel" id="film-carousel">
    <div class="carousel-track">
        <!-- diisi via JS: 5 .carousel-slide -->
    </div>
</section>
```

**JS (`landing_page.js`)**:
- Fetch 5 film top rating → buat `.carousel-slide` tiap film (gambar pakai `film.image` — lihat catatan gambar).
- Slide aktif = `index === activeIndex`. CSS `.carousel-slide.active` → scale lebih besar + z-index atas.
- `setInterval(() => next(), 5000)` → pindah `activeIndex`, re-render class.
- Cara render: buat `.carousel-track` berisi 5 slide, posisi via `transform: translateX` atau flex + `active` class menggeser tengah. Rekomendasi: track flex, slide normal width, slide `.active` diberi `transform: scale(1.15)`.
- Arrow kiri/kanan opsional, auto-slide wajib (5 detik).

**CSS (`style.css`)**:
```css
.carousel { max-width: 1150px; margin: 60px auto; padding: 0 24px; overflow: hidden; }
.carousel-track { display: flex; align-items: center; justify-content: center; gap: 16px; }
.carousel-slide { flex: 0 0 30%; transition: transform 0.5s, opacity 0.5s; opacity: 0.6; }
.carousel-slide.active { transform: scale(1.15); opacity: 1; z-index: 2; }
.carousel-slide img { width: 100%; aspect-ratio: 2/3; object-fit: cover; border-radius: 16px; }
```

---

## 2. Explore Film — 4 Film Horizontal + Grid 6 (View All)

**Tampilan awal**: section berjudul **"Explore Film"**, dengan link **"view all >"** di kanan atas. Tampil 4 film horizontal (kartu: gambar, judul, rating).

**Saat "view all >" diklik**: list melebar jadi **grid 6 film**.

**Struktur HTML**:
```html
<section class="explore">
    <div class="section-heading explore-heading">
        <h2>Explore Film</h2>
        <a href="#" class="view-all-link" id="explore-view-all">view all ></a>
    </div>
    <div id="explore-grid" class="film-row"></div>  <!-- 4 film -->
</section>
```

**JS (`landing_page.js`)**:
- Fetch `?sort=ratting` → ambil 4 pertama → render ke `.film-row`.
- Setiap kartu: `<img>`, `<h3>judul</h3>`, `<span>rating ★</span>`, klik → `FilmDetail(id)`.
- `#explore-view-all` click → `renderFilmsGrid(6)`:
  - Ubah container class: `film-row` → `film-grid` (CSS beda layout).
  - Render 6 film (atau semua bila < 6).
  - Re-render dengan rating yang sama (`?.sort=ratting`).
- Rating diambil dari response `film.ratting` (float rata-rata) — output satu desimal: `film.ratting.toFixed(1)` + `★`.

**CSS (`style.css`)**:
```css
.explore { max-width: 1150px; margin: 0 auto 80px; padding: 0 24px; }
.explore-heading { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 30px; }
.view-all-link { color: #6c3ce9; font-weight: 600; font-size: 14px; }

.film-row { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; }
.film-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 24px; } /* 6 film = 2 baris x 3 */

.film-row .film-card, .film-grid .film-card { border: 1px solid #eee; border-radius: 14px; overflow: hidden; }
.film-card img { width: 100%; aspect-ratio: 2/3; object-fit: cover; }
```

*(Catatan: `.film-card` saat ini di `landing_page.js` tidak ada `<img>` — perlu ditambah gambar + rating pada render.)*

---

## 3. Food & Drink

**Tampilan**: section dengan link **"view all >"** dan **1 gambar utama** food & drink.

**Struktur HTML**:
```html
<section class="foodndrink">
    <div class="section-heading fnb-heading">
        <h2>Food & Drink</h2>
        <a href="#" class="view-all-link">view all ></a>
    </div>
    <div class="fnb-main">
        <img src="/images/[food-drink].avif" alt="Food and Drink">
    </div>
</section>
```

**CSS (`style.css`)**:
```css
.foodndrink { max-width: 1150px; margin: 0 auto 80px; padding: 0 24px; }
.fnb-heading { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 30px; }
.fnb-main img { width: 100%; max-height: 420px; object-fit: cover; border-radius: 20px; }
```

---

## 4. Promo Membership

**Tampilan**: section dengan link **"view all >"** dan **2 gambar utama** promo membership secara horizontal dengan **perbandingan 3:1** (gambar kiri 3 bagian, kanan 1 bagian).

**Struktur HTML**:
```html
<section class="promo">
    <div class="section-heading promo-heading">
        <h2>Promo Membership</h2>
        <a href="#" class="view-all-link">view all ></a>
    </div>
    <div class="promo-images">
        <img src="/images/[promo-1].avif" alt="Promo besar">
        <img src="/images/[promo-2].avif" alt="Promo samping">
    </div>
</section>
```

**CSS (`style.css`)** — rasio 3:1 via `grid-template-columns: 3fr 1fr`:
```css
.promo { max-width: 1150px; margin: 0 auto 100px; padding: 0 24px; }
.promo-heading { display: flex; justify-content: space-between; align-items: flex-end; margin-bottom: 30px; }
.promo-images { display: grid; grid-template-columns: 3fr 1fr; gap: 20px; }
.promo-images img { width: 100%; height: 100%; object-fit: cover; border-radius: 20px; }
```

---

## Catatan Gambar Film

Model `Film` (`internal/domain/entity/film.go`) **tidak punya kolom gambar** — hanya `Name`, `Synopsis`, `Duration`, `Price`, `Status`, dll. Gambar film saat ini belum ada di API response (`FilmResponse` di usecase juga tidak kirim gambar).

**Opsi** (pilih saat implementasi):
1. **Placeholder**: pakai URL gambar statis/placeholder untuk semua film (paling cepat untuk demo tampilan).
2. **Tambah kolom `Image string`** di entity `Film` + update `FilmResponse` + `/api/v1/films` → butuh perubahan backend (migration, seeder `sql.txt`, usecase, handler).

Rekomendasikan opsi 1 dulu untuk bangun tampilan frontend, lalu tambah kolom `Image` belakangan.

---

## Urutan Section di `landing_pages.html`

1. Hero (sudah ada)
2. **Carousel film** (rating)
3. **Explore Film** (4 → grid 6)
4. **Food & Drink**
5. **Promo Membership**
6. About (sudah ada)

## Verifikasi

- `go build ./...` — pastikan backend tak rusak (tidak ada perubahan backend untuk 4 section ini kecuali kolom Image).
- Jalankan server, buka `/`, cek carousel berputar tiap 5 detik, slide tengah menonjol.
- Klik "view all >" Explore → grid jadi 6 film.
- Cek Food & Drink 1 gambar, Promo 2 gambar rasio 3:1.
