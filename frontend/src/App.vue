<template>
  <div>
    <header class="header">
      <div class="header__inner">
        <span class="header__title">Сервис анонимизации голосовых данных</span>
      </div>
    </header>
    <RouterView />
  </div>
</template>

<style>
* { box-sizing: border-box; margin: 0; padding: 0; }
body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  background: #0f172a;
  color: #e2e8f0;
  -webkit-font-smoothing: antialiased;
}
</style>

<style scoped>
.header {
  background: #0b1120;
  border-bottom: 1px solid #1e293b;
  height: 52px;
  display: flex;
  align-items: center;
  padding: 0 24px;
}

.header__inner {
  max-width: 960px;
  width: 100%;
  margin: 0 auto;
}

.header__title {
  font-size: 15px;
  font-weight: 600;
  color: #e2e8f0;
  letter-spacing: 0.02em;
}
</style>
