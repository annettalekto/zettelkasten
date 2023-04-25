быстрый старт:

[marktext/EDITING.md at master · marktext/marktext · GitHub](https://github.com/marktext/marktext/blob/master/docs/EDITING.md)

<br>

## Кратко основное:

Ctrl + E -- режим редактора, Ctrl + Shift + J -- режим фокуса, меню View.

Можно делать табличку. Как вставить незнаю, просто черточки ставишь, он сам нарисует.

###### Знаки с новой строки:

1. Заголовок # решетка от 1 до 6 в начале строки это.

2. Цитировать >, могут быть вложенными >>

3. Отделить чертой из тире ---

4. ()[] вставить ссылку сначала описание и в круглых ссылка

5. !()[] вставить картинку

6. список цифровой (как этот), кругляшками - или * или + или, а чек лист - [ ]

###### Выделять как кавычками:

- жирный `**` две звездочки **жирный** (он же заголовок 6)
* курсив одна звездочка *курсив* 

Меню Format.

###### Блоки, выделять на первой и последней строке:

- [x] код `  через три кавычки клавиши ё.
- [ ] маркер mark как html тег в угловых скобках, открыть и закрыть /mark

###### Дополнительно:

Таблицу, формулу, html блок, диаграмму вставь через меню Paragraph, Format или ПКЛМ.

:smile: вставить эмодзи через : по краям английского слова  :улыбка: 

Проверка орфографии включается в настройках File->Preferences(ctrl+,)->Spelling. 
Чтобы отобразить символ для форматирования нужно взять его в кавычки ё `<mark>`.

Для пространства между параграфами использовать `<br>` окруженный пустыми строками.

<br>

### меню View:

- Ctrl + E -- режим редактора
- Ctrl + Shift + G -- режим печатной машинки, строка которую пишешь, всегда на уровне глаз.
- Ctrl + Shift + J -- режим фокуса - только один текущий параграф четко отображается. Когда переходишь в режим редактора для правок и возвращаешся обратно, тебе не кидает на последний параграф, как в остальных ражимах на кой то хрен. Но если стоять на пустой строке, все равно перекидывает.
- Ctrl + J -- показать меню с боку
- Ctrl + Shift + B -- показать сверху файл как отдельную вкладку  

<br>

### Знак @ для вывода меню Basic block, там:

- параграф, хз что это просто символ отмечает параграф

- горизонтальная линия --- три дефиса и новая строка

- заголовки количество знаков # решетка, определяет размер 

- таблица  лучше через меню

- вставить формулу $$ два знака доллара, только не как кавычки по бокам, а как параграф на первой строке и на последней

$$
X*y=a/b
$$

- HTML блок через  < div > < /div>
- блок кода три знака ` на верхней (тут можно указать язык) и нижней строке

```go
fmt.Println(x)
```

- Цитировать
  
  > quote цитата > знак в начале строки. Это как бы цитата.
  > Можно делать перенос строки. 
  > 
  > > И даже вложенную цитату.

<br>

# Заголовок 1 (одна решетка)

## Заголовок 2

### Заголовок 3

#### Заголовок 4

##### Заголовок 5

###### Заголовок 6 (всего может быть 6, не отличается от жирного шрифта, но выделяет строку как заголовок пространством)

<br>

### Расшифровка меню пклм:

(все слитно с текстом, без пробелов)

- **Bold**                      жирный ** две звездочки по бокам 
- *Italics*                     курсив * одна звездочка по бокам 
- <u>Underline</u>             подчеркнуть < u > < / u > как теги xtml, без пробелов
- ~~Strikethrough~~      зачеркнуть ~~ две волнистые черты по бокам (тильда)
- <mark>Inline code</mark>           маркер `<mark> </mark>`
- `Inline math formulas` мат формула ` кавычка с ё которая
- [Create link]()           ссылка [ ] ( ) описание в квадратных, сама ссылка в круглых
- ![Create image]() 
  картинка так ставиться ! [ ] ( ) воскл. знак, ссылка в квадратных, описание в круглых  
- Remove formatting ??

<br>

### Списки

###### bullet list:

- список с круглешками тире или звездочка, пробел в начале (- )
* или звездочкой
  
  - такой круглешок (вложеный список) смести тире на один пробел ( - )
    - можно так
      - или так
        * хз зачем тире и звезда делаюит одно и тоже 
+ чтобы нормально смотрлся список в режиме

+ просмотра 
  
  - и в режиме
  - форматирования
  - можно для себя использовать такие вариманты

###### цифры:

1. раз

2. два

3. три

###### чек лист:

- [ ] сделать раз

- [x] сделать два

- [ ] сделать три

###### чек лист со ссылками

- [ ] [foo](#bar)
- [ ] [baz](#qux)
- [ ] [fez](#faz)

<br>

---

<mark>Как отделить один параграф от другого новой строкой? </mark>

Символ `<br>` окруженный пустыми строками строками вокруг.

<br>

----

По идее можно в начале организовать список ссылок на разные параграфы, но у меня не получается. 

# Table of Contents

* [Chapter 1](#chapter-1)
* [Chapter 2](#chapter-2)
* [Chapter 3](#chapter-3)

will jump to these sections:

## Chapter 1

Content for chapter one.

## Chapter 2

Content for chapter one.

## Chapter 3 <a name="chapter-3"></a>

Content for chapter one.

-------------------------

Доп. заметки

вот такая штука:

[image](widgets\img\progress_bar.gif)

работает для marktext, а вот на github гифку не отображает (с картинвками так же было вроде). Работает так (ширину не обязательно указывать):
<img src="widgets\img\progress_bar.gif"  width="400"/>