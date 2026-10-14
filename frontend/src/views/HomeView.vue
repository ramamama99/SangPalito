<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  Activity,
  ArrowUpRight,
  ShieldCheck,
  Clock3,
  CalendarDays,
  CheckCircle2,
  Heart,
  Search,
  MapPin,
  Phone,
  ChevronDown,
} from 'lucide-vue-next'
import { services } from '../constants/services'
import { useQueueStore } from '../stores/queue'

const router = useRouter()
const queueStore = useQueueStore()
const activeFaq = ref<number | null>(0)

const startRegister = (poli = 'Poli Umum') => {
  router.push({ path: '/register', query: { poli } })
}

const toggleFaq = (index: number) => {
  activeFaq.value = activeFaq.value === index ? null : index
}

const faqs = [
  [
    'Apakah saya perlu membuat akun?',
    'Tidak. Anda cukup mengisi identitas sesuai KTP, nomor WhatsApp, dan memilih jadwal. Simpan e-tiket untuk check-in.',
  ],
  [
    'Apa yang perlu dibawa saat berkunjung?',
    'Bawa KTP, kartu BPJS atau kartu asuransi (jika digunakan), serta e-tiket yang telah diunduh.',
  ],
  [
    'Bagaimana jika saya terlambat atau ingin membatalkan?',
    'Anda dapat membatalkan dari halaman e-tiket. Jika terlambat, hubungi petugas loket untuk penyesuaian antrean.',
  ],
  [
    'Apakah pendaftaran ini berbayar?',
    'Pendaftaran antrean gratis. Biaya pelayanan mengikuti penjamin dan ketentuan puskesmas.',
  ],
]
</script>

