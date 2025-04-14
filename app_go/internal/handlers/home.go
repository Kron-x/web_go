package handlers

import (
    "fmt"
    "net/http"
	"app_go/pkg/metrics"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	
	metrics.PageVisits.WithLabelValues("/").Inc()
    metrics.ActiveUsers.Inc()
    defer metrics.ActiveUsers.Dec()

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprintf(w, `
		<h1>Welcome to the Image Server</h1>
		<img src="/images/image1.jpg" alt="Image 1" width="300">
		<img src="/images/image2.jpg" alt="Image 2" width="300">
        <div style="position: fixed; bottom: 20px; right: 20px;">
			<button onclick="window.location.href='/new-dimension'">Перейти в новое измерение</button>
		</div>
        <div style="text-align: center; margin-top: 20px;">
			<input type="text" id="textInput" placeholder="Введите текст" style="padding: 10px; font-size: 16px;">
			<button onclick="submitText()" style="padding: 10px; font-size: 16px;">Отправить</button>
		</div>
		<div id="result" style="text-align: center; margin-top: 20px; font-size: 18px;"></div>
		<script>
			function submitText() {
				const input = document.getElementById("textInput");
				const result = document.getElementById("result");
				const text = input.value.trim();

				if (text) {
					fetch("/submit-text", {
						method: "POST",
						headers: {
							"Content-Type": "application/json",
						},
						body: JSON.stringify({ text: text }),
					})
					.then(response => response.json())
					.then(data => {
						result.textContent = "Ответ сервера: " + data.message;
						input.value = "";
					})
					.catch(error => {
						result.textContent = "Ошибка: " + error.message;
					});
				} else {
					result.textContent = "Пожалуйста, введите текст.";
				}
			}

			// Обработка нажатия Enter
			document.getElementById("textInput").addEventListener("keyup", function(event) {
				if (event.key === "Enter") {
					submitText();
				}
			});
		</script>
	`)
}