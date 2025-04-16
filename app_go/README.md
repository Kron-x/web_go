description:

app_go/                   # Корневая директория app 
├── configs/
│   └── config.json       # Конфиг с параметрами запуска (порт, путь к изображениям)
├── cmd/                  
│   └── web/          
│       └── main.go       # Главный файл приложения
├── internal/         
│   └── handlers/         # HTTP-обработчики (controllers)
|       ├── home.go   
|       └── other.go      
├── pkg/              
│   ├── config/       
|   |   └── config.go     # Обработчик конфига             
│   ├── metrics/       
|   |   └── metrics.go    # Объявление кастомных метрик
│   └── postgres/       
|       └── postgres.go   # Подключение к Postgres
├── static/               
│   └── images/           # Изображения
├── go.mod                
├── go.sum                
└── README.md             # Описание проекта