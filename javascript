// [ThitNueaHub Relay Protocol]
function doPost(e) {
  var data = JSON.parse(e.postData.contents);
  var message = data.message; // รับประโยคเดียวจากเจ้านาย

  // ส่งต่อให้ "แก้วตา" วิเคราะห์ผ่าน Gemini API
  var response = analyzeWithKaewTa(message); 
  
  // ยิงผลลัพธ์กลับไปที่ Discord Webhook
  sendToDiscord(response);
  
  return ContentService.createTextOutput("Ignite Success!");
}

function sendToDiscord(content) {
  var url = "URL_WEBHOOK_ของเจ้านาย";
  var options = {
    "method": "post",
    "contentType": "application/json",
    "payload": JSON.stringify({"content": content})
  };
  UrlFetchApp.fetch(url, options);
}
