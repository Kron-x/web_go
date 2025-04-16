# My web-server with monitoring tools:
1. Web application: [goserv.duckdns.org](http://goserv.duckdns.org)
2. Active ports: 
   - Prometeus: 9090
   - Grafana: 3000 
   - Alertmanager: 9093
   
# **Запуск**
1. `git clone github.com/Kron-x/web_go`
2. `cd web_go`
3. Add in secrets in ansible/roles/app/defaults/main.yaml
   - telegram_bot_token: your_bot_token
   - telegram_chat_id: your_chat_id
4. `docker-compose up -d`