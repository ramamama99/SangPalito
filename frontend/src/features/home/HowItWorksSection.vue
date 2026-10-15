<template>
  <section class="how-section container" id="cara">
    <div class="section-label"><span /> SEMUDAH INI</div>
    <h2>
      Antrean lebih singkat.
      <br />
      Hanya tiga langkah.
    </h2>
    <div class="how-grid">
      <div v-for="step in steps" :key="step.number">
        <span class="step-number">{{ step.number }}</span>
        <h3>{{ step.title }}</h3>
        <p>{{ step.desc }}</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const steps = [
  {
    number: "01",
    title: "Daftar dari mana saja",
    desc: "Isi identitas dan pilih poli. Tidak perlu membuat akun.",
  },
  {
    number: "02",
    title: "Simpan tiket Anda",
    desc: "Unduh e-tiket dengan QR untuk check-in di lokasi.",
  },
  {
    number: "03",
    title: "Datang, scan, dan dilayani",
    desc: "Bawa identitas, scan QR, lalu tunggu panggilan Anda.",
  },
];
</script>
