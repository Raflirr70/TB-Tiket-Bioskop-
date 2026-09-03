# Pembuatan Fitur Kelola Rooms & Seats — Spesifikasi

Dokumen perencanaan fitur "Kelola Rooms & Seats" (halaman `/manage-rooms`, dashboard admin).

## 1. Prasyarat / Perubahan Backend

Karena fitur butuh status ruangan dan soft-delete, perlu penyesuaian entity dulu.

### 1.1 Tambah kolom di entity `Room`

File: `internal/domain/entity/room.go` — tambah field:

```go
type Room struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string
	Capasity  int
	Status    string `gorm:"default:ready"` // "ready" | "not_ready" | "deleted"
	DeletedAt gorm.DeletedAt                 // soft delete (log tak hilang)

	Schedules []Schedule
	Seats     []Seat
}
```

- `Status`: "ready" (hijau) → ruangan siap; "not_ready" (merah) → rusak/dipakai.
- `DeletedAt` (GORM soft delete): hapus ruangan = `UPDATE` set `deleted_at`, bukan `DELETE` nyata → log tak hilang. Semua query GORM otomatis filter `deleted_at IS NULL`.

> Perlu `gorm.io/gorm` import. Migration `db.AutoMigrate` (cmd/migration/up/main.go) otomatis tambah kolom.

### 1.2 Aturan validasi hapus

- Ruangan **tidak bisa dihapus** jika masih ada `Schedule` **aktif** (`room_id` = id, `status != 'cancelled'`, dan `date/time` belum lewat).
- Penghapusan hanya **soft delete** (`deleted_at` diisi) + `status = "deleted"` — data riwayat tetap.

## 2. Fitur

| # | Fitur | Keterangan |
|---|---|---|
| 1 | Grid ruangan | Tiap ruangan satu kartu grid |
| 2 | Jumlah kursi | Dihitung dari jumlah row `seats` milik ruangan |
| 3 | Status ruangan | Badge berwarna: hijau = `ready`, merah = `not_ready` |
| 4 | Aksi edit/hapus | Klik kartu → muncul pilihan Edit & Hapus |
| 5 | Hapus beraturan | Cek schedule aktif → tolak; soft delete |
| 6 | Tambah ruangan | Tombol → form (bukan popup, dirender di halaman) |
| 7 | Tambah seats baris | Input text integer per label; label A, tombol + → otomatis label B, dst |
| 8 | Simpan + verifikasi | Konfirmasi menampilkan data penting (nama, total kursi, dll) sebelum commit |

## 3. Aturan Pembuatan

1. **Layout**: halaman `/manage-rooms` berbasis dashboard (sidebar + main content).
2. **Tampilkan form di halaman** (inline section), bukan modal/popup — sesuai permintaan.
3. **Verifikasi simpan**: tampilkan ringkasan data sebelum benar-benar menyimpan, dengan tombol `Konfirmasi` & `Batal`.
4. **Grid klik** → tampilkan panel aksi (Edit / Hapus) — bisa dropdown/inline di kartu, bukan popup.
5. **Soft delete**: hapus ruangan wajib lewat `deleted_at`, jangan `DELETE` fisik.
6. **Hapus tolak** bila ada schedule aktif → tampilkan pesan error jelas.
7. **Baris kursi** otomatis: mulai `A`, tombol "tambah baris" menambah label berikutnya (`B`, `C`, ...). Tiap label ada input text integer = jumlah kursi di baris itu.
8. Warna status konsisten dengan palette: hijau `#22c55e`, merah `#ef4444`.

## 4. Tampilan (UI)

### 4.1 Header halaman

```
┌────────────────────────────────────────────┐
│  Kelola Rooms & Seats        [+ Tambah Ruangan]  │
│  (subtitle: kelola ruangan & kursi bioskop)      │
└────────────────────────────────────────────┘
```

### 4.2 Grid Ruangan

Tiap kartu:

```
┌───────────────┐
│  Studio 1     │  ● READY (hijau) / ● NOT READY (merah)
│  Kapasitas:    │
│  3 baris       │
│  144 kursi     │
│  [jadwal aktif: n] │
└───────────────┘
```

