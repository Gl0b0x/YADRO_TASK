docker build -t task-app . 
Создать снимок (в текущей директории Dockerfile)
docker run --rm -v (pwd)\test_file.txt:/build/test_file.txt task-app /build/test_file.txt
Запустить docker container и вмонтировать тестовый файл; (pwd) - путь до файла 
