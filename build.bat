 cd "e:\StudyData\GO\WebSupervisor-microservice\new" 

  go build -o ./bin/cache-service.exe ./cache-service
  go build -o ./bin/crawler-service.exe ./crawler-service
  go build -o ./bin/monitor-service.exe ./monitor-service
  go build -o ./bin/notifier-service.exe ./notifier-service
  go build -o ./bin/parser-service.exe ./parser-service

