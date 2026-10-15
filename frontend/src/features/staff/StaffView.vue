<template>
  <main class="dashboard">
    <!-- Mobile backdrop for sidebar overlay -->
    <div
      v-if="!sidebarCollapsed"
      class="sidebar-mobile-backdrop"
      @click="sidebarCollapsed = true"
    />

    <!-- Reusable Sticky StaffSidebar Component (Pinned bottom for instant logout) -->
    <StaffSidebar
      v-model:collapsed="sidebarCollapsed"
      :doctor-name="doctorName"
      @logout="handleLogout"
    />

    <div class="dashboard-body">
      <div class="dashboard-heading">
        <div>
          <div class="section-label">PUSKESMAS SEHAT BERSAMA</div>
          <h2>
            {{
              queueStore.role === "loket"
                ? "Dashboard Loket"
                : queueStore.role === "medis"
                ? "Ruang Pemeriksaan"
                : queueStore.role === "database_pasien"
                ? "Database Master Pasien"
                : "Database Katalog Obat & Farmasi"
            }}
          </h2>
          <p>
            {{ currentFormattedDate }}
          </p>
        </div>
        <button class="button outline" @click="showPendingNotice">
          <Bell :size="18" /> Notifikasi
        </button>
      </div>

      <!-- Stats Grid (Role-Aware) -->
      <div v-if="queueStore.role === 'database_pasien'" class="stats-grid">
        <div class="stat-card">
          <span>Total Pasien Terdaftar</span>
          <strong>{{ queueStore.patients.length }}</strong>
          <Users :size="22" />
        </div>
        <div class="stat-card">
          <span>Pasien BPJS Kesehatan</span>
          <strong>{{ queueStore.patients.filter(p => p.insurance === 'BPJS').length }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Pasien Umum</span>
          <strong>{{ queueStore.patients.filter(p => p.insurance !== 'BPJS').length }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Catatan Riwayat Alergi</span>
          <strong>{{ queueStore.patients.filter(p => p.allergy && !/^(tidak ada|tidak|-)$/i.test(p.allergy)).length }}</strong>
          <Activity :size="22" />
        </div>
      </div>
      <div v-else-if="queueStore.role === 'database_obat'" class="stats-grid">
        <div class="stat-card">
          <span>Total Katalog Obat</span>
          <strong>{{ medications.length }}</strong>
          <Pill :size="22" />
        </div>
        <div class="stat-card">
          <span>Sediaan Oral (Minum/Kunyah)</span>
          <strong>{{ medications.filter(m => m.route?.includes('oral') || m.route?.includes('sirup')).length }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Topikal, Tetes & Kumur</span>
          <strong>{{ medications.filter(m => m.route?.includes('topikal') || m.route?.includes('tetes') || m.route?.includes('kumur')).length }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Total Stok Fisik Obat</span>
          <strong>{{ medications.reduce((acc, m) => acc + (m.stock || 0), 0) }}</strong>
          <Activity :size="22" />
        </div>
      </div>
      <div v-else class="stats-grid">
        <div class="stat-card">
          <span>Pendaftaran hari ini</span>
          <strong>{{ todayRegistrationsCount }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Menunggu verifikasi</span>
          <strong>{{ pendingVerificationCount }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Siap dilayani</span>
          <strong>{{ readyToServeCount }}</strong>
          <Activity :size="22" />
        </div>
        <div class="stat-card">
          <span>Kunjungan selesai</span>
          <strong>{{ finishedCount }}</strong>
          <Activity :size="22" />
        </div>
      </div>

      <!-- Role: Loket View -->
      <section v-if="queueStore.role === 'loket'" class="dashboard-panel">
        <div class="panel-heading">
          <h3>Antrean hari ini</h3>
          <div style="display: flex; gap: 8px">
            <RouterLink
              to="/register"
              class="button outline small"
              @click="queueStore.startRegistration(queueStore.staffPoli || 'Poli Umum', true)"
            >
              + Pasien Walk-in
            </RouterLink>
            <RouterLink to="/scan" class="button primary small">
              <QrCode :size="17" /> Scan QR Pasien
            </RouterLink>
          </div>
        </div>
        <div class="search-field">
          <Search :size="18" />
          <input
            placeholder="Cari nama, NIK, atau nomor antrean..."
            v-model="queueStore.search"
          />
        </div>
        <div class="table-scroll">
          <table>
            <thead>
              <tr>
                <th>Antrean</th>
                <th>Pasien</th>
                <th>Poli</th>
                <th>Pendaftaran</th>
                <th>Status</th>
                <th>Tindakan</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in filteredLoketPatients" :key="p.id">
                <td>
                  <strong class="queue-number">
                    {{ p.number }}
                  </strong>
                </td>
                <td>
                  <strong>{{ p.name }}</strong>
                  <small>{{ p.nik }}</small>
                </td>
                <td>{{ p.poli }}</td>
                <td>
                  <span :class="['badge', { coral: p.source !== 'Online' }]">
                    {{ p.source }}
                  </span>
                </td>
                <td>
                  <span
                    :class="[
                      'status',
                      `status-${p.status === 'Selesai' ? 'done' : 'pending'}`,
                    ]"
                  >
                    {{ p.status }}
                  </span>
                </td>
                <td>
                  <button
                    v-if="
                      ['Terdaftar', 'Menunggu verifikasi'].includes(p.status)
                    "
                    class="text-button"
                    @click="queueStore.verify = p"
                  >
                    Verifikasi <ChevronRight :size="15" />
                  </button>
                  <button
                    v-else-if="p.status === 'Hadir'"
                    class="text-button"
                    @click="queueStore.call(p)"
                  >
                    Panggil <Volume2 :size="16" />
                  </button>
                  <button
                    v-else-if="p.status === 'Selesai' && (p.prescription || p.prescription_no)"
                    class="text-button"
                    :class="p.prescription?.status === 'dispensed' ? '!text-emerald-700 !font-bold' : '!text-amber-700 !font-bold'"
                    @click="openDispenseModal(p)"
                    title="Lihat & Tebus Resep Obat Pasien"
                  >
                    <Pill :size="14" />
                    <span>{{ p.prescription?.status === 'dispensed' ? 'Obat Selesai' : 'Tebus e-Resep' }}</span>
                  </button>
                  <span v-else class="muted">—</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <!-- Role: Medis / Poli View -->
      <template v-else-if="queueStore.role === 'medis'">
        <div class="poli-toolbar">
          <label>
            Poli Anda
            <select
              v-model="queueStore.staffPoli"
              @change="queueStore.selected = ''"
            >
              <option v-for="s in services" :key="s.name">
                {{ s.name }}
              </option>
            </select>
          </label>
        </div>
        <div class="medical-grid">
          <section class="dashboard-panel patient-list">
            <h3>Antrean {{ queueStore.staffPoli }}</h3>
            <button
              v-for="p in poliPatients"
              :key="p.id"
              :class="['patient', { selected: queueStore.selected === p.id }]"
              @click="selectPatient(p)"
            >
              <strong>
                {{ p.number }}
                <span>{{ p.status }}</span>
              </strong>
              <b>{{ p.name }}</b>
              <small> {{ p.insurance }} · {{ p.gender }} </small>
            </button>
            <p v-if="!poliPatients.length" class="muted">
              Belum ada pasien yang terverifikasi.
            </p>
          </section>

          <section class="dashboard-panel soap-panel">
            <div v-if="activePatient" :key="activePatient.id">
              <div class="record-banner">
                <div class="record-identity record-identity-compact">
                  <div class="record-patient-info">
                    <span class="record-avatar" aria-hidden="true">
                      {{ getInitials(activePatient.name) }}
                    </span>
                    <div class="record-patient-details">
                      <h3>{{ activePatient.name }}</h3>
                      <p>
                        No. RM: belum tersedia ·
                        <span class="badge">{{ activePatient.gender }}</span>
                      </p>
                      <p>
                        {{ formatPatientAge(activePatient.birth) }}
                      </p>
                      <p v-if="activePatient.insurance === 'BPJS'">
                        BPJS: {{ activePatient.bpjs || "Belum diisi" }}
                      </p>
                    </div>
                  </div>
                  <RouterLink
                    to="/medical-record"
                    class="button primary small record-trigger"
                    aria-label="Lihat Riwayat Rekam Medis"
                    title="Lihat Riwayat Rekam Medis"
                  >
                    <Activity :size="17" aria-hidden="true" />
                    <span class="record-trigger-label">MEDICAL RECORD</span>
                  </RouterLink>
                </div>
                <div
                  v-if="patientAllergies.length"
                  class="record-alert"
                  role="alert"
                >
                  PERINGATAN ALERGI: {{ patientAllergies.join(", ") }}
                </div>
              </div>

              <div class="record-layout">
                <div class="record-work">
                  <div class="summary-box">
                    <strong>Pra-anamnesis pasien</strong>
                    <p>{{ activePatient.complaint }}</p>
                    <small>
                      Gejala:
                      {{ activePatient.symptoms.join(", ") || "Tidak diisi" }} ·
                      Alergi: {{ activePatient.allergy || "Tidak ada" }}
                    </small>
                  </div>

                  <form @submit.prevent="saveMedicalRecord">
                    <h3 class="field-heading">Tanda vital</h3>
                    <div class="vital-grid">
                      <label v-for="field in vitalFields" :key="field.name">
                        {{ field.label }}
                        <input
                          required
                          type="number"
                          v-model="vitals[field.name]"
                          :min="field.min"
                          :max="field.max"
                          step="any"
                          :placeholder="field.placeholder"
                        />
                      </label>
                    </div>

                    <label>
                      Dokter pemeriksa
                      <input
                        required
                        v-model="doctorName"
                        placeholder="Nama lengkap dokter"
                        autocomplete="off"
                      />
                    </label>

                    <label>
                      S · Subjektif
                      <textarea required rows="2" v-model="subjectiveText" />
                    </label>

                    <label>
                      O · Objektif
                      <textarea required rows="2" v-model="objectiveText" />
                    </label>

                    <label>
                      A · Diagnosis ICD-10
                      <input
                        required
                        list="diagnoses-list"
                        v-model="queueStore.soap.diagnosis"
                        placeholder="Cari kode atau diagnosis..."
                      />
                      <datalist id="diagnoses-list">
                        <option value="J00 — Acute nasopharyngitis" />
                        <option
                          value="J06.9 — Infeksi saluran pernapasan atas"
                        />
                        <option value="R50.9 — Demam" />
                        <option value="I10 — Hipertensi esensial" />
                        <option value="E11 — Diabetes melitus tipe 2" />
                        <option value="K02.9 — Karies gigi" />
                        <option value="R51 — Sakit kepala" />
                      </datalist>
                    </label>

                    <label>
                      P · Rencana tindakan / terapi
                      <textarea
                        required
                        rows="2"
                        v-model="queueStore.soap.plan"
                      />
                    </label>

                    <!-- KOLOM BARU: PERENCANAAN PEMBERIAN OBAT & RESEP ELEKTRONIK -->
                    <div class="med-prescription-section">
                      <div class="med-header">
                        <div class="med-title-row">
                          <Pill :size="18" class="text-emerald-700" />
                          <div>
                            <strong>Perencanaan Terapi Obat & e-Resep</strong>
                            <p class="muted text-xs">Tentukan obat dan cara penggunaan terperinci (diminum, dikunyah, dioles, dll).</p>
                          </div>
                        </div>
                        <label class="toggle-rx-label">
                          <input type="checkbox" v-model="hasPrescription" />
                          <span>Terbitkan e-Resep</span>
                        </label>
                      </div>

                      <div v-if="hasPrescription" class="med-builder-box">
                        <!-- Dropdown Pilih Obat dari Database e-FORNAS -->
                        <div class="med-field-group">
                          <label>
                            Pilih Obat dari Database e-FORNAS Kemenkes
                            <select v-model="selectedMedId" @change="onSelectMedication">
                              <option value="">-- Pilih Obat (Nama Ukuran Serta Bentuknya Saja) --</option>
                              <option v-for="m in medications" :key="m.id" :value="m.id">
                                {{ formatMedDropdownLabel(m) }}
                              </option>
                            </select>
                          </label>
                        </div>

                        <!-- Grid Rute & Frekuensi -->
                        <div class="med-grid-2">
                          <label>
                            Cara Pakai / Rute Penggunaan
                            <select v-model="itemRoute">
                              <option value="oral_minum">Diminum (Oral - Ditelan utuh dengan air)</option>
                              <option value="oral_kunyah">Dikunyah (Oral Chewable - Wajib dikunyah lumat)</option>
                              <option value="topikal_oles">Dioleskan (Topikal - Salep/Krim/Gel)</option>
                              <option value="tetes_mata">Diteteskan ke Mata</option>
                              <option value="tetes_telinga">Diteteskan ke Telinga</option>
                              <option value="kumur">Dikumur (Gargle - Tidak untuk ditelan)</option>
                              <option value="inhalasi">Dihirup (Inhalasi)</option>
                            </select>
                          </label>

                          <label>
                            Aturan Pakai / Frekuensi
                            <input v-model="itemFrequency" placeholder="Contoh: 3 x 1 tablet sehari" />
                          </label>
                        </div>

                        <!-- Grid Waktu & Jumlah -->
                        <div class="med-grid-2">
                          <label>
                            Waktu Penggunaan
                            <select v-model="itemTiming">
                              <option value="Sesudah makan">Sesudah makan</option>
                              <option value="Sebelum makan (perut kosong)">Sebelum makan (perut kosong)</option>
                              <option value="1 jam sebelum makan atau saat maag kambuh">1 jam sebelum makan / saat maag kambuh</option>
                              <option value="Sebelum tidur malam">Sebelum tidur malam</option>
                              <option value="Bersama makanan">Bersama makanan</option>
                              <option value="Bila perlu (jika nyeri/demam)">Bila perlu (jika nyeri/demam)</option>
                              <option value="Pagi dan malam sesudah makan">Pagi dan malam sesudah makan</option>
                            </select>
                          </label>

                          <div class="med-qty-row">
                            <label style="flex: 1;">
                              Jumlah
                              <input type="number" min="1" v-model.number="itemQuantity" />
                            </label>
                            <label style="flex: 1.5;">
                              Satuan
                              <input v-model="itemUnit" placeholder="Tablet, Tube, Botol" />
                            </label>
                          </div>
                        </div>

                        <!-- Instruksi Detail Pasien -->
                        <div class="med-field-group">
                          <label>
                            Instruksi Cara Menggunakan Obat (Petunjuk Pasien)
                            <textarea rows="2" v-model="itemInstruction" placeholder="Contoh: Wajib dikunyah hingga lumat di mulut sebelum ditelan, lalu minum sedikit air." />
                          </label>
                        </div>

                        <div class="med-action-row">
                          <button type="button" class="button outline small" @click="addMedicationItem" :disabled="!selectedMedId && !customMedName">
                            <Plus :size="15" /> Tambahkan Obat ke Resep
                          </button>
                        </div>

                        <!-- Daftar Item Obat yang Telah Ditambahkan -->
                        <div v-if="prescriptionItems.length" class="rx-items-summary">
                          <div class="rx-summary-title">
                            <span>Daftar Obat dalam Resep ({{ prescriptionItems.length }} item):</span>
                          </div>
                          <div v-for="(item, idx) in prescriptionItems" :key="idx" class="rx-item-card">
                            <div class="rx-item-header">
                              <div class="rx-item-badges">
                                <span class="route-badge" :data-route="item.route">
                                  {{ formatRouteLabel(item.route) }}
                                </span>
                                <strong>{{ item.medication_name }}</strong>
                              </div>
                              <div class="rx-item-actions">
                                <span class="rx-item-qty">{{ item.quantity }} {{ item.unit }}</span>
                                <button type="button" class="btn-remove-med" @click="removeMedicationItem(idx)" title="Hapus obat">
                                  <Trash2 :size="13" />
                                </button>
                              </div>
                            </div>
                            <div class="rx-item-body">
                              <p><strong>Aturan:</strong> {{ item.frequency }} · {{ item.timing }}</p>
                              <p class="rx-instruction-text"><strong>Petunjuk:</strong> {{ item.instruction }}</p>
                            </div>
                          </div>
                        </div>
                        <p v-else class="muted text-xs text-center" style="padding: 10px 0;">
                          Belum ada obat yang dimasukkan ke resep. Pilih obat dari katalog di atas dan klik "Tambahkan Obat".
                        </p>
                      </div>
                    </div>

                    <button class="button primary full-width" type="submit">
                      <CheckCircle2 :size="18" /> Simpan & Tutup Kunjungan
                    </button>
                  </form>
                </div>
              </div>
            </div>
            <div v-else class="empty-state">
              <Stethoscope :size="44" />
              <h3>Pilih pasien untuk pemeriksaan</h3>
              <p>Pasien yang telah diverifikasi akan muncul di panel kiri.</p>
            </div>
          </section>
        </div>
      </template>

      <!-- ROLE: DATABASE MASTER PASIEN CRUD VIEW -->
      <section v-else-if="queueStore.role === 'database_pasien'" class="dashboard-panel crud-master-panel">
        <div class="panel-heading">
          <div>
            <h3>Database Master Pasien</h3>
            <p class="muted text-xs">Kelola data pasien klinik, NIK terenkripsi, penjamin (BPJS/Umum), dan riwayat medis.</p>
          </div>
          <button class="button primary small" @click="openAddPatientModal">
            <Plus :size="16" /> Tambah Pasien Baru
          </button>
        </div>

        <!-- Filter & Search Toolbar -->
        <div class="crud-toolbar">
          <div class="search-input-box">
            <Search :size="16" class="search-icon" />
            <input
              v-model="patientSearch"
              placeholder="Cari pasien berdasarkan Nama, NIK, No. HP, No. BPJS..."
              class="toolbar-search-input"
            />
          </div>
          <div class="toolbar-filter-group">
            <select v-model="patientInsuranceFilter" class="toolbar-select">
              <option value="all">Semua Penjamin</option>
              <option value="BPJS">BPJS Kesehatan</option>
              <option value="Umum">Pasien Umum</option>
            </select>
          </div>
        </div>

        <!-- Tabel Pasien -->
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th>No. RM / ID</th>
                <th>Nama Pasien & NIK</th>
                <th>Gender & Usia</th>
                <th>Kontak & Alamat</th>
                <th>Penjamin</th>
                <th>Alergi</th>
                <th style="text-align: right;">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="p in paginatedPatients" :key="p.id">
                <td>
                  <span class="font-mono font-bold text-xs text-[var(--deep)]">{{ p.id.startsWith('demo-') ? 'RM-2026-' + p.id.replace('demo-', '00') : p.id.slice(0, 8).toUpperCase() }}</span>
                </td>
                <td>
                  <strong class="block text-sm text-slate-800">{{ p.name }}</strong>
                  <span class="text-xs font-mono text-slate-500">{{ p.nik || "NIK belum diisi" }}</span>
                </td>
                <td>
                  <span class="text-xs">{{ p.gender }} · {{ formatPatientAge(p.birth) }}</span>
                </td>
                <td>
                  <div class="text-xs text-slate-600">
                    <p class="font-medium">{{ p.phone || "—" }}</p>
                    <p class="text-slate-400 truncate max-w-[200px]">{{ p.address || "Alamat belum diisi" }}</p>
                  </div>
                </td>
                <td>
                  <span :class="['badge', p.insurance === 'BPJS' ? '!bg-emerald-100 !text-emerald-800' : '!bg-blue-100 !text-blue-800']">
                    {{ p.insurance }} {{ p.bpjs ? `(${p.bpjs})` : '' }}
                  </span>
                </td>
                <td>
                  <span v-if="p.allergy && !/^(tidak ada|tidak|-)$/i.test(p.allergy)" class="badge !bg-rose-100 !text-rose-800 font-bold">
                    {{ p.allergy }}
                  </span>
                  <span v-else class="text-slate-400 text-xs">—</span>
                </td>
                <td style="text-align: right;">
                  <div class="action-btn-group">
                    <button class="action-btn edit" @click="openEditPatientModal(p)" title="Edit Data Pasien">
                      <Edit2 :size="15" />
                    </button>
                    <button class="action-btn view" @click="viewPatientRecord(p)" title="Lihat Rekam Medis">
                      <Activity :size="15" />
                    </button>
                    <button class="action-btn delete" @click="confirmDeletePatient(p)" title="Hapus Pasien">
                      <Trash2 :size="15" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="masterPatients.length === 0">
                <td colspan="7" class="text-center p-6 text-slate-400 italic">
                  Tidak ada pasien ditemukan sesuai pencarian atau filter.
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Controls Pasien -->
        <div class="crud-pagination-bar" v-if="masterPatients.length > 0">
          <div class="pagination-info">
            <span>
              Menampilkan
              <strong>{{ (patientCurrentPage - 1) * patientPerPage + 1 }}</strong>
              -
              <strong>{{ Math.min(patientCurrentPage * patientPerPage, masterPatients.length) }}</strong>
              dari
              <strong>{{ masterPatients.length }}</strong>
              pasien
            </span>
            <div class="per-page-selector">
              <label>Per halaman:</label>
              <select v-model.number="patientPerPage" class="per-page-select">
                <option :value="5">5</option>
                <option :value="10">10</option>
                <option :value="15">15</option>
                <option :value="20">20</option>
              </select>
            </div>
          </div>

          <div class="pagination-controls" v-if="patientTotalPages > 1">
            <button
              type="button"
              class="page-btn nav-btn"
              :disabled="patientCurrentPage === 1"
              @click="patientCurrentPage = 1"
              title="Halaman Pertama"
            >
              <ChevronsLeft :size="16" />
            </button>
            <button
              type="button"
              class="page-btn nav-btn"
              :disabled="patientCurrentPage === 1"
              @click="patientCurrentPage--"
              title="Halaman Sebelumnya"
            >
              <ChevronLeft :size="16" />
            </button>

            <div class="page-numbers">
              <template v-for="(pNum, idx) in getVisiblePages(patientCurrentPage, patientTotalPages)" :key="idx">
                <span v-if="pNum === '...'" class="page-ellipsis">...</span>
                <button
                  v-else
                  type="button"
                  class="page-btn page-number-btn"
                  :class="{ active: patientCurrentPage === pNum }"
                  @click="patientCurrentPage = Number(pNum)"
                >
                  {{ pNum }}
                </button>
              </template>
            </div>

            <button
              type="button"
              class="page-btn nav-btn"
              :disabled="patientCurrentPage === patientTotalPages"
              @click="patientCurrentPage++"
              title="Halaman Berikutnya"
            >
              <ChevronRight :size="16" />
            </button>
            <button
              type="button"
              class="page-btn nav-btn"
              :disabled="patientCurrentPage === patientTotalPages"
              @click="patientCurrentPage = patientTotalPages"
              title="Halaman Terakhir"
            >
              <ChevronsRight :size="16" />
            </button>
          </div>
        </div>
      </section>

      <!-- ROLE: DATABASE OBAT & TERAPI CRUD VIEW -->
      <section v-else-if="queueStore.role === 'database_obat'" class="dashboard-panel crud-master-panel">
        <div class="panel-heading">
          <div>
            <h3>Database Katalog Obat & Farmasi</h3>
            <p class="muted text-xs">Kelola master sediaan obat, rute penggunaan, aturan pakai, dan indikasi terapi.</p>
          </div>
          <button class="button primary small" @click="openAddMedModal">
            <Plus :size="16" /> Tambah Obat Baru
          </button>
        </div>

        <!-- Filter & Search Toolbar -->
        <div class="crud-toolbar">
          <div class="search-input-box">
            <Search :size="16" class="search-icon" />
            <input
              v-model="medSearch"
              placeholder="Cari obat berdasarkan Nama, Kode, Nama Generik, atau Indikasi..."
              class="toolbar-search-input"
            />
          </div>
          <div class="toolbar-filter-group">
            <select v-model="medRouteFilter" class="toolbar-select">
              <option value="all">Semua Cara Pakai</option>
              <option value="oral_minum">Diminum</option>
              <option value="oral_kunyah">Dikunyah</option>
              <option value="topikal_oles">Dioleskan</option>
              <option value="tetes_mata">Diteteskan ke Mata</option>
              <option value="tetes_telinga">Diteteskan ke Telinga</option>
              <option value="kumur">Dikumur</option>
              <option value="inhalasi">Dihirup</option>
              <option value="sirup_minum">Sirup</option>
            </select>
          </div>
        </div>

        <!-- Tabel Obat -->
        <div class="table-container">
          <table class="data-table">
            <thead>
              <tr>
                <th style="width: 80px;">Kode</th>
                <th>Nama Obat & Sediaan</th>
                <th style="width: 130px;">Cara Pakai</th>
                <th style="width: 160px; max-width: 180px;">Indikasi</th>
                <th style="width: 170px;">Aturan & Waktu Pakai</th>
                <th style="text-align: right; width: 90px;">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="m in paginatedMedications" :key="m.id">
                <td>
                  <span class="font-mono font-bold text-xs text-teal-800 bg-teal-50 px-2 py-0.5 rounded border border-teal-200">{{ m.code || 'MED' }}</span>
                </td>
                <td>
                  <strong class="block text-sm text-slate-800">{{ m.name }} {{ m.strength || '' }}</strong>
                  <span class="text-xs text-slate-500">{{ m.generic_name || m.name }} · {{ m.dosage_form }}</span>
                </td>
                <td>
                  <span class="route-badge" :data-route="m.route">
                    {{ formatRouteLabel(m.route) }}
                  </span>
                </td>
                <td style="width: 160px; max-width: 180px;">
                  <span class="text-xs text-slate-700 font-medium line-clamp-2" :title="m.indication || 'Terapi umum'">{{ m.indication || "Terapi umum" }}</span>
                </td>
                <td>
                  <div class="text-xs text-slate-600">
                    <p class="font-semibold text-slate-700">{{ m.default_frequency || "3 x 1 sehari" }}</p>
                    <p class="text-slate-500">{{ m.default_timing || "Sesudah makan" }}</p>
                  </div>
                </td>
                <td style="text-align: right;">
                  <div class="action-btn-group">
                    <button class="action-btn edit" @click="openEditMedModal(m)" title="Edit Obat">
                      <Edit2 :size="15" />
                    </button>
                    <button class="action-btn delete" @click="confirmDeleteMedication(m)" title="Hapus Obat">
                      <Trash2 :size="15" />
                    </button>
                  </div>
                </td>
              </tr>
              <tr v-if="filteredMedications.length === 0">
                <td colspan="6" class="text-center p-6 text-slate-400 italic">
                  Tidak ada obat ditemukan sesuai pencarian atau filter rute.
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Controls Obat -->
        <div class="crud-pagination-bar" v-if="filteredMedications.length > 0">
          <div class="pagination-info">
            <span>
              Menampilkan
              <strong>{{ (medCurrentPage - 1) * medPerPage + 1 }}</strong>
              -
              <strong>{{ Math.min(medCurrentPage * medPerPage, filteredMedications.length) }}</strong>
              dari
              <strong>{{ filteredMedications.length }}</strong>
              obat
            </span>
            <div class="per-page-selector">
              <label>Per halaman:</label>
              <select v-model.number="medPerPage" class="per-page-select">
                <option :value="5">5</option>
                <option :value="10">10</option>
                <option :value="15">15</option>
                <option :value="20">20</option>
                <option :value="30">30</option>
              </select>
            </div>
          </div>

          <div class="pagination-controls" v-if="medTotalPages > 1">
            <button
              type="button"
              class="page-btn nav-btn"
              :disabled="medCurrentPage === 1"
              @click="medCurrentPage = 1"
              title="Halaman Pertama"
            >
              <ChevronsLeft :size="16" />
            </button>
            <button
              type="button"
              class="page-btn nav-btn"
              :disabled="medCurrentPage === 1"
              @click="medCurrentPage--"
              title="Halaman Sebelumnya"
            >
              <ChevronLeft :size="16" />
            </button>

            <div class="page-numbers">
              <template v-for="(pNum, idx) in getVisiblePages(medCurrentPage, medTotalPages)" :key="idx">
                <span v-if="pNum === '...'" class="page-ellipsis">...</span>
                <button
                  v-else
                  type="button"
                  class="page-btn page-number-btn"
                  :class="{ active: medCurrentPage === pNum }"
                  @click="medCurrentPage = Number(pNum)"
                >
                  {{ pNum }}
                </button>
              </template>
            </div>

            <button
              type="button"
              class="page-btn nav-btn"
              :disabled="medCurrentPage === medTotalPages"
              @click="medCurrentPage++"
              title="Halaman Berikutnya"
            >
              <ChevronRight :size="16" />
            </button>
            <button
              type="button"
              class="page-btn nav-btn"
              :disabled="medCurrentPage === medTotalPages"
              @click="medCurrentPage = medTotalPages"
              title="Halaman Terakhir"
            >
              <ChevronsRight :size="16" />
            </button>
          </div>
        </div>
      </section>
    </div>

    <!-- Modals -->
    <VerifyModal />

    <!-- MODAL PENEBUSAN RESEP OBAT DI LOKET FARMASI -->
    <div v-if="isDispenseModalOpen && dispenseModalPatient" class="modal-backdrop" @click.self="closeDispenseModal">
      <div class="modal-card dispense-modal">
        <div class="modal-header">
          <div class="flex items-center gap-2">
            <Pill :size="20" class="text-teal-700" />
            <h3 class="font-bold text-base">Penebusan Resep Obat (e-Resep)</h3>
          </div>
          <button type="button" class="btn-close" @click="closeDispenseModal"><X :size="18" /></button>
        </div>

        <div class="modal-body space-y-4">
          <div class="patient-rx-meta">
            <div>
              <span class="text-xs text-slate-500 uppercase font-semibold">Nama Pasien</span>
              <p class="font-bold text-sm text-[var(--deep)]">{{ dispenseModalPatient.name }}</p>
            </div>
            <div>
              <span class="text-xs text-slate-500 uppercase font-semibold">Nomor Antrean / Poli</span>
              <p class="font-bold text-sm text-slate-700">{{ dispenseModalPatient.number }} · {{ dispenseModalPatient.poli }}</p>
            </div>
            <div>
              <span class="text-xs text-slate-500 uppercase font-semibold">Nomor e-Resep</span>
              <p class="font-mono font-black text-sm text-emerald-800">{{ dispenseModalPatient.prescription?.prescription_no || dispenseModalPatient.prescription_no || 'RX-TERCETAK' }}</p>
            </div>
            <div>
              <span class="text-xs text-slate-500 uppercase font-semibold">Status Penebusan</span>
              <p>
                <span class="badge" :class="dispenseModalPatient.prescription?.status === 'dispensed' ? '!bg-emerald-100 !text-emerald-800' : '!bg-amber-100 !text-amber-800'">
                  {{ dispenseModalPatient.prescription?.status === 'dispensed' ? 'Sudah Ditebus' : 'Menunggu Penebusan' }}
                </span>
              </p>
            </div>
          </div>

          <div class="rx-items-list-box">
            <h4 class="text-xs font-bold text-slate-700 uppercase mb-2">Daftar Obat yang Harus Diserahkan:</h4>
            <div v-if="dispenseModalPatient.prescription?.items?.length" class="space-y-2">
              <div
                v-for="(item, idx) in dispenseModalPatient.prescription.items"
                :key="idx"
                class="border border-slate-200 rounded-xl p-3 bg-slate-50/80 space-y-1.5"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2">
                    <span class="route-badge" :data-route="item.route">
                      {{ formatRouteLabel(item.route) }}
                    </span>
                    <strong class="text-sm text-slate-800">{{ item.medication_name }}</strong>
                  </div>
                  <span class="font-bold text-sm font-mono text-[var(--deep)]">{{ item.quantity }} {{ item.unit }}</span>
                </div>
                <div class="text-xs text-slate-600">
                  <p><strong>Aturan Minum:</strong> {{ item.frequency }} · {{ item.timing }}</p>
                  <p class="text-slate-500 mt-0.5"><strong>Petunjuk Pasien:</strong> {{ item.instruction }}</p>
                </div>
              </div>
            </div>
            <p v-else class="text-xs text-slate-500 italic p-3 text-center">
              Detail resep tersimpan dalam catatan kunjungan.
            </p>
          </div>
        </div>

        <div class="modal-footer flex items-center justify-end gap-2 pt-3 border-t border-[var(--line)]">
          <button type="button" class="button outline small" @click="closeDispenseModal">Tutup</button>
          <button
            v-if="dispenseModalPatient.prescription?.status !== 'dispensed'"
            type="button"
            class="button primary small"
            @click="confirmDispensePrescription"
          >
            <CheckCircle2 :size="16" /> Tandai Obat Selesai Ditebus & Diserahkan
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL TAMBAH / EDIT PASIEN -->
    <div v-if="isPatientModalOpen" class="modal-backdrop" @click.self="closePatientModal">
      <div class="modal-card crud-modal">
        <div class="modal-header">
          <div class="flex items-center gap-2">
            <Database :size="20" class="text-teal-700" />
            <h3 class="font-bold text-base">{{ patientModalMode === 'add' ? 'Tambah Pasien Baru' : 'Edit Data Pasien' }}</h3>
          </div>
          <button type="button" class="btn-close" @click="closePatientModal"><X :size="18" /></button>
        </div>

        <form @submit.prevent="savePatientSubmit" class="modal-body space-y-3">
          <div class="crud-form-grid-2">
            <label>
              Nomor Induk Kependudukan (NIK) *
              <input required v-model="patientForm.nik" maxlength="16" placeholder="16 digit NIK" />
            </label>
            <label>
              Nama Lengkap Pasien *
              <input required v-model="patientForm.name" placeholder="Nama lengkap sesuai KTP" />
            </label>
          </div>

          <div class="crud-form-grid-2">
            <label>
              Tanggal Lahir *
              <input required type="date" v-model="patientForm.birth" />
            </label>
            <label>
              Jenis Kelamin *
              <select required v-model="patientForm.gender">
                <option value="Laki-laki">Laki-laki</option>
                <option value="Perempuan">Perempuan</option>
              </select>
            </label>
          </div>

          <div class="crud-form-grid-2">
            <label>
              Nomor Telepon / WhatsApp
              <input v-model="patientForm.phone" placeholder="08xxxxxxxxxx" />
            </label>
            <label>
              Penjamin / Asuransi *
              <select required v-model="patientForm.insurance">
                <option value="Umum">Pasien Umum</option>
                <option value="BPJS">BPJS Kesehatan</option>
              </select>
            </label>
          </div>

          <div v-if="patientForm.insurance === 'BPJS'" class="crud-form-group">
            <label>
              Nomor Kartu BPJS *
              <input required v-model="patientForm.bpjs" placeholder="13 digit Nomor Kartu BPJS" />
            </label>
          </div>

          <div class="crud-form-group">
            <label>
              Alamat Domisili Lengkap
              <textarea rows="2" v-model="patientForm.address" placeholder="Jalan, RT/RW, Kelurahan, Kecamatan..." />
            </label>
          </div>

          <div class="crud-form-group">
            <label>
              Riwayat Alergi Obat / Makanan (Bila Ada)
              <input v-model="patientForm.allergy" placeholder="Contoh: Alergi Penisilin, Alergi Paracetamol, dll." />
            </label>
          </div>

          <div class="modal-footer flex items-center justify-end gap-2 pt-3 border-t border-[var(--line)]">
            <button type="button" class="button outline small" @click="closePatientModal">Batal</button>
            <button type="submit" class="button primary small">
              <CheckCircle2 :size="16" /> {{ patientModalMode === 'add' ? 'Simpan Pasien' : 'Perbarui Pasien' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL TAMBAH / EDIT OBAT -->
    <div v-if="isMedModalOpen" class="modal-backdrop" @click.self="closeMedModal">
      <div class="modal-card crud-modal">
        <div class="modal-header">
          <div class="flex items-center gap-2">
            <Pill :size="20" class="text-teal-700" />
            <h3 class="font-bold text-base">{{ medModalMode === 'add' ? 'Tambah Obat Baru ke Katalog' : 'Edit Data Obat' }}</h3>
          </div>
          <button type="button" class="btn-close" @click="closeMedModal"><X :size="18" /></button>
        </div>

        <form @submit.prevent="saveMedicationSubmit" class="modal-body space-y-3">
          <div class="crud-form-grid-2">
            <label>
              Kode Obat (e-FORNAS / Katalog)
              <input v-model="medForm.code" placeholder="Contoh: FORNAS-PCT-500" />
            </label>
            <label>
              Nama Obat (Generik Resmi) *
              <input required v-model="medForm.name" placeholder="Contoh: Parasetamol" />
            </label>
          </div>

          <div class="crud-form-grid-2">
            <label>
              Ukuran / Kekuatan Sediaan *
              <input required v-model="medForm.strength" placeholder="Contoh: 500 mg, 120 mg/5 ml, 1%" />
            </label>
            <label>
              Bentuk Sediaan *
              <select required v-model="medForm.dosage_form">
                <option value="Tablet">Tablet</option>
                <option value="Tablet Kunyah">Tablet Kunyah</option>
                <option value="Tablet Dispersibel">Tablet Dispersibel</option>
                <option value="Kapsul">Kapsul</option>
                <option value="Kaplet">Kaplet</option>
                <option value="Salep Kulit">Salep Kulit</option>
                <option value="Krim">Krim</option>
                <option value="Gel">Gel</option>
                <option value="Tetes Mata">Tetes Mata</option>
                <option value="Tetes Telinga">Tetes Telinga</option>
                <option value="Larutan Kumur">Larutan Kumur</option>
                <option value="Inhaler">Inhaler Aerosol</option>
                <option value="Sirup">Sirup</option>
                <option value="Sirup Kering">Sirup Kering</option>
                <option value="Suspensi">Suspensi</option>
                <option value="Serbuk Oral">Serbuk Oral / Sachet</option>
              </select>
            </label>
          </div>

          <div class="crud-form-grid-2">
            <label>
              Rute / Cara Penggunaan Obat *
              <select required v-model="medForm.route">
                <option value="oral_minum">Diminum</option>
                <option value="oral_kunyah">Dikunyah</option>
                <option value="topikal_oles">Dioleskan</option>
                <option value="tetes_mata">Diteteskan ke Mata</option>
                <option value="tetes_telinga">Diteteskan ke Telinga</option>
                <option value="kumur">Dikumur</option>
                <option value="inhalasi">Dihirup</option>
                <option value="sirup_minum">Sirup</option>
              </select>
            </label>
            <label>
              Indikasi *
              <input required v-model="medForm.indication" placeholder="Contoh: Maag, asam lambung berlebih, gastritis" />
            </label>
          </div>

          <div class="crud-form-grid-2">
            <label>
              Aturan Pakai / Frekuensi Default *
              <input required v-model="medForm.default_frequency" placeholder="Contoh: 3 x 1 tablet kunyah sehari" />
            </label>
            <label>
              Waktu Penggunaan Default *
              <input required v-model="medForm.default_timing" placeholder="Contoh: 1 jam sebelum makan atau saat maag kambuh" />
            </label>
          </div>

          <div class="crud-form-group">
            <label>
              Petunjuk Detail Cara Pemakaian Pasien *
              <textarea
                required
                rows="2"
                v-model="medForm.default_instruction"
                placeholder="Contoh: Wajib dikunyah hingga lumat di mulut sebelum ditelan, lalu minum sedikit air hangat."
              />
            </label>
          </div>

          <div class="crud-form-grid-2">
            <label>
              Jumlah Stok Tersedia *
              <input required type="number" min="0" v-model.number="medForm.stock" placeholder="100" />
            </label>
            <label>
              Satuan Stok *
              <input required v-model="medForm.unit" placeholder="Tablet, Tube, Botol, Canister" />
            </label>
          </div>

          <div class="modal-footer flex items-center justify-end gap-2 pt-3 border-t border-[var(--line)]">
            <button type="button" class="button outline small" @click="closeMedModal">Batal</button>
            <button type="submit" class="button primary small">
              <CheckCircle2 :size="16" /> {{ medModalMode === 'add' ? 'Tambahkan ke Katalog' : 'Simpan Perubahan' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import {
  Menu,
  Users,
  Stethoscope,
  Monitor,
  LogOut,
  Bell,
  Activity,
  QrCode,
  Search,
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
  Volume2,
  CheckCircle2,
  Pill,
  Plus,
  Trash2,
  X,
  Database,
  Edit2,
} from "lucide-vue-next";
import { services } from "../../constants/services";
import { todayDate, initialMedications } from "../../constants/initialData";
import { useQueueStore } from "../../stores/queue";
import type { Patient, Medication, PrescriptionItem, Prescription } from "../../types";
import VerifyModal from "./VerifyModal.vue";
import StaffSidebar from "./StaffSidebar.vue";
import { api } from "../../services/api";

const router = useRouter();
const queueStore = useQueueStore();

const currentUser = api.getCurrentUser();
const doctorName = ref(
  currentUser?.full_name && queueStore.role === "medis"
    ? currentUser.full_name
    : "dr. Maya Putri"
);

async function handleLogout() {
  await api.logout();
  router.push("/login");
  window.scrollTo(0, 0);
}

const sidebarCollapsed = ref(false);
const subjectiveText = ref("");
const objectiveText = ref("");

// ── State Perencanaan Obat & e-Resep ──
const medications = ref<Medication[]>(initialMedications);
const hasPrescription = ref(true);
const prescriptionItems = ref<PrescriptionItem[]>([]);
const selectedMedId = ref("");
const customMedName = ref("");
const itemRoute = ref("oral_minum");
const itemDosageForm = ref("Tablet");
const itemFrequency = ref("3 x 1 tablet sehari");
const itemTiming = ref("Sesudah makan");
const itemQuantity = ref(10);
const itemUnit = ref("Tablet");
const itemInstruction = ref("Diminum utuh dengan 1 gelas air putih sesudah makan.");
const itemNotes = ref("");

// ── State Penebusan Obat di Loket ──
const isDispenseModalOpen = ref(false);
const dispenseModalPatient = ref<Patient | null>(null);

// ── State Database Pasien CRUD & Pagination ──
const patientSearch = ref("");
const patientInsuranceFilter = ref("all");
const isPatientModalOpen = ref(false);
const patientModalMode = ref<"add" | "edit">("add");
const patientCurrentPage = ref(1);
const patientPerPage = ref(10);
const patientForm = reactive({
  id: "",
  nik: "",
  name: "",
  birth: "1990-01-01",
  gender: "Laki-laki",
  phone: "",
  address: "",
  insurance: "Umum",
  bpjs: "",
  allergy: "",
});

const masterPatients = computed(() => {
  return queueStore.patients.filter((p) => {
    const q = patientSearch.value.toLowerCase();
    const matchSearch =
      !q ||
      `${p.name} ${p.nik} ${p.phone} ${p.bpjs || ""}`.toLowerCase().includes(q);
    const matchInsurance =
      patientInsuranceFilter.value === "all" ||
      p.insurance === patientInsuranceFilter.value;
    return matchSearch && matchInsurance;
  });
});

const patientTotalPages = computed(() => {
  return Math.max(1, Math.ceil(masterPatients.value.length / patientPerPage.value));
});

const paginatedPatients = computed(() => {
  const start = (patientCurrentPage.value - 1) * patientPerPage.value;
  return masterPatients.value.slice(start, start + patientPerPage.value);
});

watch([patientSearch, patientInsuranceFilter, patientPerPage], () => {
  patientCurrentPage.value = 1;
});

// ── State Database Obat CRUD & Pagination ──
const medSearch = ref("");
const medRouteFilter = ref("all");
const isMedModalOpen = ref(false);
const medModalMode = ref<"add" | "edit">("add");
const medCurrentPage = ref(1);
const medPerPage = ref(10);
const medForm = reactive<Partial<Medication>>({
  id: "",
  code: "",
  name: "",
  strength: "",
  generic_name: "",
  dosage_form: "Tablet",
  route: "oral_minum",
  indication: "",
  default_frequency: "3 x 1 tablet sehari",
  default_timing: "Sesudah makan",
  usage_instruction: "",
  default_instruction: "",
  stock: 100,
  unit: "Tablet",
});

const filteredMedications = computed(() => {
  return medications.value.filter((m) => {
    const q = medSearch.value.toLowerCase();
    const matchSearch =
      !q ||
      `${m.name} ${m.generic_name || ""} ${m.code || ""} ${m.indication || ""}`
        .toLowerCase()
        .includes(q);
    const matchRoute =
      medRouteFilter.value === "all" || m.route === medRouteFilter.value;
    return matchSearch && matchRoute;
  });
});

const medTotalPages = computed(() => {
  return Math.max(1, Math.ceil(filteredMedications.value.length / medPerPage.value));
});

const paginatedMedications = computed(() => {
  const start = (medCurrentPage.value - 1) * medPerPage.value;
  return filteredMedications.value.slice(start, start + medPerPage.value);
});

watch([medSearch, medRouteFilter, medPerPage], () => {
  medCurrentPage.value = 1;
});

// Helper Page Numbers Generator with Ellipsis
function getVisiblePages(current: number, total: number): (number | string)[] {
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }
  if (current <= 4) {
    return [1, 2, 3, 4, 5, "...", total];
  }
  if (current >= total - 3) {
    return [1, "...", total - 4, total - 3, total - 2, total - 1, total];
  }
  return [1, "...", current - 1, current, current + 1, "...", total];
}

const vitals = reactive<Record<string, string>>({
  bp_sys: "120",
  bp_dia: "80",
  hr: "80",
  temp: "36.5",
  spo2: "98",
  weight: "60",
  height: "165",
});

const vitalFields = [
  {
    name: "bp_sys",
    label: "Sistolik (mmHg)",
    placeholder: "120",
    min: "40",
    max: "300",
  },
  {
    name: "bp_dia",
    label: "Diastolik (mmHg)",
    placeholder: "80",
    min: "20",
    max: "200",
  },
  {
    name: "hr",
    label: "Detak jantung (/menit)",
    placeholder: "80",
    min: "20",
    max: "300",
  },
  {
    name: "temp",
    label: "Suhu (°C)",
    placeholder: "36.5",
    min: "25",
    max: "45",
  },
  { name: "spo2", label: "SpO₂ (%)", placeholder: "98", min: "1", max: "100" },
  {
    name: "weight",
    label: "Berat badan (kg)",
    placeholder: "60",
    min: "0.1",
    max: "500",
  },
  {
    name: "height",
    label: "Tinggi badan (cm)",
    placeholder: "165",
    min: "20",
    max: "260",
  },
];

const currentFormattedDate = computed(() =>
  new Date().toLocaleDateString("id-ID", { dateStyle: "full" }),
);

const todayRegistrationsCount = computed(
  () => queueStore.patients.filter((p) => p.date === todayDate).length,
);

const pendingVerificationCount = computed(
  () =>
    queueStore.patients.filter(
      (p) =>
        p.date === todayDate &&
        ["Terdaftar", "Menunggu verifikasi"].includes(p.status),
    ).length,
);

const readyToServeCount = computed(
  () =>
    queueStore.patients.filter(
      (p) => p.date === todayDate && p.status === "Hadir",
    ).length,
);

const finishedCount = computed(
  () =>
    queueStore.patients.filter(
      (p) => p.date === todayDate && p.status === "Selesai",
    ).length,
);

const filteredLoketPatients = computed(() =>
  queueStore.patients.filter(
    (p) =>
      p.date === todayDate &&
      `${p.name} ${p.nik} ${p.number}`
        .toLowerCase()
        .includes(queueStore.search.toLowerCase()),
  ),
);

const poliPatients = computed(() =>
  queueStore.patients.filter(
    (p) =>
      p.poli === queueStore.staffPoli &&
      p.date === todayDate &&
      ["Hadir", "Dipanggil"].includes(p.status),
  ),
);

const activePatient = computed(() => {
  const p = queueStore.chosen;
  if (
    p &&
    p.poli === queueStore.staffPoli &&
    ["Hadir", "Dipanggil"].includes(p.status)
  ) {
    return p;
  }
  return null;
});

const patientAllergies = computed(() => {
  if (!activePatient.value) return [];
  const chosen = activePatient.value;
  const allergies = [
    ...new Set(
      queueStore.patients
        .filter(
          (p) => p.id === chosen.id || (chosen.nik && p.nik === chosen.nik),
        )
        .map((p) => p.allergy?.trim())
        .filter((a) => a && !/^(tidak ada|tidak|nihil|none|-)$/i.test(a)),
    ),
  ];
  return allergies;
});

function getInitials(name: string) {
  return name
    .split(" ")
    .slice(0, 2)
    .map((n) => n[0])
    .join("");
}

function formatPatientAge(birthStr: string) {
  if (!birthStr) return "Tanggal lahir belum tersedia";
  try {
    const birth = new Date(birthStr);
    const today = new Date();
    let age = today.getFullYear() - birth.getFullYear();
    const m = today.getMonth() - birth.getMonth();
    if (m < 0 || (m === 0 && today.getDate() < birth.getDate())) {
      age--;
    }
    return `${birth.toLocaleDateString("id-ID")} · ${Math.max(0, age)} tahun`;
  } catch {
    return "Tanggal lahir belum tersedia";
  }
}

function formatRouteLabel(route?: string) {
  if (!route) return "Diminum";
  const routeLabels: Record<string, string> = {
    oral_minum: "Diminum",
    oral_kunyah: "Dikunyah",
    topikal_oles: "Dioleskan",
    tetes_mata: "Tetes Mata",
    tetes_telinga: "Tetes Telinga",
    kumur: "Dikumur",
    inhalasi: "Dihirup",
    sirup_minum: "Sirup",
  };
  return routeLabels[route] || route;
}

function formatMedDropdownLabel(m: Medication): string {
  // Hanya menampilkan nama, ukuran, serta bentuknya aja sesuai standar e-FORNAS
  const name = m.name?.trim() || "";
  const strength = m.strength?.trim() || "";
  const form = m.dosage_form?.trim() || "";

  if (strength && !name.toLowerCase().includes(strength.toLowerCase())) {
    return `${name} ${strength} - ${form}`;
  }
  return `${name} - ${form}`;
}

function onSelectMedication() {
  if (!selectedMedId.value) return;
  const found = medications.value.find((m) => m.id === selectedMedId.value);
  if (found) {
    const strengthPart = found.strength && !found.name.toLowerCase().includes(found.strength.toLowerCase()) ? ` ${found.strength}` : "";
    customMedName.value = `${found.name}${strengthPart} - ${found.dosage_form}`;
    itemRoute.value = found.route || "oral_minum";
    itemDosageForm.value = found.dosage_form || "Tablet";
    itemFrequency.value = found.default_frequency || "3 x 1 tablet sehari";
    itemTiming.value = found.default_timing || "Sesudah makan";
    itemUnit.value = found.unit || "Tablet";
    itemInstruction.value = found.usage_instruction || found.default_instruction || "";
    itemQuantity.value = 10;
  }
}

function addMedicationItem() {
  const medName = customMedName.value.trim();
  if (!medName) return;

  const newItem: PrescriptionItem = {
    medication_id: selectedMedId.value || undefined,
    medication_name: medName,
    dosage_form: itemDosageForm.value,
    route: itemRoute.value,
    frequency: itemFrequency.value,
    timing: itemTiming.value,
    quantity: itemQuantity.value || 1,
    unit: itemUnit.value || "Tablet",
    instruction: itemInstruction.value,
    notes: itemNotes.value,
  };

  prescriptionItems.value.push(newItem);

  // Sync into plan
  const planLine = `${prescriptionItems.value.length}. ${newItem.medication_name} (${formatRouteLabel(newItem.route)}) - ${newItem.frequency} ${newItem.timing}. Petunjuk: ${newItem.instruction}`;
  if (!queueStore.soap.plan) {
    queueStore.soap.plan = `Resep Terapi Obat:\n${planLine}`;
  } else if (!queueStore.soap.plan.includes(newItem.medication_name)) {
    queueStore.soap.plan += `\n${planLine}`;
  }

  // Reset inputs
  selectedMedId.value = "";
  customMedName.value = "";
  itemInstruction.value = "";
  itemNotes.value = "";
}

function removeMedicationItem(index: number) {
  prescriptionItems.value.splice(index, 1);
}

function openDispenseModal(p: Patient) {
  dispenseModalPatient.value = p;
  isDispenseModalOpen.value = true;
}

function closeDispenseModal() {
  isDispenseModalOpen.value = false;
  dispenseModalPatient.value = null;
}

async function confirmDispensePrescription() {
  if (!dispenseModalPatient.value) return;
  const p = dispenseModalPatient.value;
  const rxId = p.prescription?.id;
  if (rxId) {
    try {
      await api.dispensePrescription(rxId);
    } catch {
      // offline fallback
    }
  }
  if (p.prescription) {
    p.prescription.status = "dispensed";
    p.prescription.dispensed_at = new Date().toISOString();
  }
  queueStore.setNotice(`Resep obat untuk ${p.name} berhasil ditebus dan diserahkan.`);
  closeDispenseModal();
}

function openAddPatientModal() {
  patientModalMode.value = "add";
  patientForm.id = "p-" + Date.now();
  patientForm.nik = "";
  patientForm.name = "";
  patientForm.birth = "1995-05-15";
  patientForm.gender = "Laki-laki";
  patientForm.phone = "";
  patientForm.address = "";
  patientForm.insurance = "Umum";
  patientForm.bpjs = "";
  patientForm.allergy = "";
  isPatientModalOpen.value = true;
}

function openEditPatientModal(p: Patient) {
  patientModalMode.value = "edit";
  patientForm.id = p.id;
  patientForm.nik = p.nik || "";
  patientForm.name = p.name;
  patientForm.birth = p.birth || "1990-01-01";
  patientForm.gender = p.gender || "Laki-laki";
  patientForm.phone = p.phone || "";
  patientForm.address = p.address || "";
  patientForm.insurance = p.insurance || "Umum";
  patientForm.bpjs = p.bpjs || "";
  patientForm.allergy = p.allergy || "";
  isPatientModalOpen.value = true;
}

function closePatientModal() {
  isPatientModalOpen.value = false;
}

async function savePatientSubmit() {
  if (patientModalMode.value === "add") {
    const newPatient: Patient = {
      id: patientForm.id || "p-" + Date.now(),
      name: patientForm.name,
      nik: patientForm.nik,
      birth: patientForm.birth,
      gender: patientForm.gender,
      phone: patientForm.phone,
      address: patientForm.address,
      insurance: patientForm.insurance as any,
      bpjs: patientForm.bpjs,
      allergy: patientForm.allergy,
      complaint: "Pemeriksaan kesehatan",
      symptoms: [],
      poli: "Poli Umum",
      date: todayDate,
      number: "P" + String(queueStore.patients.length + 1).padStart(3, "0"),
      status: "Terdaftar",
      source: "Loket",
      createdAt: new Date().toISOString(),
    };
    queueStore.patients.unshift(newPatient);
    queueStore.setNotice(`Pasien ${newPatient.name} berhasil ditambahkan ke database.`);
  } else {
    const idx = queueStore.patients.findIndex((p) => p.id === patientForm.id);
    if (idx !== -1) {
      queueStore.patients[idx] = {
        ...queueStore.patients[idx],
        name: patientForm.name,
        nik: patientForm.nik,
        birth: patientForm.birth,
        gender: patientForm.gender,
        phone: patientForm.phone,
        address: patientForm.address,
        insurance: patientForm.insurance as any,
        bpjs: patientForm.bpjs,
        allergy: patientForm.allergy,
      };
      queueStore.setNotice(`Data pasien ${patientForm.name} berhasil diperbarui.`);
    }
  }
  localStorage.setItem("antrisehat-patients", JSON.stringify(queueStore.patients));
  closePatientModal();
}

function confirmDeletePatient(p: Patient) {
  if (window.confirm(`Yakin ingin menghapus pasien "${p.name}" (${p.nik}) dari database?`)) {
    queueStore.patients = queueStore.patients.filter((item) => item.id !== p.id);
    localStorage.setItem("antrisehat-patients", JSON.stringify(queueStore.patients));
    queueStore.setNotice(`Pasien ${p.name} telah dihapus dari database.`);
  }
}

function viewPatientRecord(p: Patient) {
  queueStore.selected = p.id;
  router.push("/medical-record");
}

function openAddMedModal() {
  medModalMode.value = "add";
  const nextNum = medications.value.length + 1;
  medForm.id = "fornas-" + Date.now();
  medForm.code = `FORNAS-MED-${String(nextNum).padStart(3, "0")}`;
  medForm.name = "";
  medForm.strength = "";
  medForm.generic_name = "";
  medForm.dosage_form = "Tablet";
  medForm.route = "oral_minum";
  medForm.indication = "";
  medForm.default_frequency = "3 x 1 tablet sehari";
  medForm.default_timing = "Sesudah makan";
  medForm.usage_instruction = "";
  medForm.default_instruction = "";
  medForm.stock = 100;
  medForm.unit = "Tablet";
  isMedModalOpen.value = true;
}

function openEditMedModal(m: Medication) {
  medModalMode.value = "edit";
  medForm.id = m.id;
  medForm.code = m.code;
  medForm.name = m.name;
  medForm.strength = m.strength || "";
  medForm.generic_name = m.generic_name || "";
  medForm.dosage_form = m.dosage_form || "Tablet";
  medForm.route = m.route || "oral_minum";
  medForm.indication = m.indication || "";
  medForm.default_frequency = m.default_frequency || "3 x 1 tablet sehari";
  medForm.default_timing = m.default_timing || "Sesudah makan";
  medForm.usage_instruction = m.usage_instruction || m.default_instruction || "";
  medForm.default_instruction = m.default_instruction || m.usage_instruction || "";
  medForm.stock = m.stock ?? 100;
  medForm.unit = m.unit || "Tablet";
  isMedModalOpen.value = true;
}

function closeMedModal() {
  isMedModalOpen.value = false;
}

async function saveMedicationSubmit() {
  if (medModalMode.value === "add") {
    const newMed: Medication = {
      id: medForm.id || "fornas-" + Date.now(),
      code: medForm.code || `FORNAS-MED-${String(medications.value.length + 1).padStart(3, "0")}`,
      name: medForm.name!,
      strength: medForm.strength || "",
      generic_name: medForm.generic_name || "",
      category: "Katalog e-FORNAS Klinik",
      dosage_form: medForm.dosage_form || "Tablet",
      route: medForm.route || "oral_minum",
      indication: medForm.indication || "",
      default_frequency: medForm.default_frequency || "3 x 1 sehari",
      default_timing: medForm.default_timing || "Sesudah makan",
      usage_instruction: medForm.default_instruction || medForm.usage_instruction || "",
      default_instruction: medForm.default_instruction || "",
      stock: medForm.stock ?? 100,
      unit: medForm.unit || "Tablet",
      is_prescription_required: true,
    };
    medications.value.unshift(newMed);
    try {
      await api.createMedication(newMed);
    } catch {}
    queueStore.setNotice(`Obat ${newMed.name} ${newMed.strength} berhasil ditambahkan ke katalog.`);
  } else {
    const idx = medications.value.findIndex((m) => m.id === medForm.id);
    if (idx !== -1) {
      medications.value[idx] = {
        ...medications.value[idx],
        code: medForm.code || medications.value[idx].code,
        name: medForm.name!,
        strength: medForm.strength || "",
        generic_name: medForm.generic_name || "",
        dosage_form: medForm.dosage_form || "Tablet",
        route: medForm.route || "oral_minum",
        indication: medForm.indication || "",
        default_frequency: medForm.default_frequency || "3 x 1 sehari",
        default_timing: medForm.default_timing || "Sesudah makan",
        usage_instruction: medForm.default_instruction || medForm.usage_instruction || "",
        default_instruction: medForm.default_instruction || "",
        stock: medForm.stock ?? 100,
        unit: medForm.unit || "Tablet",
      };
      if (medForm.id) {
        try {
          await api.updateMedication(medForm.id, medications.value[idx]);
        } catch {}
      }
      queueStore.setNotice(`Data obat ${medForm.name} berhasil diperbarui.`);
    }
  }
  localStorage.setItem("antrisehat-medications-v2", JSON.stringify(medications.value));
  closeMedModal();
}

async function confirmDeleteMedication(m: Medication) {
  if (window.confirm(`Yakin ingin menghapus obat "${m.name}" (${m.code}) dari katalog?`)) {
    medications.value = medications.value.filter((item) => item.id !== m.id);
    if (m.id) {
      try {
        await api.deleteMedication(m.id);
      } catch {}
    }
    localStorage.setItem("antrisehat-medications-v2", JSON.stringify(medications.value));
    queueStore.setNotice(`Obat ${m.name} telah dihapus dari katalog.`);
  }
}

onMounted(async () => {
  if (typeof window !== "undefined" && window.innerWidth < 1024) {
    sidebarCollapsed.value = true;
  }
  try {
    const savedMeds = localStorage.getItem("antrisehat-medications-v2");
    if (savedMeds) {
      const parsed = JSON.parse(savedMeds);
      if (parsed && parsed.length) {
        medications.value = parsed;
      }
    } else {
      medications.value = initialMedications;
      localStorage.setItem("antrisehat-medications-v2", JSON.stringify(initialMedications));
    }
    const res = await api.getMedications();
    if (res && res.length) {
      medications.value = res;
      localStorage.setItem("antrisehat-medications-v2", JSON.stringify(res));
    }
  } catch {
    // fallback
  }
});

function selectPatient(p: Patient) {
  queueStore.selected = p.id;
  subjectiveText.value = [
    p.complaint,
    p.symptoms.length ? `Gejala: ${p.symptoms.join(", ")}` : "",
  ]
    .filter(Boolean)
    .join("\n");
  objectiveText.value = "";
  queueStore.soap.diagnosis = "";
  queueStore.soap.plan = "";

  // Reset prescription builder
  prescriptionItems.value = [];
  selectedMedId.value = "";
  customMedName.value = "";
  hasPrescription.value = true;

  // Auto-suggest initial medication based on complaint/symptoms
  const c = (p.complaint || "").toLowerCase();
  if (c.includes("maag") || c.includes("lambung") || c.includes("perih") || c.includes("ulu hati")) {
    const antasida = medications.value.find((m) => m.name.toLowerCase().includes("antasida"));
    if (antasida) {
      selectedMedId.value = antasida.id;
      onSelectMedication();
    }
  } else if (c.includes("gatal") || c.includes("kulit") || c.includes("eksim")) {
    const hydro = medications.value.find((m) => m.name.toLowerCase().includes("hidrokortison"));
    if (hydro) {
      selectedMedId.value = hydro.id;
      onSelectMedication();
    }
  } else if (c.includes("mata") || c.includes("merah") || c.includes("belekan")) {
    const eye = medications.value.find((m) => m.dosage_form === "Tetes Mata" || m.name.toLowerCase().includes("kloramfenikol"));
    if (eye) {
      selectedMedId.value = eye.id;
      onSelectMedication();
    }
  } else if (c.includes("gigi") || c.includes("ngilu")) {
    const mef = medications.value.find((m) => m.name.toLowerCase().includes("mefenamat"));
    if (mef) {
      selectedMedId.value = mef.id;
      onSelectMedication();
    }
  } else {
    const pct = medications.value.find((m) => m.name.toLowerCase().includes("parasetamol") || m.name.toLowerCase().includes("paracetamol"));
    if (pct) {
      selectedMedId.value = pct.id;
      onSelectMedication();
    }
  }
}

async function saveMedicalRecord() {
  if (!activePatient.value) return;
  const current = activePatient.value;

  // Build medication plan text
  let medPlanText = "";
  if (hasPrescription.value && prescriptionItems.value.length > 0) {
    medPlanText = prescriptionItems.value
      .map(
        (it, idx) =>
          `${idx + 1}. ${it.medication_name} [${formatRouteLabel(it.route)}] - ${it.frequency} (${it.timing}). Jumlah: ${it.quantity} ${it.unit}. Petunjuk: ${it.instruction}`
      )
      .join("\n");
  }

  // Generate Prescription object if enabled
  const prescriptionNo = hasPrescription.value && prescriptionItems.value.length > 0
    ? `RX-${new Date().toISOString().slice(0, 10).replace(/-/g, "")}-${current.number || "001"}`
    : undefined;

  const rxObject: Prescription | undefined = prescriptionNo
    ? {
        id: "rx-" + Date.now(),
        prescription_no: prescriptionNo,
        appointment_id: current.id,
        patient_id: current.id,
        patient_name: current.name,
        doctor_name: doctorName.value,
        polyclinic: current.poli,
        status: "pending_dispense",
        items: [...prescriptionItems.value],
        created_at: new Date().toISOString(),
      }
    : undefined;

  const record = {
    ...vitals,
    pressure: `${vitals.bp_sys}/${vitals.bp_dia}`,
    temperature: vitals.temp,
    pulse: vitals.hr,
    doctor: doctorName.value,
    subjective: subjectiveText.value,
    objective: objectiveText.value,
    diagnosis: queueStore.soap.diagnosis,
    plan: queueStore.soap.plan,
    medication_plan: medPlanText,
    prescription_no: prescriptionNo,
    prescription: rxObject,
    recordedAt: new Date().toISOString(),
  };

  queueStore.changeStatus(current.id, "Selesai");
  queueStore.patients = queueStore.patients.map((p) =>
    p.id === current.id
      ? {
          ...p,
          status: "Selesai",
          soap: JSON.stringify(record),
          medication_plan: medPlanText,
          prescription_no: prescriptionNo,
          prescription: rxObject,
        }
      : p,
  );

  // If patient has a valid UUID, persist visit to backend
  if (
    /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i.test(
      current.id,
    )
  ) {
    try {
      await api.createVisit({
        appointment_id: current.id,
        patient_id: current.id,
        doctor_id: "00000000-0000-0000-0000-000000000000",
        vitals: {
          bp_sys: Number(vitals.bp_sys) || 120,
          bp_dia: Number(vitals.bp_dia) || 80,
          temp: Number(vitals.temp) || 36.5,
          hr: Number(vitals.hr) || 80,
          spo2: Number(vitals.spo2) || 98,
          weight: Number(vitals.weight) || 65,
          height: Number(vitals.height) || 168,
        },
        subjective:
          subjectiveText.value || current.complaint || "Pemeriksaan rutin",
        objective:
          objectiveText.value || "Pemeriksaan fisik umum dalam batas normal",
        assessment: queueStore.soap.diagnosis || "Observasi Klinis",
        plan: queueStore.soap.plan || "Terapi simtomatik dan edukasi",
        medication_plan: medPlanText,
        prescription_no: prescriptionNo,
        prescription: rxObject
          ? {
              status: "pending_dispense",
              items: rxObject.items.map((it) => ({
                medication_id: it.medication_id,
                medication_name: it.medication_name,
                dosage_form: it.dosage_form,
                route: it.route,
                frequency: it.frequency,
                timing: it.timing,
                quantity: it.quantity,
                unit: it.unit,
                instruction: it.instruction,
                notes: it.notes,
              })),
            }
          : undefined,
        diagnoses: [
          {
            icd10_code: "Z00.0",
            is_primary: true,
            notes: queueStore.soap.diagnosis,
          },
        ],
      });
    } catch {
      // offline fallback
    }
  }

  queueStore.selected = "";
  queueStore.setNotice(
    prescriptionNo
      ? `Pemeriksaan disimpan. e-Resep (${prescriptionNo}) diterbitkan untuk ditebus di Farmasi.`
      : "Pemeriksaan disimpan dan kunjungan ditutup."
  );
}

function showPendingNotice() {
  queueStore.setNotice(
    `${pendingVerificationCount.value} pasien menunggu verifikasi`,
  );
}

function handleWalkinRegistration() {
  queueStore.startRegistration(queueStore.staffPoli || "Poli Umum", true);
  router.push("/register");
  window.scrollTo(0, 0);
}

function goToScan() {
  router.push("/scan");
  window.scrollTo(0, 0);
}

function goToTv() {
  router.push("/tv");
  window.scrollTo(0, 0);
}

function goToHome() {
  router.push("/");
  window.scrollTo(0, 0);
}

function goToMedicalRecord() {
  router.push("/medical-record");
  window.scrollTo(0, 0);
}
</script>

<style scoped>
.record-banner {
  margin-bottom: 24px;
}

.record-identity {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  width: 100%;
}

.record-patient-info {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.record-avatar {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background-color: #1e464c;
  color: #ffffff;
  font-weight: 700;
  font-size: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  letter-spacing: 0.5px;
}

.record-patient-details h3 {
  margin: 0 0 4px 0;
  font-size: 18px;
  font-weight: 700;
  color: #111827;
  line-height: 1.25;
}

.record-patient-details p {
  margin: 2px 0 0 0;
  font-size: 13px;
  color: #64748b;
  line-height: 1.4;
}

.record-trigger {
  margin-left: auto;
  align-self: center;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 18px;
  border-radius: 9999px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.5px;
  background-color: #1e464c;
  color: #ffffff;
  border: none;
  cursor: pointer;
  white-space: nowrap;
  transition: all 0.2s ease;
  box-shadow: 0 2px 6px rgba(30, 70, 76, 0.2);
}

.record-trigger:hover {
  background-color: #13353a;
  transform: translateY(-1px);
}

.record-alert {
  background-color: #fef2f2;
  border: 1px solid #fee2e2;
  color: #b91c1c;
  border-radius: 8px;
  padding: 10px 14px;
  font-size: 12px;
  font-weight: 600;
  margin-top: 14px;
}

@media (max-width: 640px) {
  .record-identity {
    flex-direction: column;
    align-items: stretch;
  }

  .record-trigger {
    margin-left: 0;
    align-self: stretch;
    justify-content: center;
    margin-top: 8px;
  }
}

/* ── Styling Medication Section & Builder ── */
.med-prescription-section {
  margin: 18px 0;
  padding: 16px;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 14px;
}

.med-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  flex-wrap: wrap;
  gap: 10px;
}

.med-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.toggle-rx-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 700;
  color: #195a5f;
  cursor: pointer;
  background: #e6fffa;
  padding: 6px 12px;
  border-radius: 20px;
  border: 1px solid #b2f5ea;
}

.med-builder-box {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.med-field-group label,
.med-grid-2 label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #334155;
}

.med-field-group select,
.med-field-group input,
.med-field-group textarea,
.med-grid-2 select,
.med-grid-2 input {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 13px;
  color: #0f172a;
  background: #f8fafc;
  outline: none;
  transition: border-color 0.2s;
}

.med-field-group select:focus,
.med-field-group input:focus,
.med-field-group textarea:focus,
.med-grid-2 select:focus,
.med-grid-2 input:focus {
  border-color: #195a5f;
  background: #ffffff;
}

.med-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.med-qty-row {
  display: flex;
  gap: 8px;
}

.med-action-row {
  display: flex;
  justify-content: flex-end;
  margin-top: 4px;
}

/* ── Route Badges ── */
.route-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.route-badge[data-route="oral_minum"],
.route-badge[data-route="sirup_minum"] {
  background: #d1fae5;
  color: #065f46;
  border: 1px solid #a7f3d0;
}

.route-badge[data-route="oral_kunyah"] {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fcd34d;
}

.route-badge[data-route="topikal_oles"] {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid #bfdbfe;
}

.route-badge[data-route="tetes_mata"],
.route-badge[data-route="tetes_telinga"] {
  background: #f3e8ff;
  color: #6b21a8;
  border: 1px solid #e9d5ff;
}

.route-badge[data-route="kumur"] {
  background: #ffedd5;
  color: #9a3412;
  border: 1px solid #fed7aa;
}

.route-badge[data-route="inhalasi"] {
  background: #e0f2fe;
  color: #0369a1;
  border: 1px solid #bae6fd;
}

/* ── Rx Items Summary ── */
.rx-items-summary {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rx-summary-title {
  font-size: 12px;
  font-weight: 700;
  color: #475569;
}

.rx-item-card {
  padding: 10px 14px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  font-size: 12px;
}

.rx-item-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6px;
}

.rx-item-badges {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rx-item-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.rx-item-qty {
  font-weight: 700;
  font-family: monospace;
  color: #0f172a;
}

.btn-remove-med {
  background: transparent;
  border: none;
  color: #ef4444;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
}

.btn-remove-med:hover {
  background: #fee2e2;
}

.rx-item-body {
  color: #475569;
  line-height: 1.4;
}

.rx-instruction-text {
  color: #0f766e;
  margin-top: 2px;
}

/* ── Dispense Modal ── */
.dispense-modal {
  max-width: 580px;
  width: 90%;
  background: #ffffff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.patient-rx-meta {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
  padding: 12px;
  background: #f1f5f9;
  border-radius: 10px;
}

.rx-items-list-box {
  max-height: 280px;
  overflow-y: auto;
}

/* ── CRUD Master Panel Styles ── */
.crud-master-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.crud-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
  padding: 12px 16px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
}

.search-input-box {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 260px;
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  padding: 8px 12px;
}

.search-input-box .search-icon {
  color: #64748b;
  flex-shrink: 0;
}

.toolbar-search-input {
  border: none;
  outline: none;
  background: transparent;
  width: 100%;
  font-size: 13px;
  color: #0f172a;
}

.toolbar-filter-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.toolbar-select {
  padding: 8px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  background: #ffffff;
  font-size: 13px;
  font-weight: 500;
  color: #334155;
  outline: none;
  cursor: pointer;
}

.action-btn-group {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn.edit {
  color: #0f766e;
}
.action-btn.edit:hover {
  background: #e6fffa;
  border-color: #99f6e4;
}

.action-btn.view {
  color: #2563eb;
}
.action-btn.view:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
}

.action-btn.delete {
  color: #ef4444;
}
.action-btn.delete:hover {
  background: #fef2f2;
  border-color: #fecaca;
}

/* ── Modal CRUD Styles ── */
.crud-modal {
  max-width: 650px;
  width: 90%;
  background: #ffffff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 22px;
  border-bottom: 1px solid #e2e8f0;
  background: #ffffff;
}

.modal-header h3 {
  font-size: 16px;
  font-weight: 700;
  color: #0f172a;
  margin: 0;
}

.btn-close {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: #f8fafc;
  color: #64748b;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-close:hover {
  background: #fee2e2;
  color: #ef4444;
  border-color: #fca5a5;
}

.modal-body {
  padding: 20px 22px;
}

.modal-footer {
  padding: 14px 22px;
  background: #f8fafc;
  border-top: 1px solid #e2e8f0;
}

.crud-form-grid-2 {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.crud-form-grid-2 label,
.crud-form-group label {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  color: #334155;
}

.crud-form-grid-2 input,
.crud-form-grid-2 select,
.crud-form-group input,
.crud-form-group select,
.crud-form-group textarea {
  padding: 8px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 8px;
  font-size: 13px;
  color: #0f172a;
  background: #f8fafc;
  outline: none;
  transition: border-color 0.2s;
}

.crud-form-grid-2 input:focus,
.crud-form-grid-2 select:focus,
.crud-form-group input:focus,
.crud-form-group select:focus,
.crud-form-group textarea:focus {
  border-color: #195a5f;
  background: #ffffff;
}

@media (max-width: 640px) {
  .crud-form-grid-2 {
    grid-template-columns: 1fr;
  }
}

/* ── CRUD Pagination Bar ── */
.crud-pagination-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
  padding: 14px 18px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  margin-top: 12px;
}

.pagination-info {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 13px;
  color: #475569;
  flex-wrap: wrap;
}

.per-page-selector {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12.5px;
  color: #64748b;
}

.per-page-select {
  padding: 4px 8px;
  border: 1px solid #cbd5e1;
  border-radius: 6px;
  background: #f8fafc;
  font-size: 12.5px;
  font-weight: 600;
  color: #0f172a;
  outline: none;
  cursor: pointer;
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: 6px;
}

.page-numbers {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 34px;
  height: 34px;
  padding: 0 8px;
  border-radius: 8px;
  border: 1px solid #e2e8f0;
  background: #ffffff;
  color: #334155;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.18s ease;
  user-select: none;
}

.page-btn:hover:not(:disabled) {
  background: #f1f5f9;
  border-color: #cbd5e1;
  color: #0f172a;
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: #f8fafc;
}

.page-number-btn.active {
  background: #195a5f;
  border-color: #195a5f;
  color: #ffffff;
  box-shadow: 0 2px 4px rgba(25, 90, 95, 0.25);
}

.page-ellipsis {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  font-size: 13px;
  color: #94a3b8;
  user-select: none;
}

.sidebar-mobile-backdrop {
  display: none;
}

@media (max-width: 768px) {
  .sidebar-mobile-backdrop {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(15, 23, 42, 0.45);
    backdrop-filter: blur(4px);
    -webkit-backdrop-filter: blur(4px);
    z-index: 35;
  }
  .crud-pagination-bar {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
  }
  .pagination-info {
    justify-content: space-between;
  }
  .pagination-controls {
    justify-content: center;
    flex-wrap: wrap;
  }
}
</style>
