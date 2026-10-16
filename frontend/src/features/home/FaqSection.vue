<template>
  <section class="faq-section container" id="faq">
    <div>
      <div class="section-label"><span /> ADA PERTANYAAN?</div>
      <h2>Kami siap membantu.</h2>
      <p>Hal yang perlu diketahui sebelum berkunjung.</p>
    </div>
    <div>
      <div v-for="(item, i) in faqItems" :key="item.question" class="faq-item">
        <button @click="toggleFaq(i)">
          {{ item.question }}
          <ChevronDown :size="19" />
        </button>
        <p v-if="activeFaq === i">{{ item.answer }}</p>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ChevronDown } from "lucide-vue-next";
import { faqItems } from "../../constants/faq";

const activeFaq = ref<number | null>(0);

function toggleFaq(index: number) {
  activeFaq.value = activeFaq.value === index ? null : index;
}
</script>
