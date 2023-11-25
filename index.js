import './index.css';

const handleClick = (element) => {

    // Получаем все элементы с классом "emoji"
    const emojis = document.querySelectorAll(element);

    // Добавляем обработчик события для каждого элемента
    emojis.forEach(function (emoji) {
        emoji.addEventListener('click', function () {
            // Удаляем класс "selected" у всех элементов
            emojis.forEach(function (el) {
                el.classList.remove('selected');
            });

            // Добавляем класс "selected" к выбранному элементу
            emoji.classList.add('selected');
        });
    });
}


const handleSubmit = (value) => {
    console.log(value);
}

const renderCsat = () => {
    document.body.innerHTML = `
    <form name="csat" class="iframe">
        <div class="iframe__title">Насколько вы удовлетворены подборками фильмов?</div>

        <div class="emoji__bar">
            <i class="emoji emoji__very-bad" data-value="1"></i>
            <i class="emoji emoji__bad" data-value="2"></i>
            <i class="emoji emoji__normal" data-value="3"></i>
            <i class="emoji emoji__good" data-value="4"></i>
            <i class="emoji emoji__very-good" data-value="5"></i>
        </div>

        <button class="iframe__button" type="submit">
            Отправить
        </button>
    </form>`

    handleClick('.emoji')
}

const renderCsatProfile = () => {
    document.body.innerHTML = `
    <form name="csat" class="iframe">
        <div class="iframe__title">Вам понятно как обновлять профиль?</div>

        <div class="emoji__bar">
            <i class="emoji emoji__very-bad" data-value="1"></i>
            <i class="emoji emoji__bad" data-value="2"></i>
            <i class="emoji emoji__normal" data-value="3"></i>
            <i class="emoji emoji__good" data-value="4"></i>
            <i class="emoji emoji__very-good" data-value="5"></i>
        </div>

        <button class="iframe__button" type="submit">
            Отправить
        </button>
    </form>`

    handleClick('.emoji')
}

const renderNps = () => {
    document.body.innerHTML = `
    <form name="csat" class="iframe">
        <div class="iframe__title">Насколько вы готовы рекомендовать наш сервис?</div>

        <div class="number__bar">
            <div class="number__item" data-value="1">1</div>
            <div class="number__item" data-value="2">2</div>
            <div class="number__item" data-value="3">3</div>
            <div class="number__item" data-value="4">4</div>
            <div class="number__item" data-value="5">5</div>
            <div class="number__item" data-value="6">6</div>
            <div class="number__item" data-value="7">7</div>
            <div class="number__item" data-value="8">8</div>
            <div class="number__item" data-value="9">9</div>
            <div class="number__item" data-value="10">10</div>
        </div>

        <button class="iframe__button" type="submit">
            Отправить
        </button>
    </form>`

    handleClick('.number__item')
}

document.addEventListener('DOMContentLoaded', function (e) {
    const path = window.location.pathname;
    switch (path) {
        case '/csi/feed':
            renderCsat();
            break;
        case '/csi/profile':
            renderCsatProfile()
            break;
        case '/nps':
            renderNps();
            break;

    }

    document.querySelector('.iframe').addEventListener('submit', (event) => {
        event.preventDefault();
        console.log(document.querySelector('.selected').getAttribute('data-value'))

    })


})

