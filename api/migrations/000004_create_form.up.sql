CREATE TABLE form_pengajuan (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bmt_id          UUID NOT NULL REFERENCES bmt(id),
    cabang_id       UUID NOT NULL REFERENCES cabang(id),
    jenis_form      VARCHAR(50) NOT NULL,
    nomor_form      VARCHAR(40) UNIQUE NOT NULL,
    status          VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    data_form       JSONB NOT NULL DEFAULT '{}',
    catatan         TEXT,
    diajukan_oleh   UUID NOT NULL REFERENCES pengguna(id),
    diajukan_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE form_approval (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    form_id         UUID NOT NULL REFERENCES form_pengajuan(id),
    urutan          SMALLINT NOT NULL,
    role_approver   VARCHAR(50) NOT NULL,
    approver_id     UUID REFERENCES pengguna(id),
    status          VARCHAR(20) NOT NULL DEFAULT 'MENUNGGU',
    catatan         TEXT,
    approved_at     TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_form_pengajuan_bmt ON form_pengajuan(bmt_id, status, created_at DESC);
CREATE INDEX idx_form_approval_form ON form_approval(form_id, urutan);