- Warna badge mengikuti `status`.
- Klik kartu → muncul `Edit` & `Hapus`.

### 4.3 Form Tambah Ruangan (inline)

```
┌────────────────────────────────────────┐
│ Tambah Ruangan                          │
│ ┌────────────────────────────────────┐  │
│ │ Nama Ruangan: [ Studio 4         ] │  │
│ └────────────────────────────────────┘  │
│                                          │
│ BARIS KURSI        JUMLAH KURSI          │
│ A                  [  10 ]               │
│ B                  [  12 ]   [+ Tambah Baris]  │
│ C                  [  12 ]               │
│                                          │
│ [Batal]  [Simpan]                        │
└────────────────────────────────────────┘
```

- Tombol `+ Tambah Baris` setelah baris terakhir → otomatis label huruf berikutnya.
- Tiap baris: `input type="text"` label (A/B/C...) + `input type="number"` jumlah kursi.
- Baris boleh dihapus (`×`).
- `Simpan` → buka panel **verifikasi** (bukan langsung simpan).

### 4.4 Verifikasi Sebelum Simpan

```
┌────────────────────────────────────────┐
│ Konfirmasi Ruangan Baru                 │
│  Nama      : Studio 4                   │
│  Baris     : A, B, C                    │
│  Total Kursi: 34                        │
│  Status    : ready                      │
│  ┌────────────────────────────┐         │
│  │ [Batal]        [Konfirmasi] │         │
│  └────────────────────────────┘         │
└────────────────────────────────────────┘
```

- Menampilkan data penting. `Konfirmasi` → submit ke backend.

## 5. Backend (endpoint yang akan dibuat)

> Belum dibuat sekarang — hanya API target agar frontend siap.

| Method | Path | Fungsi |
|---|---|---|
| GET | `/api/v1/rooms` | list ruangan (non-deleted) + jumlah kursi + jumlah jadwal |
| GET | `/api/v1/rooms/:id` | detail ruangan + daftar seats |
| POST | `/api/v1/rooms` | buat ruangan + seats (body: `name`, `seats: [{label, count}]`) |
| PUT | `/api/v1/rooms/:id` | update nama/status |
| DELETE | `/api/v1/rooms/:id` | soft delete, tolak bila ada schedule aktif |
| POST | `/api/v1/rooms/:id/seats` | tambah baris kursi ke ruangan |

Body contoh POST `/api/v1/rooms`:

```json
{
  "name": "Studio 4",
  "seats": [
    { "label": "A", "count": 10 },
    { "label": "B", "count": 12 }
  ]
}
```

Respons list room (yang dipakai grid):

```json
[
  {
    "id": 1,
    "name": "Studio 1",
    "status": "ready",
    "total_seats": 144,
    "seat_rows": ["A", "B", "C"],
    "active_schedules": 3
  }
]
```

## 6. File yang Akan Diubah/Dibuat

**Backend**
- `internal/domain/entity/room.go` — tambah `Status`, `DeletedAt`
- **Baru** `internal/delivery/http/handler/room_handler.go`
- `internal/delivery/http/router/route.go` — daftar `/api/v1/rooms*`
- *(repo/usecase room baru — sesuai instruksi, belum dibuat, ditambah belakangan)*

**Frontend**
- `web/templates/pages/dashboard.html` — section `Kelola Rooms & Seats` (ganti placeholder)
- `web/static/js/room_page.js` — baru (render grid, form, verifikasi)
- `web/static/css/style.css` / `navbar.css` — style grid, badge status, form, verifikasi

## 7. Struktur Data Kursi (baris)

- `Seat.Name` = label baris + nomor, contoh `A1`, `A2`, ..., `A10`, `B1`, ..., `B12`.
- Untuk `label: "A", count: 10` → generate 10 seat: `A1`..`A10`.
- `total_seats` = jumlah semua seat.

## 8. Aturan Tambahan (dipertimbangkan)

- Perbarui `Capasity` otomatis = total kursi saat simpan/tambah baris.
- Input tak boleh nol/negatif (min 1) dan nama ruangan wajib.
- Dua ruangan tak boleh nama sama (validasi unique).
- Edit memungkinkan ubah nama & status (`ready`/`not_ready`).
