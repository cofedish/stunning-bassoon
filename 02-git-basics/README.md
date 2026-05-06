# 02. Системы контроля версий — основы Git

## Что такое VCS

Система контроля версий хранит историю изменений файлов: кто, когда и что поменял. Позволяет откатиться, посмотреть разницу, работать командой без перезаписей.

**Git** — самая популярная распределённая VCS. «Распределённая» = у каждого разработчика полная копия истории, не нужен постоянно живой сервер.

## Основные понятия

| Термин | Что это |
|---|---|
| Repository (репозиторий) | Папка с проектом + папка `.git` с историей |
| Working tree (рабочая копия) | Файлы, которые ты видишь и редактируешь |
| Index / staging area (индекс) | «Корзина» того, что попадёт в следующий коммит |
| Commit | Снимок состояния + автор + сообщение + ссылка на родителя |
| HEAD | Указатель на текущий коммит / ветку |
| Branch (ветка) | Подвижный указатель на коммит |
| Remote | Удалённая копия репы (на GitHub / GitVerse / Gitea) |
| Origin | Стандартное имя для основного remote |

## Жизненный цикл файла

```
[не отслеживается] --- git add ---> [staged] --- git commit ---> [committed]
        ^                              |                            |
        |--- редактирование ----- [modified] <----------------------|
```

## Базовые команды

### Создать репу

```bash
git init                       # из существующей папки
git clone <url>                # склонировать чужую
git clone git@github.com:u/r.git
```

### Посмотреть состояние

```bash
git status                     # что изменено / staged / untracked
git diff                       # diff между working tree и index
git diff --staged              # diff между index и последним коммитом
git log                        # история
git log --oneline --graph --all # компактно с веткой
git show <commit>              # содержимое конкретного коммита
```

### Записать изменения

```bash
git add file.py                # добавить файл в index
git add .                      # всё в текущей папке
git add -p                     # интерактивно по кускам
git commit -m "сообщение"      # создать коммит
git commit -am "..."           # add+commit для уже отслеживаемых
git commit --amend             # переписать последний коммит (если ещё не запушен!)
```

### Откатиться

```bash
git restore file.py            # откатить изменения в working tree
git restore --staged file.py   # убрать из index, оставить в working tree
git reset --hard HEAD          # ОПАСНО: всё откатить к последнему коммиту
git revert <commit>            # создать новый коммит, отменяющий старый
```

## SSH-ключи (нужны для push/pull без пароля)

```bash
# 1. Сгенерировать
ssh-keygen -t ed25519 -C "you@example.com"
# жмёшь Enter трижды (путь по умолчанию + без passphrase)

# 2. Посмотреть публичный ключ
cat ~/.ssh/id_ed25519.pub
# на Windows: type %USERPROFILE%\.ssh\id_ed25519.pub

# 3. Скопировать содержимое и добавить в GitHub:
#    Settings → SSH and GPG keys → New SSH key

# 4. Проверить
ssh -T git@github.com
```

## Pull request / Merge request

PR (GitHub) и MR (GitLab/GitVerse) — это одно и то же: запрос на слияние твоей ветки в основную. Workflow:

1. Делаешь форк / создаёшь ветку.
2. Коммитишь, пушишь.
3. Идёшь в веб-интерфейс и нажимаешь **New pull request**.
4. Описываешь что и зачем поменял.
5. Ревьюеры комментируют → правишь → коммитишь → пуш.
6. Когда ок — merge.

## Конфиг (один раз поставить)

```bash
git config --global user.name "Имя Фамилия"
git config --global user.email "you@example.com"
git config --global init.defaultBranch main
git config --global pull.rebase false   # стратегия pull: merge (или true для rebase)
```

## .gitignore

Файлы, которые git игнорирует. Лежит в корне репы.

```gitignore
# Python
__pycache__/
*.pyc
.venv/

# Go
/bin/
*.exe

# IDE
.idea/
.vscode/

# секреты
.env
*.key
```

## Полезные ссылки

- [Что такое система управления версиями Git (GitVerse)](https://gitverse.ru/blog/articles/development/195-chto-takoe-sistema-kontrolya-versij-git/)
- [О системе контроля версий (официальная книга Pro Git, ru)](https://git-scm.com/book/ru/v2/Введение-О-системе-контроля-версий)

См. также: [03-git-branches](../03-git-branches), [09-git-repos](../09-git-repos).
