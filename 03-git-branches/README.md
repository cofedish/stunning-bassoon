# 03. Git: ветки, merge, rebase

Финал — это скорость. Все нужные команды для работы с ветками — на одной странице.

## Что такое ветка

Ветка — это **подвижный указатель** на коммит. Создание ветки — это создание ещё одного указателя, копирования файлов нет. Это очень дёшево.

```
       main
        ↓
A ──── B ──── C
              ↑
            HEAD
```

После `git switch -c feature`:

```
       main    feature
        ↓        ↓
A ──── B ──── C
                ↑
              HEAD
```

## Команды для веток

### Список / создание / переключение

```bash
git branch                       # список локальных
git branch -a                    # + удалённые
git branch -v                    # + последний коммит каждой ветки

git switch main                  # переключиться (новая команда, рекомендуется)
git switch -c feature            # создать ветку и переключиться
git checkout feature             # старая команда, делает то же
git checkout -b feature          # старый эквивалент switch -c
```

### Удаление

```bash
git branch -d feature            # безопасно (только если влита)
git branch -D feature            # принудительно
git push origin --delete feature # удалить на remote
```

### Переименование

```bash
git branch -m old new            # переименовать (если ты на другой ветке)
git branch -m new                # переименовать текущую ветку
```

## Слияние (merge)

```bash
# из main: подтянуть feature
git switch main
git merge feature
```

Что произойдёт:
- **Fast-forward** — если main не двигалась с момента создания feature, указатель просто переедет вперёд.
- **Three-way merge** — если main ушла вперёд параллельно, создастся merge-коммит с двумя родителями.

```bash
git merge --no-ff feature        # всегда merge-коммит, даже если можно ff
git merge --squash feature       # все коммиты feature → один новый коммит на main
```

## Rebase (перебазирование)

Берёт твои коммиты с одной ветки и **переписывает** их поверх другой.

```bash
git switch feature
git rebase main
```

```
До:
       feature
A ── B ── C ── D ── E
        \         ↑
         X ── Y
              ↑
            main

После git rebase main:
                       feature
A ── B ── C ── D ── E ── X' ── Y'
                      ↑
                    main
```

Плюсы: чистая линейная история. Минусы: меняет хеши, **никогда не делай rebase коммитов, которые уже запушены и используются другими**.

```bash
git rebase -i HEAD~5             # интерактивный: squash, reword, drop
```

## Конфликты слияния

Когда git не может автоматически объединить — он отмечает места конфликта в файле:

```
<<<<<<< HEAD
твоя версия
=======
их версия
>>>>>>> feature
```

Алгоритм решения:
1. Открыть файл, выбрать правильную версию (или объединить руками).
2. Удалить маркеры `<<<<<<<`, `=======`, `>>>>>>>`.
3. `git add <file>`.
4. `git commit` (для merge) или `git rebase --continue` (для rebase).

```bash
git merge --abort                # передумал, откатить merge
git rebase --abort               # передумал, откатить rebase
```

## Stash — временно отложить изменения

```bash
git stash                        # отложить незакоммиченное
git stash push -m "wip auth"     # с подписью
git stash list                   # список заначек
git stash pop                    # вернуть последнюю и удалить из стека
git stash apply stash@{1}        # применить конкретную, не удаляя
git stash drop stash@{1}         # удалить заначку
```

## Cherry-pick — взять отдельный коммит из другой ветки

```bash
git cherry-pick <commit-hash>
git cherry-pick A^..B            # диапазон
```

## Типичные сценарии

### Я закоммитил не в ту ветку

```bash
git switch correct-branch
git cherry-pick <hash>
git switch wrong-branch
git reset --hard HEAD~1          # снять коммит с неправильной ветки
```

### Хочу обновить feature свежим main

```bash
git switch main
git pull
git switch feature
git rebase main                  # или git merge main
```

### Запушил кривое имя коммита

```bash
git commit --amend -m "новое сообщение"
git push --force-with-lease      # ! только если ветка только твоя
```

## Полезные ссылки

- [Работа с Git (официальный туториал)](https://git-scm.com/book/ru/v2/Ветвление-в-Git-Основы-ветвления-и-слияния)
