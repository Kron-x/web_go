description:

app_go/                   # Корневая директория проекта
├── configs/
│   └── config.json       # Конфиг с параметрами запуска (порт, путь к изображениям)
├── cmd/                  # Основные исполняемые файлы (точки входа)
│   └── web/          
│       └── main.go       # Главный файл приложения
├── internal/         
│   └── handlers/         # HTTP-обработчики (controllers)
|       ├── home.go   
|       └── other.go      
├── pkg/              
│   ├── config/       
|   |   └── config.go     # Обработчик конфига             
│   └── postgres/       
|       └── postgres.go   # Подключение к Postgres
├── static/               # Статические файлы (изображения)
│   └── images/       
├── go.mod                
├── go.sum                
└── README.md             # Описание проекта