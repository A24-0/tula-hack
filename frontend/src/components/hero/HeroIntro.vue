<template>
  <section id="product" class="hero container">
    <div class="hero__copy">
      <h1>Анонимизация голосовых записей</h1>
      <p>
        Загрузите звонок или запишите речь в браузере. Voice Redaction распознаёт аудио,
        находит персональные данные и создаёт безопасную версию текста, записи и отчёта.
      </p>
    </div>

    <div class="product-stage" aria-label="Демонстрация интерфейса Voice Redaction">
      <div class="product-window">
        <aside class="product-sidebar">
          <div class="sidebar-brand">
            <span class="sidebar-brand__mark"></span>
            <strong>Redaction</strong>
          </div>

        </aside>

        <section class="workbench">
          <div class="workbench__bar">
            <strong>call-042.wav</strong>
            <span>03 / 14</span>
          </div>

          <article class="audio-card">
            <div class="audio-card__head">
              <div>
                <span class="eyebrow">Secure audio review</span>
                <h2>Проверка звонка перед передачей в CRM</h2>
              </div>
              <span class="status-pill">обработка</span>
            </div>

            <div class="wave-strip">
              <i
                v-for="(bar, index) in bars"
                :key="`${bar}-${index}`"
                :class="{ muted: maskedBars.includes(index) }"
                :style="{ height: `${bar}%` }"
              ></i>
            </div>

            <div class="transcript-preview">
              <p>
                Добрый день, меня зовут <mark>Иван Петров</mark>, номер договора
                <mark>77-1842</mark>. Телефон для связи <mark>+7 900 000-00-00</mark>.
              </p>
            </div>
          </article>

          <div class="activity">
            <h3>Activity</h3>
            <div v-for="event in events" :key="event.title" class="activity__row">
              <span :class="event.tone"></span>
              <div>
                <strong>{{ event.title }}</strong>
                <p>{{ event.copy }}</p>
              </div>
            </div>
          </div>
        </section>

        <aside class="inspector">
          <div class="inspector__top">
            <span>VR-2703</span>
            <span class="dot"></span>
          </div>
          <dl>
            <div>
              <dt>Status</dt>
              <dd>In progress</dd>
            </div>
            <div>
              <dt>Priority</dt>
              <dd>High</dd>
            </div>
            <div>
              <dt>Owner</dt>
              <dd>Security team</dd>
            </div>
          </dl>

          <div class="labels">
            <span>PERSON</span>
            <span>PHONE</span>
            <span>PASSPORT</span>
          </div>
        </aside>
      </div>

      <div class="assistant-card">
        <div class="assistant-card__head">
          <strong>AI redactor</strong>
          <button type="button" aria-label="Закрыть">×</button>
        </div>
        <p>Запустил транскрибацию, нашёл 7 фрагментов ПД и подготовил маскирование аудио.</p>
        <div class="assistant-card__steps">
          <span>speech-to-text</span>
          <span>pii-detect</span>
          <span>audio-mask</span>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
const bars = [24, 46, 66, 38, 74, 52, 30, 84, 68, 42, 78, 56, 34, 64, 88, 48, 32, 72, 92, 58, 28, 62, 44, 80, 52, 36, 70, 86, 40, 60]
const maskedBars = [4, 5, 10, 14, 18, 19, 26]
const events = [
  { title: 'Речь распознана', copy: 'Транскрипт синхронизирован с таймкодами', tone: 'green' },
  { title: 'PII обнаружены', copy: 'ФИО, телефон и номер договора выделены', tone: 'blue' },
  { title: 'Аудио замаскировано', copy: 'Чувствительные интервалы приглушены', tone: 'red' },
]
</script>

<style scoped>
.hero {
  padding-top: clamp(82px, 13vw, 178px);
  padding-bottom: 52px;
}

.hero__copy {
  display: grid;
  justify-items: start;
  max-width: 920px;
  animation: heroReveal 760ms var(--ease-out) both;
}

.hero h1 {
  max-width: 930px;
  margin: 0;
  font-size: clamp(48px, 7.2vw, 82px);
  font-weight: 650;
  line-height: 0.98;
  letter-spacing: -0.067em;
}

.hero p {
  max-width: 720px;
  margin: 26px 0 0;
  color: var(--text-muted);
  font-size: clamp(15px, 1.8vw, 18px);
  font-weight: 600;
  line-height: 1.65;
}

.product-stage {
  position: relative;
  margin-top: clamp(58px, 8vw, 82px);
  animation: heroReveal 860ms 120ms var(--ease-out) both;
}

.product-stage::before {
  position: absolute;
  inset: auto -6% -18% -6%;
  height: 42%;
  pointer-events: none;
  content: "";
  background: radial-gradient(ellipse at 50% 100%, rgba(145, 153, 255, 0.26), transparent 68%);
  filter: blur(18px);
}

.product-window {
  position: relative;
  display: grid;
  grid-template-columns: 240px minmax(0, 1fr) 270px;
  min-height: 520px;
  overflow: hidden;
  background: linear-gradient(180deg, rgba(17, 18, 22, 0.98), rgba(9, 10, 12, 0.94));
  border: 1px solid var(--border-strong);
  border-radius: 14px;
  box-shadow: 0 40px 130px rgba(0, 0, 0, 0.62), 0 0 0 1px rgba(255, 255, 255, 0.02) inset;
}

.product-sidebar,
.inspector {
  padding: 20px;
  background: rgba(255, 255, 255, 0.018);
}

.product-sidebar {
  border-right: 1px solid var(--border-soft);
}

.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-bottom: 28px;
}

