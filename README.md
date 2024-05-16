docker build -t task-app . 
docker run --rm -v (pwd)\test_file.txt:/build/test_file.txt task-app /build/test_file.txt
(pwd) - путь до файла 
