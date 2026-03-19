CREATE TABLE sesi_teller (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    teller_id       UUID NOT NULL REFERENCES pengguna(id),
    tanggal         DATE NOT NULL,
    saldo_awal      BIGINT NOT NULL DEFAULT 0,
    redenominasi    JSONB NOT NULL DEFAULT '[]',
    saldo_akhir     BIGINT,
    redenominasi_akhir JSONB,
    status          VARCHAR(20) NOT NULL DEFAULT 'AKTIF',
    toleransi_selisih BIGINT NOT NULL DEFAULT 0,
    selisih         BIGINT,
    dibuka_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ditutup_at      TIMESTAMPTZ,
    UNIQUE (teller_id, tanggal)
);

CREATE INDEX idx_sesi_teller_bmt ON sesi_teller(bmt_id, cabang_id, tanggal);
