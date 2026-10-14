<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter, useRoute, RouterView } from 'vue-router'
import {
  Activity,
  Users,
  ArrowUpRight,
  Menu,
  X,
  CheckCircle2,
  ShieldCheck,
  Check,
  Search,
  ArrowRight,
} from 'lucide-vue-next'
import { useQueueStore } from './stores/queue'

const router = useRouter()
const route = useRoute()
const queueStore = useQueueStore()

const menuOpen = ref(false)
const lookupNik = ref('')
const lookupNumber = ref('')

const isSpecialPage = computed(() => {
  const name = route.name as string
  return ['tv', 'staff', 'medical-record', 'scan'].includes(name)
})

const navigate = (path: string) => {
  menuOpen.value = false
  router.push(path)
}

const scrollToSection = (sectionId: string) => {
  menuOpen.value = false
  if (route.path !== '/') {
    router.push('/').then(() => {
      setTimeout(() => {
        document.getElementById(sectionId)?.scrollIntoView({ behavior: 'smooth' })
      }, 50)
    })
  } else {
    document.getElementById(sectionId)?.scrollIntoView({ behavior: 'smooth' })
  }
}

const handleLookupSubmit = () => {
  const found = queueStore.patients.find(
    (p) =>
      p.nik === lookupNik.value.trim() &&
      p.number.toUpperCase() === lookupNumber.value.trim().toUpperCase(),
  )
  if (found) {
    queueStore.ticket = found.id
    queueStore.lookupOpen = false
    lookupNik.value = ''
    lookupNumber.value = ''
    router.push('/ticket')
  } else {
    queueStore.setNotice('Tiket tidak ditemukan. Periksa kembali NIK dan nomor antrean.')
  }
}

const confirmVerify = () => {
  if (!queueStore.verifyPatient) return
  queueStore.changeStatus(queueStore.verifyPatient.id, 'Hadir')
  queueStore.verifyPatient = null
  queueStore.setNotice('Identitas terverifikasi. Pasien siap dilayani.')
}

const rejectVerify = () => {
  if (!queueStore.verifyPatient) return
  queueStore.changeStatus(queueStore.verifyPatient.id, 'Ditolak')
  queueStore.verifyPatient = null
  queueStore.setNotice('Data pasien ditolak. Silakan lakukan pendaftaran ulang.')
}
</script>

<template>
  <div>
    <!-- Header -->
    <header v-if="!isSpecialPage" class="header">
      <div class="nav-wrap">
        <button class="brand" @click="navigate('/')" aria-label="AntriSehat beranda">
          <span class="brand-mark">
            <Activity :size="27" :stroke-width="2.6" />
          </span>
          <span>
            Antri<span class="teal">Sehat</span>
            <i />
          </span>
        </button>

        <nav :class="['nav', menuOpen ? 'open' : '']">
          <button
            :class="route.path === '/' ? 'nav-active' : ''"
            @click="navigate('/')"
          >
            Beranda
          </button>
          <button @click="scrollToSection('layanan')">Layanan</button>
          <button @click="scrollToSection('cara')">Cara Daftar</button>
          <button @click="scrollToSection('faq')">FAQ</button>
          <button
            :class="route.path === '/live' ? 'nav-active' : ''"
            @click="navigate('/live')"
          >
            <span class="live-dot" /> Antrean Live
          </button>
        </nav>

        <div class="nav-actions">
          <button class="staff-link" @click="router.push('/login')">
            <Users :size="16" /> Login Staf
          </button>

          <button class="button primary small" @click="router.push('/register')">
            Daftar Online <ArrowUpRight :size="17" />
          </button>

          <button class="menu-button" @click="menuOpen = !menuOpen" aria-label="Buka menu">
            <X v-if="menuOpen" />
            <Menu v-else />
          </button>
        </div>
      </div>
    </header>

    <!-- Toast Notification -->
    <div v-if="queueStore.notice" class="toast" role="status">
      <CheckCircle2 :size="20" />
      {{ queueStore.notice }}
      <button @click="queueStore.notice = ''" aria-label="Tutup">
        <X :size="16" />
      </button>
    </div>

    <!-- Active Route View -->
    <RouterView />

    <!-- Verification Modal -->
    <div
      v-if="queueStore.verifyPatient"
      class="modal-overlay"
      @click="queueStore.verifyPatient = null"
    >
      <div class="modal" role="dialog" aria-modal="true" aria-label="Verifikasi pasien" @click.stop>
        <button class="modal-close" @click="queueStore.verifyPatient = null" aria-label="Tutup">
          <X />
        </button>
        <span class="service-icon">
          <ShieldCheck :size="30" />
        </span>
        <h2>Verifikasi identitas pasien</h2>
        <p>Cocokkan data dengan KTP dan kartu BPJS fisik.</p>

        <div class="verify-data">
          <div>
            <span>Nomor antrean</span>
            <strong>{{ queueStore.verifyPatient.number }}</strong>
          </div>
          <div>
            <span>Nama lengkap</span>
            <strong>{{ queueStore.verifyPatient.name }}</strong>
          </div>
          <div>
            <span>NIK</span>
            <strong>{{ queueStore.verifyPatient.nik }}</strong>
          </div>
          <div>
            <span>Penjamin</span>
            <strong>{{ queueStore.verifyPatient.insurance }}</strong>
          </div>
          <div>
            <span>Nomor BPJS</span>
            <strong>{{ queueStore.verifyPatient.bpjs || '—' }}</strong>
          </div>
        </div>

        <button class="button primary full-width" @click="confirmVerify">
          <Check :size="18" /> Konfirmasi Hadir
        </button>
        <button class="button danger full-width" @click="rejectVerify">
          Tolak · Data Tidak Valid
        </button>
      </div>
    </div>

    <!-- Lookup Modal -->
    <div v-if="queueStore.lookupOpen" class="modal-overlay" @click="queueStore.lookupOpen = false">
      <div class="modal" role="dialog" aria-modal="true" aria-label="Cek antrean" @click.stop>
        <button class="modal-close" @click="queueStore.lookupOpen = false" aria-label="Tutup">
          <X />
        </button>
        <span class="service-icon">
          <Search :size="26" />
        </span>
        <h2>Cek antrean Anda</h2>
        <p>Gunakan NIK dan nomor antrean untuk menemukan tiket.</p>

        <form @submit.prevent="handleLookupSubmit">
          <label>
            NIK
            <input
              required
              pattern="[0-9]{16}"
              v-model="lookupNik"
              placeholder="16 digit NIK Anda"
            />
          </label>
          <label>
            Nomor antrean
            <input required placeholder="Contoh: A-012" v-model="lookupNumber" />
          </label>
          <button class="button primary full-width" type="submit">
            Temukan Tiket <ArrowRight :size="18" />
          </button>
        </form>
      </div>
    </div>
  </div>
</template>
