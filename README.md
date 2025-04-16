## My web-server with monitoring tools:
1. Web application: [goserv.duckdns.org](http://goserv.duckdns.org)
2. Active ports: 
   - Prometeus: 9090
   - Grafana: 3000 
   - Alertmanager: 9093
 
## **Starting**

### **Pre-install**
- git
- python
- ansible (after python: `pip install ansible`)

### **Steps**
1. `git clone github.com/Kron-x/web_go`
2. `cd web_go`
3. Add secrets in ansible/roles/app/defaults/main.yaml
   - `echo -e "\ntelegram_bot_token: your_bot_token" >> ansible/roles/app_local/defaults/main.yaml`
   - `echo -e "\ntelegram_chat_id: your_chat_id" >> ansible/roles/app_local/defaults/main.yaml`  
   (replace **your_bot_token** and **your_chat_id**)
4. `ansible-playbook playbook.yaml -e "target=localhost"`