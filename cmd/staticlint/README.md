# staticlint

staticlint это агрегатор уже существующих анализаторов кода позволяющий добавить дополнительные проверки

## Analyzer assign

origin
doc [https://pkg.go.dev/golang.org/x/tools@v0.19.0/go/analysis/passes/assign#section-documentation](https://pkg.go.dev/golang.org/x/tools@v0.19.0/go/analysis/passes/assign#section-documentation)

**assign:** check for useless assignments

This checker reports assignments of the form x = x or a[i] = a[i]. These are almost always useless, and even when they
aren't they are usually a mistake.

## Analyzer SA (staticcheck)

Ссылка на все SA проверки и их обозначения [https://staticcheck.io/docs/checks/](https://staticcheck.io/docs/checks/)

Проверки SA анализатора staticcheck разбиты на группы и имеют индекс вида SA1001, SA2002, SA3003:

- **SA1???** — неправильное использование стандартных библиотек;
- **SA2???** — проблемы с многопоточностью;
- **SA3???** — проблемы с тестами;
- **SA4???** — бесполезный код;
- **SA5???** — ошибочный код;
- **SA6???** — проблемы с производительностью;
- **SA9???** — сомнительные конструкции кода, c высокой вероятностью ошибочные;

## Analyzer QF1006 (staticcheck) - Lift if+break into loop condition

origin doc [https://staticcheck.io/docs/checks/#QF1006](https://staticcheck.io/docs/checks/#QF1006)

**Before:**

```go
for {
if done {
break
}
...
}
```

**After:**

```go
for !done {
...
}
```


