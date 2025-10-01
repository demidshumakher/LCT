import torch
from torch import nn
from transformers import AutoModel 
from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
from transformers import AutoTokenizer

# ======== Настройки ========
MODEL_NAME = "DeepPavlov/rubert-base-cased"
ALL_TOPICS = [
    "Вклады", "Кредиты наличными", "Мобильное приложение",
    "Дистанционное обслуживание", "Обслуживание",
    "Кредитные карты", "Дебетовые карты", "Кешбэк",
    "Ипотека", "Накопительный счет", "Сайт"
]
ID2SENTIMENT = {0: "Отрицательно", 1: "Нейтрально", 2: "Положительно"}

# ======== Модель ========
class ABSA_RuBERT(nn.Module):
    def __init__(self, model_name, num_topics, num_sentiments):
        super().__init__()
        self.encoder = AutoModel.from_pretrained(model_name)
        hidden_size = self.encoder.config.hidden_size
        # Для каждой темы своя "голова" на 3 класса
        self.classifier = nn.Linear(hidden_size, num_topics * num_sentiments)
        self.num_topics = num_topics
        self.num_sentiments = num_sentiments

    def forward(self, input_ids, attention_mask):
        outputs = self.encoder(input_ids=input_ids, attention_mask=attention_mask)
        pooled = outputs.last_hidden_state[:, 0]  # CLS токен
        logits = self.classifier(pooled)
        # Преобразуем в [batch, topics, sentiments]
        return logits.view(-1, self.num_topics, self.num_sentiments)

device = torch.device("cuda" if torch.cuda.is_available() else "cpu")

tokenizer = AutoTokenizer.from_pretrained("tokenizer/")
model = ABSA_RuBERT(MODEL_NAME, num_topics=len(ALL_TOPICS), num_sentiments=3)
model.load_state_dict(torch.load("absa_rubert.pt", map_location=device))
model.to(device)
model.eval()

# ======== FastAPI ========
app = FastAPI()

class Review(BaseModel):
    id: int
    text: str

class RequestData(BaseModel):
    data: List[Review]

@app.post("/predict")
def predict(req: RequestData, threshold: float = 0.3):
    results = []

    for review in req.data:
        inputs = tokenizer(review.text, return_tensors="pt", truncation=True, padding=True).to(device)
        with torch.no_grad():
            logits = model(inputs["input_ids"], inputs["attention_mask"])
            probs = torch.sigmoid(logits).cpu().numpy()[0]

        topics, sentiments = [], []
        for t_idx, topic in enumerate(ALL_TOPICS):
            if probs[t_idx].max() > threshold:
                topics.append(topic)
                sentiments.append(ID2SENTIMENT[probs[t_idx].argmax()])

        results.append({
            "id": review.id,
            "text": review.text,
            "topics": topics,
            "sentiments": sentiments
        })

    return {"results": results}