.sidebar-brand__mark {
  width: 16px;
  height: 16px;
  background: linear-gradient(135deg, #f4f4f2 0 20%, transparent 20% 36%, #f4f4f2 36% 58%, transparent 58% 72%, #f4f4f2 72% 100%);
  border-radius: 999px;
  transform: rotate(36deg);
}

.sidebar-brand strong {
  font-size: 14px;
}

.workbench {
  min-width: 0;
}

.workbench__bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 48px;
  padding: 0 28px;
  color: var(--text-muted);
  border-bottom: 1px solid var(--border-soft);
}

.workbench__bar strong {
  color: var(--text-secondary);
}

.audio-card {
  padding: clamp(26px, 4vw, 54px);
}

.audio-card__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
}

.audio-card h2 {
  max-width: 620px;
  margin: 14px 0 0;
  font-size: clamp(26px, 3.2vw, 42px);
  font-weight: 650;
  line-height: 1.08;
  letter-spacing: -0.045em;
}

.status-pill {
  display: inline-flex;
  min-height: 28px;
  align-items: center;
  padding: 0 10px;
  color: var(--accent-amber);
  font-family: var(--font-mono);
  font-size: 11px;
  background: rgba(228, 189, 115, 0.1);
  border: 1px solid rgba(228, 189, 115, 0.2);
  border-radius: 999px;
}

.wave-strip {
  display: flex;
  align-items: center;
  height: 150px;
  gap: 7px;
  margin-top: 34px;
  padding: 18px 0;
}

.wave-strip i {
  flex: 1;
  min-width: 3px;
  background: linear-gradient(180deg, rgba(246, 246, 243, 0.5), rgba(246, 246, 243, 0.16));
  border-radius: 999px;
}

.wave-strip i.muted {
  background: linear-gradient(180deg, rgba(255, 139, 139, 0.72), rgba(255, 139, 139, 0.22));
  box-shadow: 0 0 24px rgba(255, 139, 139, 0.16);
}

.transcript-preview {
  max-width: 760px;
  padding: 20px 22px;
  color: var(--text-secondary);
  background: rgba(255, 255, 255, 0.035);
  border: 1px solid var(--border-soft);
  border-radius: 12px;
}

.transcript-preview p {
  margin: 0;
  font-size: 16px;
  line-height: 1.75;
}

mark {
  padding: 2px 7px;
  color: #f8eeee;
  background: rgba(255, 139, 139, 0.24);
  border: 1px solid rgba(255, 139, 139, 0.28);
  border-radius: 6px;
}

.activity {
  display: grid;
  gap: 12px;
  padding: 0 clamp(26px, 4vw, 54px) 42px;
}

.activity h3 {
  margin: 0 0 4px;
  color: var(--text-secondary);
  font-size: 15px;
}

.activity__row {
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 12px;
  align-items: start;
  max-width: 780px;
  padding: 13px 0;
  border-top: 1px solid var(--border-soft);
}

.activity__row > span {
  width: 8px;
  height: 8px;
  margin-top: 6px;
  border-radius: 999px;
}

.activity__row .green {
  background: var(--accent-green);
}

.activity__row .blue {
  background: var(--accent-blue);
}

.activity__row .red {
  background: var(--accent-red);
}

.activity__row strong {
  color: var(--text-secondary);
  font-size: 14px;
}

.activity__row p {
  margin: 4px 0 0;
  color: var(--text-muted);
  font-size: 13px;
}

.inspector {
  border-left: 1px solid var(--border-soft);
}

.inspector__top {
  display: flex;
  justify-content: space-between;
  color: var(--text-muted);
  font-family: var(--font-mono);
  font-size: 12px;
}

.dot {
  width: 8px;
  height: 8px;
  background: var(--accent-green);
  border-radius: 999px;
}

dl {
  display: grid;
  gap: 20px;
  margin: 54px 0 0;
}

dt {
  color: var(--text-muted);
  font-size: 12px;
}

dd {
  margin: 8px 0 0;
  color: var(--text-secondary);
  font-size: 14px;
}

.labels {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 42px;
}

.labels span {
  padding: 7px 9px;
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 11px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid var(--border-soft);
  border-radius: 999px;
}

.assistant-card {
  position: absolute;
  right: 22px;
  bottom: -28px;
  z-index: 2;
  width: min(420px, 52vw);
  padding: 16px;
  background: rgba(13, 14, 17, 0.94);
  border: 1px solid var(--border-soft);
  border-radius: 12px;
  box-shadow: var(--shadow-raised);
  backdrop-filter: blur(22px);
  animation: softOrbit 6s ease-in-out infinite;
}

.assistant-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border-soft);
}

.assistant-card__head strong {
  font-size: 14px;
}

.assistant-card__head button {
  width: 28px;
  height: 28px;
  color: var(--text-muted);
  background: transparent;
  border: 0;
  cursor: pointer;
}

.assistant-card p {
  margin: 12px 0 0;
  color: var(--text-muted);
  font-size: 13px;
  line-height: 1.6;
}

.assistant-card__steps {
  display: flex;
  flex-wrap: wrap;
  gap: 7px;
  margin-top: 14px;
}

.assistant-card__steps span {
  padding: 5px 8px;
  color: var(--text-secondary);
  font-family: var(--font-mono);
  font-size: 10px;
  background: rgba(255, 255, 255, 0.045);
  border-radius: 999px;
}

@media (max-width: 1060px) {
  .product-window {
    grid-template-columns: 190px minmax(0, 1fr);
  }

  .inspector {
    display: none;
  }
}

@media (max-width: 760px) {
  .hero {
    padding-top: 52px;
  }

  .product-window {
    grid-template-columns: 1fr;
  }

  .product-sidebar {
    display: none;
  }

  .audio-card__head {
    flex-direction: column;
  }

  .assistant-card {
    position: relative;
    right: auto;
    bottom: auto;
    width: 100%;
    margin-top: 14px;
    animation: none;
  }
}
</style>
