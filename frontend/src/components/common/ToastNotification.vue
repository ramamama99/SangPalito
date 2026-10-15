<template>
  <div v-if="queueStore.notice" class="toast" role="status">
    <CheckCircle2 :size="20" />
    <span>{{ queueStore.notice }}</span>
    <button @click="queueStore.clearNotice()" aria-label="Tutup">
      <X :size="16" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { CheckCircle2, X } from "lucide-vue-next";
import { useQueueStore } from "../../stores/queue";

const queueStore = useQueueStore();
</script>
