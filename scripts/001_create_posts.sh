#!/usr/bin/env bash

CURL="curl"
CALL="grpcurl --plaintext"
URL="localhost:6090"

function call_CreateOriginalPost() {
    $CALL -d @ "$URL" post_processor.PostsProcessor/CreateOriginalPost
}

echo "start"
call_CreateOriginalPost <<EOF
{
  "text": "Сегодня был прекрасный закат в горах. Природа напоминает, как важно иногда замедляться и наслаждаться моментом. 🏔️",
  "link_original_post": "https://example.com/post/12345",
  "original_channel": "telegram",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Только что попробовал новый рецепт пасты с соусом песто. Просто, быстро и невероятно вкусно! Делитесь своими любимыми рецептами в комментариях. 🍝",
  "link_original_post": "",
  "original_channel": "instagram",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Завершил курс по машинному обучению. Ощущение, что открыл новую вселенную возможностей. Следующая цель — применить знания на реальном проекте. 🤖",
  "link_original_post": "https://medium.com/@user/ml-journey",
  "original_channel": "twitter",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Утренняя пробежка в парке — лучшее начало дня. Свежий воздух, пение птиц и никаких забот. Рекомендую всем попробовать! 🏃‍♂️",
  "link_original_post": "https://strava.com/activities/987654321",
  "original_channel": "facebook",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Новый трек от моего любимого исполнителя вышел! Готов залипнуть на весь вечер. Какую музыку слушаете вы? 🎧",
  "link_original_post": "https://open.spotify.com/track/abc123",
  "original_channel": "vk",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Прочитал книгу «Атомные привычки». Очень вдохновляет на маленькие, но постоянные изменения. Советую всем, кто хочет стать лучше. 📚",
  "link_original_post": "",
  "original_channel": "telegram",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Вчера впервые покатался на сноуборде. Упал 50 раз, но улыбка не сходила с лица. Спорт — это радость, даже когда не сразу получается. 🏂",
  "link_original_post": "https://youtu.be/xyz789",
  "original_channel": "instagram",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Цитата дня: «Ваша работа заполнит большую часть жизни, и единственный способ быть полностью довольным — делать то, что вы считаете великим делом». Стив Джобс. 💡",
  "link_original_post": "https://linkedin.com/posts/user-456",
  "original_channel": "twitter",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Фотография ночного города с крыши небоскрёба. Москва выглядит как светящаяся карта сокровищ. 🌃",
  "link_original_post": "",
  "original_channel": "facebook",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Собрал свой первый ПК своими руками. Запустился с первого раза! Спецификации в описании. Кто тоже собирал — делимся опытом. 💻",
  "link_original_post": "https://pcpartpicker.com/b/123abc",
  "original_channel": "vk",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Только что вернулся из поездки в горы! Незабываемые виды и чистый воздух 🏔️ #путешествия",
  "link_original_post": "https://t.me/travel_diary/123",
  "original_channel": "travel_channel",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Приготовил сегодня пасту карбонара по домашнему рецепту. Получилось очень вкусно! 🍝 #еда",
  "link_original_post": "https://t.me/food_lover/456",
  "original_channel": "food_blog",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Обновил смартфон до новой прошивки — батарея держится на 20% дольше. Технологии творят чудеса! 📱 #tech",
  "link_original_post": "https://t.me/tech_news/789",
  "original_channel": "gadgets",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Никогда не сдавайся. Даже маленький шаг вперёд — это прогресс. 💪 #мотивация",
  "link_original_post": "https://t.me/inspire/101",
  "original_channel": "motivation_daily",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Утренняя пробежка в парке подарила заряд бодрости на весь день. Бег — лучшее лекарство от стресса 🏃‍♂️ #спорт",
  "link_original_post": "https://t.me/run_world/202",
  "original_channel": "fitness_life",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Дочитал «1984» Оруэлла. Пугающе актуально, даже спустя десятилетия. Советую всем! 📚 #книги",
  "link_original_post": "https://t.me/book_club/303",
  "original_channel": "reading_room",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Вчера посмотрел новый «Дюна 2» — визуальный шедевр и отличный сюжет. Однозначно в топ года! 🎬 #кино",
  "link_original_post": "https://t.me/movie_talk/404",
  "original_channel": "cinema_fan",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Открыл для себя группу «The Midnight» – синтвейв на вечер. Создаёт атмосферу ночного города 🎧 #музыка",
  "link_original_post": "https://t.me/synth_waves/505",
  "original_channel": "music_discovery",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Посадил сегодня дуб у дома. Через 20 лет здесь будет настоящая аллея 🌳 #природа #экология",
  "link_original_post": "https://t.me/green_life/606",
  "original_channel": "eco_blog",
  "theme": "test"
}
EOF

call_CreateOriginalPost <<EOF
{
  "text": "Минимализм в гардеробе: 5 базовых вещей, которые подходят ко всему. Освобождает голову и место в шкафу 👕 #мода",
  "link_original_post": "https://t.me/style_tips/707",
  "original_channel": "fashion_guide",
  "theme": "test"
}
EOF
