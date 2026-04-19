# QA-пакет проекта voice-redaction

Документация подготовлена по стандарту QA Service Lab ([rules_testing.md](../../rules_testing.md)).

## Содержимое
- [Test Plan](test-plan.md) — стратегия, виды тестирования, критерии входа/выхода.
- [Test Cases](test-cases/README.md) — 13 кейсов, покрывают все 5 классификаций.
- [Bug Reports](bug-reports/README.md) — найденные дефекты по шаблону раздела 4.1.
- [Daily Report](reports/daily-2026-04-19.md) — ежедневный отчёт.
- [Bug Summary](reports/bug-summary.md) — агрегированный отчёт по багам.
- [Test Summary](reports/test-summary.md) — итоговый отчёт с метриками и рекомендацией GO/NO-GO.

## Соответствие стандарту

| Раздел стандарта | Артефакт |
|------------------|----------|
| 2. Виды тестирования | [test-plan.md §3.1](test-plan.md) |
| 4. Шаблон баг-репорта | [bug-reports/BUG-*.md](bug-reports/) |
| 5. Жизненный цикл бага | [bug-reports/README.md](bug-reports/README.md) (статусы) |
| 6.1 Daily Report | [reports/daily-2026-04-19.md](reports/daily-2026-04-19.md) |
| 6.2 Test Summary | [reports/test-summary.md](reports/test-summary.md) |
| 6.3 Bug Report | [reports/bug-summary.md](reports/bug-summary.md) |
| 7. Критерии | [test-plan.md §6](test-plan.md) + [test-summary.md §6](reports/test-summary.md) |
| 8. Метрики | [test-summary.md §5](reports/test-summary.md) |
| 10.1 Test Plan | [test-plan.md](test-plan.md) |
| 10.2 Test Case | [test-cases/TC-*.md](test-cases/) |

## Краткая сводка
- 13 тест-кейсов / 12 Passed / 1 Failed / 0 Blocked.
- 2 найденных бага: 1 High (BUG-001), 1 Low (BUG-002).
- Покрытие требований: 100%.
- Рекомендация: **Conditional GO** для MVP, **NO-GO** для production до закрытия BUG-001.
