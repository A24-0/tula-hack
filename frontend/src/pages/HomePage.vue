<template>
  <main class="page-shell">
    <header class="site-header container">
      <div class="brand" aria-label="Voice Redaction">
        <span class="brand__mark"></span>
        <strong>Voice Redaction</strong>
      </div>
    </header>
    <HeroIntro />
    <MainActionPanel />
    <ProcessingPipeline />
    <ResultSection />
  </main>
</template>

<script setup lang="ts">
import MainActionPanel from '@/components/actions/MainActionPanel.vue'
import HeroIntro from '@/components/hero/HeroIntro.vue'
import ProcessingPipeline from '@/components/processing/ProcessingPipeline.vue'
import ResultSection from '@/components/result/ResultSection.vue'
</script>

<style scoped>
.site-header {
  position: sticky;
  top: 0;
  z-index: 20;
  display: flex;
  align-items: center;
  min-height: 64px;
  padding: 12px 0;
  animation: heroReveal 720ms var(--ease-out) both;
}

.site-header::before {
  position: fixed;
  top: 0;
  left: 0;
  z-index: -1;
  width: 100%;
  height: 72px;
  pointer-events: none;
  content: "";
  background: linear-gradient(180deg, rgba(5, 6, 7, 0.86), rgba(5, 6, 7, 0));
  backdrop-filter: blur(14px);
}

.brand {
  display: inline-flex;
  align-items: center;
  width: max-content;
  gap: 10px;
  color: var(--text-primary);
}

.brand__mark {
  position: relative;
  width: 22px;
  height: 22px;
  overflow: hidden;
  background: linear-gradient(135deg, #f4f4f2 0 18%, transparent 18% 30%, #f4f4f2 30% 48%, transparent 48% 60%, #f4f4f2 60% 100%);
  border-radius: 999px;
  box-shadow: 0 0 24px rgba(255, 255, 255, 0.18);
  transform: rotate(36deg);
}

.brand strong {
  font-size: 15px;
  font-weight: 760;
  letter-spacing: 0;
}
</style>
