cd "E:\StudyData\GO\WebSupervisor\WebSupervisor-microservice" 

go build -o ./microservice/cache-service/run.exe ./microservice/cache-service
go build -o ./microservice/crawler-service/run.exe ./microservice/crawler-service
go build -o ./microservice/monitor-service/run.exe ./microservice/monitor-service
go build -o ./microservice/notifier-service/run.exe ./microservice/notifier-service
go build -o ./microservice/parser-service/run.exe ./microservice/parser-service


cd microservice/cache-service
./run.bat 
cd ../monitor-service
./run.bat 
cd ../parser-service
./run.bat 
cd ../notifier-service
./run.bat
cd ../crawler-service
./run.bat