<template>
  <div>
    <main class="hero container">
      <div class="hero-copy">
        <h1>
          Lebih sedikit
          <br />
          menunggu.
          <br />
          <span>
            Lebih banyak
            <br class="mobile-break" />
            hidup sehat.
          </span>
        </h1>
        <p>
          Daftar dari rumah, datang saat giliran Anda.
          <br class="desktop-break" />
          Pelayanan puskesmas kini lebih mudah, cepat, <br class="desktop-break" />
          dan nyaman untuk semua.
        </p>
        <div class="hero-actions">
          <button class="button primary" @click="startRegister()">
            Ambil Antrean
            <span class="arrow-disc">
              <ArrowUpRight :size="20" />
            </span>
          </button>
          <button class="button outline" @click="queueStore.lookupOpen = true">
            <Search :size="17" /> Cek Antrean Saya
          </button>
        </div>
        <div class="social-proof">
          <div class="proof-stars">★★★★★</div>
        </div>
      </div>

      <div class="hero-visual">
        <div class="photo-frame">
          <img src="/images/doctor.jpg" alt="Dokter siap melayani kesehatan Anda" />
          <div class="photo-tint" />
          <span class="decor-plus">+</span>
        </div>
        <div class="floating-benefits">
          <div>
            <span>
              <ShieldCheck :size="20" />
            </span>
            <p>Data aman<small>& terlindungi</small></p>
          </div>
          <div>
            <span>
              <Clock3 :size="20" />
            </span>
            <p>Tanpa antre<small>lebih lama</small></p>
          </div>
        </div>
        <div class="floating-queue">
          <div class="queue-icon">
            <CalendarDays :size="24" />
          </div>
          <div>
            <small>Antrean Anda, dalam genggaman</small>
            <strong>Daftar kurang dari 2 menit</strong>
            <span> <span class="live-dot" /> Mudah. Cepat. Tanpa login. </span>
          </div>
          <CheckCircle2 :size="22" />
        </div>
      </div>
    </main>

    <div class="trust-strip container">
      <span> <ShieldCheck :size="18" /> Data pasien terlindungi </span>
      <i />
      <span> <Clock3 :size="18" /> Waktu tunggu 30% lebih singkat* </span>
      <i />
      <span> <CheckCircle2 :size="18" /> Terintegrasi dengan layanan BPJS </span>
      <i />
      <span> <Heart :size="18" /> Ramah untuk semua usia </span>
    </div>

    <section class="services-section" id="layanan">
      <div class="container">
        <div class="section-heading">
          <div>
            <div class="section-label"><span /> LAYANAN KESEHATAN</div>
            <h2>
              Perawatan yang tepat,
              <br />
              untuk setiap kebutuhan.
            </h2>
          </div>
          <p>
            Pilih layanan yang Anda butuhkan.
            <br />
            Kami siap mendampingi perjalanan sehat Anda.
          </p>
        </div>
        <div class="service-grid">
          <button
            v-for="(s, i) in services"
            :key="s.name"
            class="service-card"
            @click="startRegister(s.name)"
          >
            <span :class="['service-icon', `color-${i}`]">
              <component :is="s.icon" :size="26" />
            </span>
            <span class="availability"> <span class="live-dot" /> Tersedia </span>
            <h3>{{ s.name }}</h3>
            <p>{{ s.desc }}</p>
            <div class="card-bottom">
              <span> <Clock3 :size="14" /> 08.00 – 14.00 </span>
              <span class="card-arrow">
                <ArrowUpRight :size="21" />
              </span>
            </div>
          </button>
        </div>
        <p class="section-footnote">
          <MapPin :size="14" /> Puskesmas Sehat Bersama · Senin – Sabtu
          <span>*Target efisiensi pelayanan AntriSehat</span>
        </p>
      </div>
    </section>

    <section class="how-section container" id="cara">
      <div class="section-label"><span /> SEMUDAH INI</div>
      <h2>
        Antrean lebih singkat.
        <br />
        Hanya tiga langkah.
      </h2>
      <div class="how-grid">
        <div
          v-for="[n, title, desc] in [
            [
              '01',
              'Daftar dari mana saja',
              'Isi identitas dan pilih poli. Tidak perlu membuat akun.',
            ],
            ['02', 'Simpan tiket Anda', 'Unduh e-tiket dengan QR untuk check-in di lokasi.'],
            [
              '03',
              'Datang, scan, dan dilayani',
              'Bawa identitas, scan QR, lalu tunggu panggilan Anda.',
            ],
          ]"
          :key="n"
        >
          <span class="step-number">{{ n }}</span>
          <h3>{{ title }}</h3>
          <p>{{ desc }}</p>
        </div>
      </div>
    </section>

    <section class="faq-section container" id="faq">
      <div>
        <div class="section-label"><span /> ADA PERTANYAAN?</div>
        <h2>Kami siap membantu.</h2>
        <p>Hal yang perlu diketahui sebelum berkunjung.</p>
        <a class="contact-link" href="https://wa.me/6281234567890" target="_blank" rel="noreferrer">
          <Phone :size="17" /> Hubungi pusat bantuan
          <ArrowUpRight :size="17" />
        </a>
      </div>
      <div>
        <div v-for="([q, a], i) in faqs" :key="q" class="faq-item">
          <button @click="toggleFaq(i)">
            {{ q }}
            <ChevronDown :size="19" />
          </button>
          <p v-if="activeFaq === i">{{ a }}</p>
        </div>
      </div>
    </section>

    <section class="cta container">
      <div>
        <div class="section-label">KESEHATAN TIDAK PERLU MENUNGGU</div>
        <h2>Mulai langkah sehat Anda hari ini.</h2>
        <p>Ambil antrean dari rumah. Kami menantikan Anda.</p>
      </div>
      <button class="button white" @click="startRegister()">
        Daftar Sekarang <ArrowUpRight :size="20" />
      </button>
    </section>

    <footer class="container">
      <button class="brand" @click="router.push('/')" aria-label="AntriSehat beranda">
        <span class="brand-mark">
          <Activity :size="27" :stroke-width="2.6" />
        </span>
        <span>
          Antri<span class="teal">Sehat</span>
          <i />
        </span>
      </button>
    </footer>
  </div>
</template>
