# 09. Команды для работы с удалёнными репозиториями

Удалённый репозиторий (remote) — копия твоей репы где-то в сети: на GitHub, GitVerse, Gitea, Gitlab. Git синхронизирует историю между локалью и remote.

## Понятия

| Термин | Что |
|---|---|
| **remote** | Зарегистрированный на твоей локальной машине адрес удалённой репы |
| **origin** | Стандартное имя главного remote (то, откуда клонировали) |
| **upstream** | Часто — родительская репа форка |
| **tracking branch** | Локальная ветка, привязанная к ветке remote (`main` ↔ `origin/main`) |
| **fork** | Серверная копия чужой репы под твоим аккаунтом |

## Управление remote

```bash
git remote -v                                        # список + URL
git remote add origin git@github.com:user/repo.git   # добавить
git remote set-url origin git@github.com:u/r.git     # сменить URL
git remote rename origin upstream                    # переименовать
git remote remove old                                # удалить
```

## Получить чужие изменения

```bash
git fetch origin                # скачать всё, локальная ветка не меняется
git fetch origin main           # только одну ветку
git fetch --all                 # все remote
git fetch --prune               # ещё и удалить ссылки на удалённые ветки

git pull                        # = fetch + merge (по умолчанию)
git pull --rebase               # = fetch + rebase (часто чище)
git pull origin main            # явно указать
```

**Что выбрать:** `pull --rebase` для своей feature-ветки (чистая история). Обычный `pull` (= merge) для основных веток.

## Отправить свои изменения

```bash
git push origin main            # запушить main в origin
git push                        # если ветка отслеживает remote
git push -u origin feature      # первый раз: создать ветку на remote и привязать
git push --tags                 # запушить теги
git push origin --delete feat   # удалить ветку на remote
git push --force-with-lease     # принудительно, но безопаснее force
git push --force                # ОПАСНО: переписать историю на remote
```

**`--force-with-lease` vs `--force`:** lease проверяет, что remote-ветка с тех пор не менялась кем-то ещё. Если менялась — пуш отклоняется, ничего не теряем.

## Связь локальных и удалённых веток

```bash
git branch -vv                  # локальные ветки + tracking
git branch --set-upstream-to=origin/main main
git push -u origin feature      # коротко: создать tracking при первом пуше
```

## Стандартный workflow (feature branch)

```bash
# 1. Получить актуальный main
git switch main
git pull

# 2. Создать ветку под задачу
git switch -c feat/json-import

# 3. Работать, коммитить
git add .
git commit -m "добавил парсер JSON"

# 4. Запушить
git push -u origin feat/json-import

# 5. Открыть Pull Request в веб-интерфейсе

# 6. Если попросили доработать — коммитим, пушим — PR обновится сам
git commit -am "правки по ревью"
git push

# 7. После merge — почистить локальные ветки
git switch main
git pull
git branch -d feat/json-import
git fetch --prune
```

## Работа с форком

```bash
# 1. Форкнул чужую репу через UI GitHub.
# 2. Клонировал свой форк
git clone git@github.com:me/some-repo.git
cd some-repo

# 3. Добавил оригинал как upstream — чтобы подтягивать обновления
git remote add upstream git@github.com:original-author/some-repo.git
git remote -v
# origin    git@github.com:me/some-repo.git
# upstream  git@github.com:original-author/some-repo.git

# 4. Регулярно синхронизируешь main с upstream
git fetch upstream
git switch main
git merge upstream/main
git push origin main
```

## Теги (для релизов)

```bash
git tag v1.0.0                  # лёгкий тег
git tag -a v1.0.0 -m "релиз"    # аннотированный (с автором/сообщением)
git push origin v1.0.0          # запушить один тег
git push --tags                 # запушить все
git tag -d v1.0.0               # удалить локально
git push origin :refs/tags/v1.0.0   # удалить на remote
```

## Решение типовых проблем

### `! [rejected]   main -> main (non-fast-forward)`

На remote есть коммиты, которых у тебя нет. Подтяни:
```bash
git pull --rebase
git push
```

### `Updates were rejected because the tip of your current branch is behind`

Тоже самое что выше — кто-то запушил раньше тебя.

### Случайно запушил секрет в репу

1. Сменить ключ/токен прямо сейчас (он скомпрометирован).
2. `git filter-repo` или BFG для удаления из истории.
3. `git push --force-with-lease` после очистки.

Просто новый коммит «удалил .env» **не помогает** — старые коммиты остаются в истории и доступны через `git log`.

## Полезные ссылки

- [Работа с репозиториями GitVerse](https://gitverse.ru/docs/web-interface/repos/)
- [Работа с удалёнными репозиториями (Pro Git, ru)](https://git-scm.com/book/ru/v2/Основы-Git-Работа-с-удалёнными-репозиториями)

См. также: [02-git-basics](../02-git-basics), [03-git-branches](../03-git-branches).
