<template>
  <div v-if="queueStore.lookupOpen" class="modal-overlay" @click="closeModal">
    <div
      class="modal"
      role="dialog"
      aria-modal="true"
      aria-label="Cek antrean"
      @click.stop
    >
      <button class="modal-close" @click="closeModal" aria-label="Tutup">
        <X />
      </button>
      <span class="service-icon">
        <Search :size="26" />
      </span>
      <h2>Cek antrean Anda</h2>
      <p>Gunakan NIK dan nomor antrean untuk menemukan tiket.</p>
      <form @submit.prevent="handleSearch">
        <label>
          NIK
          <input
            required
            pattern="[0-9]{16}"
            v-model="nikInput"
            placeholder="16 digit NIK Anda"
          />
        </label>
        <label>
          Nomor antrean
          <input required placeholder="Contoh: A-012" v-model="numberInput" />
        </label>
        <button class="button primary full-width" type="submit">
          Temukan Tiket <ArrowRight :size="18" />
        </button>
      </form>
    </div>
  </div>
</template>

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
