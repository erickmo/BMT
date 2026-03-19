CREATE TABLE akun (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    kode            VARCHAR(10) NOT NULL,
    nama            VARCHAR(100) NOT NULL,
    tipe            VARCHAR(20) NOT NULL,
    posisi_normal   VARCHAR(10) NOT NULL CHECK (posisi_normal IN ('DEBIT','KREDIT')),
    parent_kode     VARCHAR(10),
    level           SMALLINT NOT NULL DEFAULT 1,
    is_aktif        BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (bmt_id, kode)
);

CREATE TABLE jurnal (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    nomor_jurnal    VARCHAR(40) UNIQUE NOT NULL,
    tanggal         DATE NOT NULL,
    keterangan      TEXT NOT NULL,
    referensi_id    UUID,
    referensi_tipe  VARCHAR(30),
    total_debit     BIGINT NOT NULL,
    total_kredit    BIGINT NOT NULL,
    is_balanced     BOOLEAN NOT NULL GENERATED ALWAYS AS (total_debit = total_kredit) STORED,
    created_by      UUID,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE jurnal_entry (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jurnal_id       UUID NOT NULL REFERENCES jurnal(id),
    akun_kode       VARCHAR(10) NOT NULL,
    posisi          VARCHAR(10) NOT NULL CHECK (posisi IN ('DEBIT','KREDIT')),
    nominal         BIGINT NOT NULL CHECK (nominal > 0),
    keterangan      TEXT
);

CREATE INDEX idx_jurnal_bmt ON jurnal(bmt_id, cabang_id, tanggal DESC);
CREATE INDEX idx_jurnal_entry ON jurnal_entry(jurnal_id);
