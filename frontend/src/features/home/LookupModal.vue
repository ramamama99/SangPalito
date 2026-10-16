<template>
  <div v-if="queueStore.lookupOpen" class="modal-overlay" @click="closeModal">
    <div
      class="modal lookup-modal"
      role="dialog"
      aria-modal="true"
      aria-label="Cek antrean"
      @click.stop
    >
      <div class="modal-top-bar">
        <div class="modal-badge-icon">
          <Search :size="22" />
        </div>
        <button class="lookup-close-btn" @click="closeModal" aria-label="Tutup" title="Tutup Modal">
          <X :size="18" />
        </button>
      </div>

      <h2 class="lookup-heading">Cek antrean Anda</h2>
      <p class="lookup-sub">Gunakan NIK dan nomor antrean untuk menemukan tiket aktif Anda.</p>

      <form @submit.prevent="handleSearch" class="lookup-form">
        <label class="lookup-field">
          <span class="field-label">Nomor Induk Kependudukan (NIK)</span>
          <input
            required
            pattern="[0-9]{16}"
            v-model="nikInput"
            placeholder="Masukkan 16 digit NIK Anda"
            class="lookup-input"
          />
        </label>

        <label class="lookup-field">
          <span class="field-label">Nomor Antrean</span>
          <input
            required
            placeholder="Contoh: A-012"
            v-model="numberInput"
            class="lookup-input"
          />
        </label>

        <button class="lookup-submit-btn" type="submit">
          <span>Temukan Tiket</span>
          <ArrowRight :size="18" />
        </button>
      </form>
    </div>
  </div>
</template>

<style scoped>
.lookup-modal {
  max-width: 460px;
  width: 100%;
  padding: 28px 30px;
  border-radius: 22px;
  border: 1px solid var(--line);
  box-shadow: 0 25px 60px -15px rgba(25, 90, 95, 0.2);
  background: #ffffff;
}

.modal-top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.modal-badge-icon {
  width: 46px;
  height: 46px;
  border-radius: 14px;
  background: #eaf6f2;
  border: 1px solid #c8e8df;
  color: var(--deep);
  display: flex;
  align-items: center;
  justify-content: center;
}

.lookup-close-btn {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: #f6f8f5;
  border: 1px solid var(--line);
  color: var(--muted);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.lookup-close-btn:hover {
  background: #fdf3f0;
  border-color: #f7d6cf;
  color: var(--coral);
}

.lookup-heading {
  font-family: 'Manrope', sans-serif !important;
  font-size: 24px !important;
  font-weight: 700 !important;
  letter-spacing: -0.3px !important;
  color: var(--deep) !important;
  margin: 0 0 6px 0 !important;
  line-height: 1.25 !important;
}

.lookup-sub {
  font-size: 13.5px;
  color: var(--muted);
  margin: 0 0 22px 0;
  line-height: 1.5;
}

.lookup-form {
  display: flex;
  flex-direction: column;
}

.lookup-field {
  display: flex;
  flex-direction: column;
  gap: 7px;
  margin-bottom: 16px;
  text-align: left;
}

.field-label {
  font-size: 13px;
  font-weight: 700;
  color: var(--deep);
}

.lookup-input {
  width: 100%;
  height: 46px;
  border: 1px solid #dce5df;
  border-radius: 9px;
  padding: 10px 14px;
  background: #ffffff;
  font-size: 14px;
  color: #243e40;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
  box-sizing: border-box;
}

.lookup-input:focus {
  border-color: var(--teal);
  box-shadow: 0 0 0 2px rgba(18, 184, 176, 0.2);
}

.lookup-input::placeholder {
  color: #a0aba6;
}

.lookup-submit-btn {
  width: 100%;
  height: 48px;
  border-radius: 50px;
  background: var(--deep);
  color: #ffffff;
  font-size: 15px;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-top: 6px;
  box-shadow: 0 4px 14px rgba(25, 90, 95, 0.2);
}

.lookup-submit-btn:hover {
  background: #237079;
  transform: translateY(-1px);
}
</style>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { Search, ArrowRight, X } from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";

const router = useRouter();
const queueStore = useQueueStore();

const nikInput = ref("");
const numberInput = ref("");

function closeModal() {
  queueStore.lookupOpen = false;
}

function handleSearch() {
  const p = queueStore.patients.find(
    (item) =>
      item.nik === nikInput.value.trim() &&
      item.number.toUpperCase() === numberInput.value.trim().toUpperCase(),
  );

  if (p) {
    queueStore.ticket = p.id;
    queueStore.lookupOpen = false;
    router.push("/ticket");
    window.scrollTo(0, 0);
  } else {
    queueStore.setNotice(
      "Tiket tidak ditemukan. Periksa kembali NIK dan nomor antrean.",
    );
  }
}
</script>
