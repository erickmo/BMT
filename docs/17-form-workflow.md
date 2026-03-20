# 17 — Form Workflow & Approval Engine

> **Terakhir diperbarui:** 20 Maret 2026

## Prinsip
**Nasabah dan rekening tidak pernah diubah langsung.** Semua mutasi via form ber-status.

## State Machine
```
DRAFT → DIAJUKAN → DISETUJUI → DITERAPKAN
                ↘ DITOLAK → (revisi) → DIAJUKAN
```

## Jenis Form
| Kode | Efek Setelah DITERAPKAN |
|------|------------------------|
| `FORM_DAFTAR_NASABAH` | INSERT nasabah + generate nomor |
| `FORM_UBAH_NASABAH` | UPDATE kolom nasabah |
| `FORM_BUKA_REKENING` | INSERT rekening + generate nomor |
| `FORM_UBAH_REKENING` | UPDATE rekening (non-finansial saja) |
| `FORM_BLOKIR_REKENING` | UPDATE status → `BEKU` |
| `FORM_BUKA_BLOKIR` | UPDATE status → `AKTIF` |
| `FORM_TUTUP_REKENING` | UPDATE status → `TUTUP` |
| `FORM_BUKA_PEMBIAYAAN` | INSERT pembiayaan |
| `FORM_AKAD_PEMBIAYAAN` | UPDATE status pembiayaan → `AKAD` |

## Konfigurasi Approver (BMT Settings)
```json
"approval.FORM_DAFTAR_NASABAH":  ["TELLER", "MANAJER_CABANG"]
"approval.FORM_UBAH_NASABAH":    ["MANAJER_CABANG"]
"approval.FORM_BUKA_REKENING":   ["TELLER", "MANAJER_CABANG"]
"approval.FORM_BLOKIR_REKENING": ["MANAJER_CABANG"]
"approval.FORM_TUTUP_REKENING":  ["MANAJER_CABANG"]
"approval.FORM_BUKA_PEMBIAYAAN": ["KOMITE"]
```

## Skema Tabel
```sql
form_pengajuan (bmt_id, cabang_id,
                jenis_form VARCHAR(30),
                nomor_form VARCHAR(50) UNIQUE,
                -- Format: {KODE_BMT}/{KODE_CAB}/{FORM}/{TAHUN}/{SEQ:05d}
                status VARCHAR(20),
                data_form JSONB,         -- snapshot seluruh field
                referensi_id UUID,
                referensi_tipe VARCHAR(30),
                alasan_tolak TEXT,
                diajukan_oleh, diajukan_at,
                disetujui_oleh, disetujui_at,
                ditolak_oleh, ditolak_at,
                diterapkan_at)

form_riwayat (form_id, status_dari, status_ke,
              catatan, oleh UUID, created_at)
```

## Aturan Penting
- `FORM_UBAH_REKENING` hanya bisa ubah kolom non-finansial
- Saldo rekening **hanya** berubah via `transaksi_rekening`
- Form `DITERAPKAN` tidak bisa dibatalkan
- Setiap transisi dicatat di `form_riwayat`
- Notifikasi FCM ke approver saat form baru `DIAJUKAN`
