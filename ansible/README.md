```plaintext
ansible/                    # Корневая директория ansible
├── ansible.cfg             # Конфиг с параметрами запуска
├── group_vars/                  
│   └── all.yaml/           # Переменные по умолчанию (заглушки)
├── inventory/         
│   ├── group_vars/       
|   |   └── webserver.yaml  # Переменные вебсервера   
|   └── hosts               # Инвентарь для плейбука   
├── playbook.yaml           # Плейбук
├── README.md  
└── roles/              
    ├── app_local/          # Запуск при target=localhost
    |   ├── defaults/            
    │   │   └── main.yaml
    │   └── tasks/
    │       └── main.yaml
    ├── app_webserver/      # Запуск default (для pipeline)
    |   ├── defaults/            
    │   │   └── main.yaml
    │   └── tasks/
    │       └── main.yaml
    └── docker/             # Проверка установки docker compose
        ├── defaults/
        │   └── main.yaml
        └── tasks/
            └── main.yaml
